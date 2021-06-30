package monitor

import (
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/pkg/errors"
)

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func SendMail(cfg *Config, subject string, msg string) error {
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	auth := smtp.PlainAuth("", cfg.Mail.From, cfg.Mail.Password, smtpServer.host)

	msgBody := fmt.Sprintf("Subject: %v\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n"+
		"%v\n", subject, msg)

	if err := smtp.SendMail(smtpServer.Address(), auth, cfg.Mail.From, cfg.Mail.To, []byte(msgBody)); err != nil {
		return err
	}

	return nil
}

type ModelMap map[string]float32
type ProductMap map[string]ModelMap
type StockMap map[string]ProductMap

// add a model by either inserting a new one, or comparing with the price of an existing model in the map
func (modelMap ModelMap) AddModel(modelNumber string, price float32) {
	if existingPrice, ok := modelMap[modelNumber]; ok {
		if price < existingPrice {
			modelMap[modelNumber] = price
		}
	} else {
		modelMap[modelNumber] = price
	}
}

// adds a range of models
func (modelMap ModelMap) AddModels(models []Model) {
	for _, model := range models {
		if model.Error != nil {
			modelMap.AddModel(model.Number, model.Price)
		}
	}
}

// converts the list to a map. We don't reduce by price, because the model slice should already be price-reduced
func ToModelMap(models []Model) (ModelMap, error) {
	modelMap := make(ModelMap)
	for _, model := range models {
		if model.Error == nil {
			modelMap[model.Number] = model.Price
		}
	}

	if len(modelMap) == 0 {
		return modelMap, errors.New("all models have an error, couldn't create map")
	}

	return modelMap, nil
}

func (m StockMap) AddResult(r PriceCheckResponse) {
	if productMap, ok := m[r.Job.SiteName]; ok {
		if modelMap, ok := productMap[r.Job.ProductName]; ok {
			modelMap.AddModels(r.Models)
		} else if len(r.Models) > 0 {
			if mp, err := ToModelMap(r.Models); err == nil {
				productMap[r.Job.ProductName] = mp
			}
		}
	} else if len(r.Models) > 0 {
		productMap = make(ProductMap)
		if mp, err := ToModelMap(r.Models); err == nil {
			productMap[r.Job.ProductName] = mp
		}
		m[r.Job.SiteName] = productMap
	}
}

func GenerateResultsEmail(stockMap StockMap) (string, bool) {
	hasStock := false
	body := "<h1>Looks like some GPUs are finally available!</h1>"
	for siteName, productMap := range stockMap {
		body += fmt.Sprintf("<div><h2>stock at <strong>%v</strong></h2>", siteName)
		for productName, models := range productMap {
			body += fmt.Sprintf("<h3>Models/price available for <strong>%v</strong>:</h3><ul>", productName)
			for modelNumber, price := range models {
				hasStock = true
				body += Hprintf("<li><strong>%v</strong> is available for <strong>%v</strong></li>",
					modelNumber,
					strconv.Itoa(int(price)))
			}
			body += "</ul>"
		}
		body += "</div>"
	}
	return body, hasStock
}

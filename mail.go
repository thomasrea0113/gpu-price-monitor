package monitor

import (
	"fmt"
	"net/smtp"
	"strconv"
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

type StockMap map[string]map[string]map[string]float32

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

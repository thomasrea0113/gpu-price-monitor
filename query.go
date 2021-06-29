package monitor

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/thomasrea0113/gpu-price-monitor/domain"
)

func execPuppeteer(url string) (string, error) {
	cmd := exec.Command("node", "index.js", url)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	var outString string
	for {
		if line, err := out.ReadString(byte('\n')); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Printf("unexpected error: %v", err)
				return "", err
			}
		} else {
			outString += line
		}
	}

	return outString, nil
}

// ensures all arguments to sprinf are properly escaped
func uSprintf(format string, vv ...string) string {
	vCopy := make([]interface{}, len(vv))
	for i, v := range vv {
		vCopy[i] = url.PathEscape(v)
	}
	return fmt.Sprintf(format, vCopy...)
}

func QueryBestBuy(j domain.PriceCheckJob, keyword string) []domain.Model {
	models := make([]domain.Model, 0)

	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	url := uSprintf(j.Site.UrlFormat, keyword, j.Product.Name)

	html, err := execPuppeteer(url)
	if err != nil {
		log.Printf("error execing puppeteer: %v", err)
		return models
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("error parsing document: %v", err)
		return models
	}

	doc.Find(".sku-item").Each(func(i int, s *goquery.Selection) {
		// if sold out... we're done here.
		if strings.ToLower(s.Find(".fulfillment-fulfillment-summary").Text()) == "sold out" {
			return
		}

		// price text has the leading $, need to exclude it
		priceStr := s.Find(".priceView-customer-price span[aria-hidden=true]").Text()

		if priceStr != "" {
			priceStr = strings.ReplaceAll(priceStr[1:], ",", "")
		}

		price, err := strconv.ParseFloat(priceStr, 32)
		if err != nil {
			log.Printf("Error converting price to int: %v", err)
			return
		}

		modelNumber := s.Find(".sku-value").Text()

		url, ok := s.Find(".sku-header a").Attr("href")
		if !ok {
			log.Printf("model has no details link")
			return
		}

		if price < float64(j.Product.PriceThreshhold) {
			models = append(models, domain.Model{
				Number: modelNumber,
				Price:  float32(price),
				Url:    "https://bestbuy.com" + url})
		}
	})

	return models
}

func QueryWalMart(j domain.PriceCheckJob, keyword string) []domain.Model {
	models := make([]domain.Model, 0)

	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	return models
}

func QueryNewegg(j domain.PriceCheckJob, keyword string) []domain.Model {
	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	return make([]domain.Model, 0)
}

func QueryMicroCenter(j domain.PriceCheckJob, keyword string) []domain.Model {
	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	return make([]domain.Model, 0)
}

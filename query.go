package monitor

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func QueryBestBuy(j PriceCheckJob) []Model {
	models := make([]Model, 0)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*j.PageContent))
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

		if priceStr == "" {
			log.Println("couldn't get item price")
			return
		}

		priceStr = strings.ReplaceAll(priceStr[1:], ",", "")
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

		if price < float64(j.PriceThreshhold) {
			models = append(models, Model{
				Number: modelNumber,
				Price:  float32(price),
				Url:    "https://bestbuy.com" + url})
		}
	})

	return models
}

func QueryWalMart(j PriceCheckJob) []Model {
	models := make([]Model, 0)

	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	return models
}

func QueryNewegg(j PriceCheckJob) []Model {
	models := make([]Model, 0)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*j.PageContent))
	if err != nil {
		log.Printf("error parsing document: %v", err)
		return models
	}

	doc.Find(".item-cell").Each(func(i int, s *goquery.Selection) {
		priceText := s.Find(".price-current strong,sup").Text()

		if priceText == "" {
			log.Println("couldn't get item price")
			return
		}

		price, err := strconv.ParseFloat(strings.ReplaceAll(priceText, ",", ""), 32)
		if err != nil {
			log.Printf("error getting price: %v", err)
			return
		}

		// navigating to details URL, since the search page doesn't include model number
		detailsUrl, ok := s.Find("a[title='View Details']").Attr("href")
		if detailsUrl == "" || !ok {
			log.Println("couldn't get details url")
			return
		}

		// TODO get model number
		modelNumber := strings.Split(detailsUrl, "/")[3]
		if modelNumber == "" {
			log.Println("model number not found on details page")
			return
		}

		models = append(models, Model{Price: float32(price), Number: modelNumber, Url: detailsUrl})
	})

	return models
}

func QueryMicroCenter(j PriceCheckJob) []Model {
	// TODO use colly to get all the models returned for a given product + keyword, across first couple pages
	return make([]Model, 0)
}

package monitor

import (
	"strconv"

	"github.com/pkg/errors"
)

// a function that is responsible for scraping one keyword query for a given job
type KeywordScrapeFunc = func(PriceCheckJob) []Model

// HOF that will break up the keyword queries, and aggregrate the returned models
// TODO bubble errors to resp.Error
func reduce(work KeywordScrapeFunc) func(PriceCheckJob) []Model {
	return func(j PriceCheckJob) []Model {
		// reduce keyword model-map into one flat list of models. If the same model appears twice, take
		// the one with the lower price
		models := work(j)
		modelMap := make(map[string]Model)

		// handle model duplicates, by keeping the listing with the lowest price
		reducedModels := make([]Model, 0, len(models))
		i := 0
		for _, model := range models {
			existing, ok := modelMap[model.Number]
			if !ok || model.Price < existing.Price {
				modelMap[model.Number] = model
				reducedModels = append(reducedModels, model)
				i++
			}
		}

		return reducedModels
	}
}

func GenerateUrl(site Site, product Product, keyword string) (string, error) {
	switch site.Name {
	case "Best Buy":
		return Uprintf(site.UrlFormat, keyword, product.Name), nil
	case "Newegg":
		return Uprintf(site.UrlFormat, keyword, product.Name, strconv.Itoa(product.PriceThreshhold)), nil
	default:
		return "", errors.Errorf("Unrecognized site name: %v", site.Name)
	}
}

func scrape(job interface{}) interface{} {
	j := job.(PriceCheckJob)
	resp := PriceCheckResponse{Job: j}

	switch j.SiteName {
	case "Best Buy":
		resp.Models = reduce(QueryBestBuy)(j)
	case "Wal-Mart":
		resp.Models = reduce(QueryWalMart)(j)
	case "Micro Center":
		resp.Models = reduce(QueryMicroCenter)(j)
	case "Newegg":
		resp.Models = reduce(QueryNewegg)(j)
	default:
		resp.Error = errors.Errorf("Unknown site name: %v", j.SiteName)
	}

	return resp
}

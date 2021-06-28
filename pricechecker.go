package monitor

import (
	"github.com/pkg/errors"
	"github.com/thomasrea0113/gpu-price-monitor/domain"
)

// a function that is responsible for scraping one keyword query for a given job
type KeywordScrapeFunc = func(domain.PriceCheckJob, string) []domain.Model

// HOF that will break up the keyword queries, and aggregrate the returned models
// TODO bubble errors to resp.Error
func scraper(work KeywordScrapeFunc) func(domain.PriceCheckJob) []domain.Model {
	return func(j domain.PriceCheckJob) []domain.Model {
		// a model map of model number to the model struct. used to easily identify duplicate models
		modelMap := make(map[string]domain.Model)

		// TODO is there a simpler way to cast a slice to a slice of empty interfaces?
		keywordJobs := make([]interface{}, len(j.Product.AdditionalKeywords))
		for i, k := range j.Product.AdditionalKeywords {
			keywordJobs[i] = k
		}

		worker := NewWorkerPool(5, keywordJobs, func(t interface{}) interface{} {
			return work(j, t.(string))
		})

		worker.StartGroup().Wait()

		// reduce keyword model-map into one flat list of models. If the same model appears twice, don't overwrite. Assume
		// that the keyword search returned the same exact results as a previous query
		for model := range worker.Results {
			mod := model.(domain.Model)
			if _, ok := modelMap[mod.Number]; !ok {
				modelMap[mod.Number] = mod
			}
		}

		// get just the model value, and create a new slice of just the values
		models := make([]domain.Model, 0, len(modelMap))
		for _, v := range modelMap {
			models = append(models, v)
		}

		return models
	}
}

func scrape(job interface{}) interface{} {
	j := job.(domain.PriceCheckJob)
	resp := domain.PriceCheckResponse{Job: j}

	switch j.Site.Name {
	case "Best Buy":
		resp.Models = scraper(QueryBestBuy)(j)
	case "Wal-Mart":
		resp.Models = scraper(QueryWalMart)(j)
	case "Micro Center":
		resp.Models = scraper(QueryMicroCenter)(j)
	case "Newegg":
		resp.Models = scraper(QueryNewegg)(j)
	default:
		return domain.PriceCheckResponse{Error: errors.Errorf("Unknown site name: %v", j.Site.Name)}
	}

	return resp
}

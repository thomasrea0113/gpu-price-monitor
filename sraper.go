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

		worker := NewWorkerPool(100, keywordJobs, func(t interface{}) interface{} {
			return work(j, t.(string))
		})

		worker.Start()

		// reduce keyword model-map into one flat list of models. If the same model appears twice, take
		// the one with the lower price
		for i := 0; i < worker.JobCount; i++ {
			result := <-worker.Results
			models := result.([]domain.Model)
			for _, model := range models {
				if existing, ok := modelMap[model.Number]; !ok {
					if model.Price < existing.Price {
						modelMap[model.Number] = model
					}
				}
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
		resp.Error = errors.Errorf("Unknown site name: %v", j.Site.Name)
	}

	return resp
}

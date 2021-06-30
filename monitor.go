package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getConfigWithOverrides(overrideCfg *Config) (*Config, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return cfg, fmt.Errorf("error loading config: %v", err)
	}

	if overrideCfg != nil {
		if err := cfg.Merge(*overrideCfg); err != nil {
			return cfg, fmt.Errorf("error merging override config with base config: %v", err)
		}
	}

	return cfg, nil
}

func MonitorProducts(w http.ResponseWriter, r *http.Request) {
	// reading message body
	var msg RequestMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		writeAndLog(w, "could not read request body: %v", err)
		return
	}

	cfg, err := getConfigWithOverrides(msg.ConfigOverrides)
	if err != nil {
		writeAndLog(w, "couldn't get config: %v", err)
		return
	}

	jobs, err := cfg.GetJobs()
	if err != nil {
		writeAndLog(w, "error getting jobs: %v", err)
	}

	pool := NewWorkerPool(100, jobs, scrape)
	pool.Start()

	// TODO add some functions to reduce copied code
	StockMap := make(StockMap)
	for i := 0; i < pool.JobCount; i++ {
		result := (<-pool.Results).(PriceCheckResponse)

		if productMap, ok := StockMap[result.Job.SiteName]; ok {
			if modelMap, ok := productMap[result.Job.ProductName]; ok {
				for _, model := range result.Models {
					if existingPrice, ok := modelMap[model.Number]; ok {
						if model.Price < existingPrice {
							modelMap[model.Number] = model.Price
						}
					} else {
						modelMap[model.Number] = model.Price
					}
				}
			} else if len(result.Models) > 0 {
				modelMap = make(map[string]float32)
				for _, model := range result.Models {
					if existingPrice, ok := modelMap[model.Number]; ok {
						if model.Price < existingPrice {
							modelMap[model.Number] = model.Price
						}
					} else {
						modelMap[model.Number] = model.Price
					}
				}

				productMap[result.Job.ProductName] = modelMap
			}
		} else if len(result.Models) > 0 {
			productMap = make(map[string]map[string]float32)
			modelMap := make(map[string]float32)
			for _, model := range result.Models {
				if existingPrice, ok := modelMap[model.Number]; ok {
					if model.Price < existingPrice {
						modelMap[model.Number] = model.Price
					}
				} else {
					modelMap[model.Number] = model.Price
				}
			}
			productMap[result.Job.ProductName] = modelMap

			StockMap[result.Job.SiteName] = productMap
		}
	}

	email, hasStock := GenerateResultsEmail(StockMap)

	if hasStock {
		// if *cfg.SendEmails {
		SendMail(cfg, "GPU Stock Available!", email)
		// } else {
		// 	log.Println("skipping email")
		// }
		fmt.Fprint(w, email)
	} else {
		writeAndLog(w, "No stock available")
	}

}

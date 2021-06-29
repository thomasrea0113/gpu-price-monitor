package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	hasStock := false
	body := "<h2>Looks like some GPUs are finally available!</h2>"
	for i := 0; i < pool.JobCount; i++ {
		result := (<-pool.Results).(PriceCheckResponse)

		if len(result.Models) > 0 {
			hasStock = true
			modelStr := "<ul>"
			for _, model := range result.Models {
				if model.Error == nil {
					modelStr += Hprintf("<li>%v (model: %v) for %v </li>",
						result.Job.ProductName,
						model.Number,
						strconv.Itoa(int(model.Price)))
				}
			}
			modelStr += "</ul>"

			body += Hprintf("<div><h3>%v</h3>%v</div>", result.Job.SiteName, modelStr)
		}
	}

	if hasStock {
		SendMail(cfg, "GPU Stock Available!", body)
		fmt.Fprint(w, body)
	} else {
		writeAndLog(w, "No stock available")
	}

}

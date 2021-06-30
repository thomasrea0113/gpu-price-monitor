package monitor

import (
	"encoding/json"
	"fmt"
	"log"
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

	stockMap := make(StockMap)
	for i := 0; i < pool.JobCount; i++ {
		stockMap.AddResult((<-pool.Results).(PriceCheckResponse))
	}

	email, hasStock := GenerateResultsEmail(stockMap)

	if hasStock {
		if *cfg.SendEmails {
			SendMail(cfg, "GPU Stock Available!", email)
		} else {
			log.Println("skipping email")
		}
		fmt.Fprint(w, email)
	} else {
		writeAndLog(w, "No stock available")
	}

}

package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thomasrea0113/gpu-price-monitor/domain"
)

func getConfigWithOverrides(overrideCfg *domain.Config) (*domain.Config, error) {
	cfg, err := domain.LoadConfig()
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
	var msg domain.RequestMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		writeAndLog(w, "could not read request body: %v", err)
		return
	}

	cfg, err := getConfigWithOverrides(msg.ConfigOverrides)
	if err != nil {
		writeAndLog(w, "couldn't get config: %v", err)
		return
	}

	pool := NewWorkerPool(100, cfg.GetJobs(), scrape)
	pool.Start()

	for i := 0; i < pool.JobCount; i++ {
		result := (<-pool.Results).(domain.PriceCheckResponse)
		// TODO do something meaningful with the results, like generate an email
		fmt.Printf("Result: %v\n\n", result.Models)
	}

	fmt.Fprintf(w, "Okay")
}

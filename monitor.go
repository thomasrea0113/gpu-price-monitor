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
		return cfg, fmt.Errorf("Error loading config: %v", err)
	}

	if overrideCfg != nil {
		if err := domain.Merge(cfg, *overrideCfg); err != nil {
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

	// TODO process each product for each site

	fmt.Fprintf(w, "Okay")
}

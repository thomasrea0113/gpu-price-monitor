package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thomasrea0113/gpu-price-monitor/domain"
)

func MonitorProducts(w http.ResponseWriter, r *http.Request) {
	message := domain.RequestMessage{}

	cfg, err := domain.LoadConfig()
	if err != nil {
		writeAndLog(w, "Error loading config: %v", err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		writeAndLog(w, "could not read request body")
		return
	}

	var overrides []byte
	if message.ConfigOverrides != nil {
		if overrides, err = json.Marshal(message.ConfigOverrides); err != nil {
			writeAndLog(w, "error marshalling config: %v", message.ConfigOverrides)
			return
		}

		overrideCfg := domain.Config{}
		if err = json.NewDecoder(bytes.NewReader(overrides)).Decode(&overrideCfg); err != nil {
			writeAndLog(w, "error applying config overrides")
			return
		}

		if err := domain.Merge(cfg, overrideCfg); err != nil {
			writeAndLog(w, "error merging override config with base config: %v", err)
			return
		}
	}

	fmt.Fprint(w, "Okay")
}

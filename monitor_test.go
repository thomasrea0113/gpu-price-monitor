package monitor_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	monitor "github.com/thomasrea0113/gpu-price-monitor"
)

func TestConfig(t *testing.T) {
	cfg, err := monitor.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	if cfg.Environment != "Test" {
		t.Fatal(cfg.Environment)
	}

	if len(cfg.Sites) == 0 {
		t.Fatal(cfg.Sites)
	}

	if val := cfg.Sites[0].Name; val == "" {
		t.Fatal(val)
	}

	if len(cfg.Products) == 0 {
		t.Fatal(cfg.Products)
	}

	if cfg.SendEmails == nil || *cfg.SendEmails == true {
		t.Fatal(cfg.SendEmails)
	}

	if val := cfg.Sites[0].Name; val == "" {
		t.Fatal(val)
	}

	if val := cfg.Mail.From; val == "" {
		t.Fatal(val)
	}

	if val := cfg.Mail.Password; val == "" {
		t.Fatal(val)
	}
}

func TestMail(t *testing.T) {
	cfg, err := monitor.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	if err := monitor.SendMail(cfg, "Test Message", "<h3>Hello world!</h3><p>some <i>text</i></p>"); err != nil {
		t.Fatalf("failed sending mail: %v", err)
	}
}

func TestMonitorProducts(t *testing.T) {
	// TODO add more tests
	tests := []monitor.RequestMessage{
		{},
		// TODO is there a simpler way to initialize nested structs?
		{ConfigOverrides: &monitor.Config{SendEmails: monitor.NewFalse()}},
	}

	for _, test := range tests {
		var message []byte
		var err error

		if message, err = json.Marshal(test); err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("GET", "/", bytes.NewReader(message))
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		monitor.MonitorProducts(rr, req)

		if got := rr.Body.String(); got == "" {
			t.Fatalf("MonitorProducts(%q) = response was empty!", string(message))
		}
	}
}

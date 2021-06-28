package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Site struct {
	Name      string `json:"name"`
	UrlFormat string `json:"urlFormat"`
}

type Product struct {
	Name               string   `json:"name"`
	AdditionalKeywords []string `json:"additionalKeywords"`
	PriceThreshhold    int      `json:"priceThreshold"`
}

type Mail struct {
	To       []string
	From     string `json:"from"`
	Password string `json:"password"`
}

type Config struct {
	Environment string    `envconfig:"GOLANG_ENVIRONMENT"`
	SendEmails  *bool     `json:"sendEmails"`
	Sites       []Site    `json:"sites"`
	Products    []Product `json:"products"`
	Mail        Mail      `json:"mail"`
}

func loadConfigFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	decoder := json.NewDecoder(f)

	cfg := &Config{}

	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg Config) GetJobs() []interface{} {
	products := make([]interface{}, len(cfg.Sites)*len(cfg.Products))

	i := 0
	for _, site := range cfg.Sites {
		for _, product := range cfg.Products {
			products[i] = PriceCheckJob{Site: site, Product: product}
			i++
		}
	}

	return products
}

// a simple merge function that applies each field value on top of dest, or the value of the previous Config in the array
// zero and nil values are not overwritten
func (dest *Config) Merge(ss ...Config) error {
	// TODO use reflection to cut down on maintaince/boilterplate?
	for _, s := range ss {
		if !reflect.ValueOf(s.Environment).IsZero() {
			dest.Environment = s.Environment
		}
		if !reflect.ValueOf(s.SendEmails).IsNil() {
			dest.SendEmails = s.SendEmails
		}
		if !reflect.ValueOf(s.Products).IsZero() {
			dest.Products = s.Products
		}
		if !reflect.ValueOf(s.Sites).IsZero() {
			dest.Sites = s.Sites
		}
		if !reflect.ValueOf(s.Mail.To).IsZero() {
			dest.Mail.To = s.Mail.To
		}
		if !reflect.ValueOf(s.Mail.From).IsZero() {
			dest.Mail.From = s.Mail.From
		}
		if !reflect.ValueOf(s.Mail.Password).IsZero() {
			dest.Mail.Password = s.Mail.Password
		}
	}

	return nil
}

func (dest *Config) loadAndMerge(path string) error {
	mergingCfg, err := loadConfigFile(path)
	if err != nil {
		return err
	}

	err = dest.Merge(*mergingCfg)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig() (*Config, error) {
	// the destination config that all others will be merged into
	destCfg := &Config{}

	// loading environment variables
	if err := envconfig.Process("", destCfg); err != nil {
		return nil, err
	}

	if reflect.ValueOf(destCfg.Environment).IsZero() {
		return nil, errors.New("GOLANG_ENVIRONMENT variable not set")
	}

	// base config loaded, always required
	if err := destCfg.loadAndMerge("config.json"); err != nil {
		return nil, err
	}

	// base secret config loaded, always required since it has credentials
	if err := destCfg.loadAndMerge("config.secret.json"); err != nil {
		return nil, err
	}

	// environment specific config loaded, not required
	envCfgPath := fmt.Sprintf("config.%v.json", strings.ToLower(destCfg.Environment))
	if _, err := os.Stat(envCfgPath); !os.IsNotExist(err) {
		// the below would NOT work, because stat won't return an error if the file exists. IsExist would be appropriate if
		// you call an os method that creates a new file, but the file you're creating already exists
		// if _, err := os.Stat(envCfgPath); os.IsExist(err) {
		if err = destCfg.loadAndMerge(envCfgPath); err != nil {
			return nil, err
		}
	}

	envSecretCfgPath := fmt.Sprintf("config.%v.secret.json", strings.ToLower(destCfg.Environment))
	if _, err := os.Stat(envSecretCfgPath); !os.IsNotExist(err) {
		if err = destCfg.loadAndMerge(envSecretCfgPath); err != nil {
			return nil, err
		}
	}

	return destCfg, nil
}

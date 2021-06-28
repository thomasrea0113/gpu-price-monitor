package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Environment string `envconfig:"GOLANG_ENVIRONMENT"`
	SendEmails  *bool  `json:"sendEmails"`
	Sites       []struct {
		Name      string `json:"name"`
		UrlFormat string `json:"urlFormat"`
	} `json:"sites"`
	Products []struct {
		Name               string   `json:"name"`
		AdditionalKeywords []string `json:"additionalKeywords"`
		PriceThreshhold    float32  `json:"priceThreshold"`
	} `json:"products"`
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

// a simple merge function that applies each field value on top of dest, or the value of the previous Config in the array
// zero and nil values are not overwritten
func Merge(dest *Config, ss ...Config) error {
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
	}

	return nil
}

func LoadConfig() (*Config, error) {
	env, ok := os.LookupEnv("GOLANG_ENVIRONMENT")
	if !ok {
		return nil, errors.New("GOLANG_ENVIRONMENT variable not set")
	}

	osEnvCfg := &Config{}

	var err error
	var baseCfg, envCfg *Config

	// loading environment variables
	if err = envconfig.Process("", osEnvCfg); err != nil {
		return nil, err
	}

	// base config loaded, always required
	if baseCfg, err = loadConfigFile("config.json"); err != nil {
		return nil, err
	}

	// merge os config with base, both should exist and succeed
	if err = Merge(osEnvCfg, *baseCfg); err != nil {
		return nil, err
	}

	// environment specific config loaded, not required
	envCfgPath := fmt.Sprintf("config.%v.json", env)
	if _, err := os.Stat(envCfgPath); os.IsExist(err) {
		if envCfg, err = loadConfigFile(fmt.Sprintf("config.%v.json", env)); err != nil {
			return nil, err
		}

		if err = Merge(osEnvCfg, *envCfg); err != nil {
			return nil, err
		}
	}

	// all merges would applies to the osEnv config, so we can just return it
	return osEnvCfg, nil
}

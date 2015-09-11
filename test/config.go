package test

import (
	"fmt"
	"os"
)

// Pull in required config from the environment.
// Return a map instead of an api.Config object to avoid import cycles.
func LoadConfig() (map[string]string, error) {
	cfg := map[string]string{}

	baseUrl := os.Getenv("SPARKPOST_BASEURL")
	if baseUrl == "" {
		return nil, fmt.Errorf("Base URL not set in environment: SPARKPOST_BASEURL")
	}
	cfg["baseurl"] = baseUrl

	apiKey := os.Getenv("SPARKPOST_APIKEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API Key not set in environment: SPARKPOST_APIKEY")
	}
	cfg["apikey"] = apiKey

	return cfg, nil
}

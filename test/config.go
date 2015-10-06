package test

import (
	"fmt"
	"os"
	"strconv"
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

	apiKey := os.Getenv("SPARKPOST_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API Key not set in environment: SPARKPOST_API_KEY")
	}
	cfg["apikey"] = apiKey

	apiVer := os.Getenv("SPARKPOST_APIVER")
	if apiVer == "" {
		apiVer = "1"
	}
	if _, err := strconv.Atoi(apiVer); err != nil {
		return nil, fmt.Errorf("API Version must be an integer: SPARKPOST_APIVER")
	}
	cfg["apiver"] = apiVer

	return cfg, nil
}

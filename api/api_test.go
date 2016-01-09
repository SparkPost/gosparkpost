package api

import (
	"testing"

	"github.com/SparkPost/go-sparkpost/test"
)

func TestAPI(t *testing.T) {
	cfgMap, err := test.LoadConfig()
	if err != nil {
		t.Error(err)
	}
	cfg := &Config{
		BaseUrl: cfgMap["baseurl"],
		ApiKey:  cfgMap["apikey"],
	}

	var client Client
	err = client.Init(cfg)
	if err != nil {
		t.Error(err)
	}
}

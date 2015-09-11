package api

import (
	"testing"

	"bitbucket.org/yargevad/go-sparkpost/test"
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

	var a API
	err = a.Init(cfg)
	if err != nil {
		t.Error(err)
	}
}

package api

import (
	"testing"

	"bitbucket.org/yargevad/go-sparkpost/config"
)

func TestAPI(t *testing.T) {
	// FIXME: hardcoded config
	cfg, err := config.Load("../config.json")
	if err != nil {
		t.Error(err)
	}

	_, err = New(cfg)
	if err != nil {
		t.Error(err)
	}
}

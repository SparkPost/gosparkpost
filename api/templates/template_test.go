package templates

import (
	"testing"

	"bitbucket.org/yargevad/go-sparkpost/config"
)

func TestTemplates(t *testing.T) {
	// FIXME: hardcoded config
	cfg, err := config.Load("../../config.json")
	if err != nil {
		t.Error(err)
	}

	var tAPI Templates
	err = tAPI.Init(cfg)
	if err != nil {
		t.Error(err)
	}
}

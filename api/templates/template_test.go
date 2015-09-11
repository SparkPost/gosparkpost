package templates

import (
	"fmt"
	"testing"

	"bitbucket.org/yargevad/go-sparkpost/config"
)

func TestTemplates(t *testing.T) {
	// FIXME: hardcoded config
	cfg, err := config.Load("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}

	T, err = New(cfg)
	if err != nil {
		t.Error(err)
		return
	}

	content := Content{
		Subject: "this is a test template",
		//Text:    "text part of the test template {{a}",
		Text: "text part of the test template",
		From: map[string]string{
			"name":  "test name",
			"email": "test@email.com",
		},
	}

	id, err := T.Create("test template", content)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Created Template with id=%s\n", id)
}

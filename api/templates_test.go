package api_test

import (
	"fmt"
	"testing"

	"github.com/SparkPost/go-sparkpost/api"
	"github.com/SparkPost/go-sparkpost/test"
)

func TestTemplates(t *testing.T) {
	cfgMap, err := test.LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	cfg, err := api.NewConfig(cfgMap)
	if err != nil {
		t.Error(err)
		return
	}

	var client api.Client
	err = client.Init(cfg)
	if err != nil {
		t.Error(err)
		return
	}

	tlist, _, err := client.Templates()
	if err != nil {
		t.Error(err)
		return
	}

	t.Error(fmt.Errorf("%s", tlist))
	return

	content := api.Content{
		Subject: "this is a test template",
		// NB: deliberate syntax error
		//Text: "text part of the test template {{a}",
		Text: "text part of the test template",
		From: map[string]string{
			"name":  "test name",
			"email": "test@email.com",
		},
	}
	template := &api.Template{Content: content, Name: "test template"}

	id, _, err := client.TemplateCreate(template)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Created Template with id=%s\n", id)

	_, err = client.TemplateDelete(id)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Deleted Template with id=%s\n", id)
}

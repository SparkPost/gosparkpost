package templates

import (
	"fmt"
	"testing"

	"bitbucket.org/yargevad/go-sparkpost/api"
	"bitbucket.org/yargevad/go-sparkpost/test"
)

func TestTemplates(t *testing.T) {
	cfgMap, err := test.LoadConfig()
	if err != nil {
		t.Error(err)
	}
	cfg := &api.Config{
		BaseUrl: cfgMap["baseurl"],
		ApiKey:  cfgMap["apikey"],
	}

	Template, err := New(cfg)
	if err != nil {
		t.Error(err)
		return
	}

	tlist, err := Template.List()
	if err != nil {
		t.Error(err)
		return
	}
	t.Error(fmt.Errorf("%s", tlist))
	return

	content := &Content{
		Subject: "this is a test template",
		// NB: deliberate syntax error
		//Text: "text part of the test template {{a}",
		Text: "text part of the test template",
		From: map[string]string{
			"name":  "test name",
			"email": "test@email.com",
		},
	}

	id, err := Template.Create("test template", content)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Created Template with id=%s\n", id)

	err = Template.Delete(id)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Deleted Template with id=%s\n", id)
}

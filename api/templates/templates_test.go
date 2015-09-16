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

	tmpl, err := New(cfg)
	if err != nil {
		t.Error(err)
		return
	}

	//tlist, err := tmpl.List()
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	//t.Error(fmt.Errorf("%s", tlist))
	//return

	templ, err := tmpl.Build(map[string]string{
		"from_email": "a@b.com",
		"from_name":  "a b",
	})
	if err != nil {
		t.Error(err)
		return
	}

	templ.SetHeaders(map[string]string{
		"x-binding": "foo",
	})

	t.Error(fmt.Errorf("%s", templ))
	return

	content := Content{
		Subject: "this is a test template",
		// NB: deliberate syntax error
		//Text: "text part of the test template {{a}",
		Text: "text part of the test template",
		From: map[string]string{
			"name":  "test name",
			"email": "test@email.com",
		},
	}
	template := &Template{Content: content, Name: "test template"}

	id, err := tmpl.Create(template)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Created Template with id=%s\n", id)

	err = tmpl.Delete(id)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Deleted Template with id=%s\n", id)
}

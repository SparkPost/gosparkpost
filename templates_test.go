package gosparkpost_test

import (
	"fmt"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/SparkPost/gosparkpost/test"
)

func TestTemplates(t *testing.T) {
	cfgMap, err := test.LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	cfg, err := sp.NewConfig(cfgMap)
	if err != nil {
		t.Error(err)
		return
	}

	var client sp.Client
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
	t.Logf("templates listed: %+v", tlist)

	content := sp.Content{
		Subject: "this is a test template",
		// NB: deliberate syntax error
		//Text: "text part of the test template {{a}",
		Text: "text part of the test template",
		From: map[string]string{
			"name":  "test name",
			"email": "test@email.com",
		},
	}
	template := &sp.Template{Content: content, Name: "test template"}

	id, _, err := client.TemplateCreate(template)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Created Template with id=%s\n", id)

	d := map[string]interface{}{}
	res, err := client.TemplatePreview(id, &sp.PreviewOptions{d})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Preview Template with id=%s and response %+v\n", id, res)

	_, err = client.TemplateDelete(id)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Deleted Template with id=%s\n", id)
}

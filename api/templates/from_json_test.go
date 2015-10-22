package templates_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SparkPost/go-sparkpost/api"
	"github.com/SparkPost/go-sparkpost/api/templates"
)

// Build a native Go Template structure from a JSON string
func ExampleTemplate() {
	cfg := api.Config{BaseUrl: "https://example.com", ApiKey: "foo"}
	T, err := templates.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	template := &templates.Template{}
	jsonStr := `{
		"name": "testy template",
		"content": {
			"html": "this is a <b>test</b> email!",
			"subject": "test email",
			"from": {
				"name": "tester",
				"email": "tester@example.com"
			},
			"reply_to": "tester@example.com"
		}
	}`
	err = json.Unmarshal([]byte(jsonStr), template)
	if err != nil {
		log.Fatal(err)
	}
}

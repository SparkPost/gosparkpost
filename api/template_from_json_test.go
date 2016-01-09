package api_test

import (
	"encoding/json"
	"log"

	"github.com/SparkPost/go-sparkpost/api"
)

// Build a native Go Template structure from a JSON string
func ExampleTemplate() {
	template := &api.Template{}
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
	err := json.Unmarshal([]byte(jsonStr), template)
	if err != nil {
		log.Fatal(err)
	}
}

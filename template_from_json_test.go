package gosparkpost_test

import (
	"encoding/json"
	"log"

	sp "github.com/SparkPost/gosparkpost"
)

// Build a native Go Template structure from a JSON string
func ExampleTemplate() {
	template := &sp.Template{}
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

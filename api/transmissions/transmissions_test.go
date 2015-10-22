package transmissions

import (
	"testing"

	"github.com/SparkPost/go-sparkpost/api"
	"github.com/SparkPost/go-sparkpost/api/templates"
	"github.com/SparkPost/go-sparkpost/test"
)

func TestTransmissions(t *testing.T) {
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

	TransAPI, err := New(*cfg)
	if err != nil {
		t.Error(err)
		return
	}

	// TODO: 404 from Transmission Create could mean either
	// Recipient List or Content wasn't found - open doc ticket
	// to make error message more specific

	T := &Transmission{
		CampaignID: "msys_smoke",
		ReturnPath: "dgray@messagesystems.com",
		Recipients: []string{"dgray@messagesystems.com", "dgray@sparkpost.com"},
		// Single-recipient Transmissions are transient - Retrieve will 404
		//Recipients: []string{"dgray@messagesystems.com"},
		Content: templates.Content{
			Subject: "this is a test message",
			HTML:    "this is the <b>HTML</b> body of the test message",
			From: map[string]string{
				"name":  "Dave Gray",
				"email": "dgray@messagesystems.com",
			},
		},
		Metadata: map[string]interface{}{
			"binding": "example",
		},
	}
	err = T.Validate()
	if err != nil {
		t.Error(err)
		return
	}

	id, err := TransAPI.Create(T)
	if err != nil {
		t.Error(err)
		return
	}

	t.Errorf("Transmission created with id [%s]", id)

	tr, err := TransAPI.Retrieve(id)
	if err != nil {
		t.Error(err)
		return
	}

	if TransAPI.Response != nil {
		t.Errorf("Retrieve returned HTTP %s\n", TransAPI.Response.HTTP.Status)
		if len(TransAPI.Response.Errors) > 0 {
			for _, e := range TransAPI.Response.Errors {
				json, err := e.Json()
				if err != nil {
					t.Error(err)
				}
				t.Errorf("%s\n", json)
			}
		} else {
			t.Errorf("Transmission retrieved: %s=%s\n", tr.ID, tr.State)
		}
	}

	err = TransAPI.Delete(id)
	if err != nil {
		t.Error(err)
		return
	}

	t.Errorf("Delete returned HTTP %s\n%s\n", TransAPI.Response.HTTP.Status, TransAPI.Response.Body)

}

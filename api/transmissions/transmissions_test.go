package transmissions

import (
	"testing"

	"github.com/SparkPost/go-sparkpost/api"
	"github.com/SparkPost/go-sparkpost/api/recipient_lists"
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
		Recipients: []recipient_lists.Recipient{
			{Address: map[string]interface{}{
				"email": "dgray@messagesystems.com",
			}},
		},
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
}

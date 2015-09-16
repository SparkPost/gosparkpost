package recipient_lists

import (
	"testing"

	"github.com/SparkPost/go-sparkpost/api"
	"github.com/SparkPost/go-sparkpost/test"
)

func TestRecipients(t *testing.T) {
	cfgMap, err := test.LoadConfig()
	if err != nil {
		t.Error(err)
	}
	cfg := &api.Config{
		BaseUrl: cfgMap["baseurl"],
		ApiKey:  cfgMap["apikey"],
	}

	//Recipients, err := New(cfg)
	_, err = New(cfg)
	if err != nil {
		t.Error(err)
		return
	}
}

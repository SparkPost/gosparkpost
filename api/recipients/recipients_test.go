package recipients

import (
	"testing"

	"bitbucket.org/yargevad/go-sparkpost/api"
	"bitbucket.org/yargevad/go-sparkpost/test"
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

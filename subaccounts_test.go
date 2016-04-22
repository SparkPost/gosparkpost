package gosparkpost_test

import (
	"fmt"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/SparkPost/gosparkpost/test"
)

func TestSubaccounts(t *testing.T) {
	if true {
		// Temporarily disable test so TravisCI reports build success instead of test failure.
		// NOTE: need travis to set sparkpost base urls etc, or mock http request
		return
	}

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

	tlist, _, err := client.Subaccounts()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("subaccounts listed: %+v", tlist)

	return
}

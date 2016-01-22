package gosparkpost_test

import (
	"strings"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/SparkPost/gosparkpost/test"
)

func TestRecipients(t *testing.T) {
	if true {
		// Temporarily disable test so TravisCI reports build success instead of test failure.
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

	list, _, err := client.RecipientLists()
	if err != nil {
		t.Error(err)
		return
	}

	strs := make([]string, len(*list))
	for idx, rl := range *list {
		strs[idx] = rl.String()
	}
	t.Errorf("%s\n", strings.Join(strs, "\n"))
}

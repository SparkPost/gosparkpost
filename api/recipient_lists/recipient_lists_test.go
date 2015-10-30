package recipient_lists

import (
	"strings"
	"testing"

	"github.com/SparkPost/go-sparkpost/api"
	"github.com/SparkPost/go-sparkpost/test"
)

func TestRecipients(t *testing.T) {
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

	RLAPI, err := New(*cfg)
	if err != nil {
		t.Error(err)
		return
	}

	/*
		R1, err := RLAPI.BuildRecipient(map[string]interface{}{
			"return_path": "a@b.com",
			"email":       "abc@example.com",
			"name":        "A B C",
			"tags":        []string{"abc", "def", "ghi"},
			"metadata": map[string]interface{}{
				"abc": "def",
				"ghi": []interface{}{"j", "k", "l"},
				"mno": map[string]interface{}{
					"p": []interface{}{"q": "r"},
					"s": "t",
				},
			},
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Error(fmt.Errorf("%s", R1))
	*/

	list, res, err := RLAPI.List()
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

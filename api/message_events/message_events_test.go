package message_events

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/SparkPost/go-sparkpost/api"
	"github.com/SparkPost/go-sparkpost/events"
	"github.com/SparkPost/go-sparkpost/test"
)

func TestMessageEvents(t *testing.T) {
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

	mev, err := New(*cfg)
	if err != nil {
		t.Error(err)
		return
	}

	types := []string{"open", "click", "bounce"}
	e, err := mev.Samples(&types)
	if err != nil {
		t.Error(err)
		return
	}

	for _, ev := range *e {
		switch event := ev.(type) {
		case *events.Open, *events.Click:
			t.Error(fmt.Errorf("%s", event))

		case *events.Bounce:
			t.Error(fmt.Errorf("%s", event.ECLog()))

		default:
			t.Errorf("Unsupported type [%s]", reflect.TypeOf(ev))
		}
	}
}

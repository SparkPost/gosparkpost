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

	types := []string{"open"}
	e, err := mev.Samples(&types)
	if err != nil {
		t.Error(err)
		return
	}

	for _, ev := range *e {
		switch event := ev.(type) {
		case *events.Open:
			g := event.GeoIP
			if g != nil {
				t.Error(fmt.Errorf("%s (%s, %s)", e, g.Latitude, g.Longitude))
			} else {
				t.Error(fmt.Errorf("%s", e))
			}
		default:
			t.Errorf("Unsupported type [%s]", reflect.TypeOf(ev))
		}
	}
}

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

	//types := []string{"open", "click", "bounce"}
	//e, err := mev.Samples(&types)
	e, err := mev.Samples(nil)
	if err != nil {
		t.Error(err)
		return
	}

	for _, ev := range *e {
		//t.Error(fmt.Errorf("%s", ev))
		switch event := ev.(type) {
		case *events.Click, *events.Open, *events.GenerationFailure, *events.GenerationRejection,
			*events.ListUnsubscribe, *events.LinkUnsubscribe, *events.PolicyRejection,
			*events.RelayInjection, *events.RelayRejection, *events.RelayDelivery,
			*events.RelayTempfail, *events.RelayPermfail, *events.SpamComplaint:
			t.Error(fmt.Errorf("%s", event))

		case *events.Bounce, *events.Delay, *events.Delivery, *events.Injection, *events.OutOfBand:
			t.Error(fmt.Errorf("%s", events.ECLog(event)))

		default:
			t.Errorf("Unsupported type [%s]", reflect.TypeOf(ev))
		}
	}
}

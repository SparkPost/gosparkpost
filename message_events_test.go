package gosparkpost_test

import (
	"fmt"
	"reflect"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/SparkPost/gosparkpost/events"
	"github.com/SparkPost/gosparkpost/test"
)

func TestMessageAllEvents(t *testing.T) {
	// if true {
	// 	// Temporarily disable test so TravisCI reports build success instead of test failure.
	// 	return
	// }

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

	e, err := client.EventSamples(nil)
	if err != nil {
		t.Error(err)
		return
	}

	if len(*e) == 0 {
		t.Error("expected non-empty result")
	}

	for _, ev := range *e {
		t.Error(fmt.Errorf("%s", ev))
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

func TestMessageFilteredEvents(t *testing.T) {
	// if true {
	// 	// Temporarily disable test so TravisCI reports build success instead of test failure.
	// 	return
	// }

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

	types := []string{"open", "click", "bounce"}
	e, err := client.EventSamples(&types)
	if err != nil {
		t.Error(err)
		return
	}

	if len(*e) == 0 {
		t.Error("expected non-empty result")
	}

	for _, ev := range *e {
		t.Error(fmt.Errorf("%s", ev))
		switch event := ev.(type) {
		case *events.Click, *events.Open, *events.Bounce:
			// Expected, ok.
		default:
			t.Errorf("Unexpected type %T, should have been filtered out.", event)
		}
	}
}

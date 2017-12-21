package events

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGeoIP(t *testing.T) {
	for idx, test := range []struct {
		in  []byte
		err error
	}{
		{[]byte(`{"country":"USA","region":"CO","city":"Denver","latitude":39.7392,"longitude":104.9903}`), nil},
	} {
		geo := GeoIP{}
		err := json.Unmarshal(test.in, &geo)
		if err != nil {
			t.Fatal(err)
		}

		jsonBytes, err := json.Marshal(geo)
		if err != nil {
			if test.err != nil && test.err == err {
				// ignore expected errors
			} else {
				t.Fatal(err)
			}
		}
		if bytes.Compare(jsonBytes, test.in) != 0 {
			t.Errorf("Marshal[%d] => got/want:\n%s\n%s", idx, string(jsonBytes), string(test.in))
		}
	}
}

func TestSampleEvents(t *testing.T) {
	payload, err := ioutil.ReadFile("test/json/sample-events.json")
	if err != nil {
		t.Fatal(err)
	}

	var events Events
	err = json.Unmarshal(payload, &events)
	if err != nil {
		t.Fatal(err)
	}

	for _, event := range events {
		if unknown, ok := event.(*Unknown); ok {
			t.Fatal(unknown)
		}
	}
}

func TestSampleWebhookValidationRequest(t *testing.T) {
	payload, err := ioutil.ReadFile("test/json/sample-webhook-validation.json")
	if err != nil {
		t.Fatal(err)
	}

	var events Events
	err = json.Unmarshal(payload, &events)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected zero events, got %d: %v", len(events), events)
	}
}

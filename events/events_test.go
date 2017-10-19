package events

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGeoIPMarshalling(t *testing.T) {
	for idx, test := range []struct {
		in   GeoIP
		err  error
		json []byte
	}{
		{
			GeoIP{"USA", "CO", "Denver", 39.7392, 104.9903}, nil,
			[]byte(`{"country":"USA","region":"CO","city":"Denver","latitude":39.7392,"longitude":104.9903}`),
		},
	} {
		jsonBytes, err := json.Marshal(test.in)
		if err != nil {
			if test.err != nil && test.err == err {
				// ignore expected errors
			} else {
				t.Fatal(err)
			}
		}
		if bytes.Compare(jsonBytes, test.json) != 0 {
			t.Errorf("Marshal[%d] => got/want:\n%s\n%s", idx, string(jsonBytes), string(test.json))
		}
		geo := GeoIP{}
		err = json.Unmarshal(jsonBytes, &geo)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(test.in, geo) {
			t.Errorf("Unmarshal[%d] => got/want:\n%q\n%q", idx, geo, test.in)
		}
	}
}

func TestSampleEvents(t *testing.T) {
	file, err := os.Open("sample-events.json")
	if err != nil {
		t.Fatal(err)
	}

	payload, err := ioutil.ReadAll(file)
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
	file, err := os.Open("sample-webhook-validation.json")
	if err != nil {
		t.Fatal(err)
	}

	payload, err := ioutil.ReadAll(file)
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

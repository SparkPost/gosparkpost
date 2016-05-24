package events

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

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

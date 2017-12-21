// Package events defines a struct for each type of event and provides various other helper functions.
package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Event is a generic event.
type Event interface {
	EventType() string
}

// Events is a list of generic events. Useful for decoding events from API webhooks.
type Events []Event

var (
	ErrNotImplemented = errors.New("not implemented")
)

// ValidEventType returns true if the event name parameter is valid.
func ValidEventType(eventType string) bool {
	if _, ok := EventForName(eventType).(*Unknown); ok {
		return false
	}
	return true
}

// EventForName returns a struct matching the passed-in type.
func EventForName(eventType string) Event {
	switch eventType {
	case "bounce":
		return &Bounce{}
	case "click":
		return &Click{}
	case "creation":
		return &Creation{}
	case "delay":
		return &Delay{}
	case "delivery":
		return &Delivery{}
	case "generation_failure":
		return &GenerationFailure{}
	case "generation_rejection":
		return &GenerationRejection{}
	case "injection":
		return &Injection{}
	case "list_unsubscribe":
		return &ListUnsubscribe{}
	case "link_unsubscribe":
		return &LinkUnsubscribe{}
	case "open":
		return &Open{}
	case "out_of_band":
		return &OutOfBand{}
	case "policy_rejection":
		return &PolicyRejection{}
	case "spam_complaint":
		return &SpamComplaint{}
	case "relay_delivery":
		return &RelayDelivery{}
	case "relay_injection":
		return &RelayInjection{}
	case "relay_message":
		return &RelayMessage{}
	case "relay_permfail":
		return &RelayPermfail{}
	case "relay_rejection":
		return &RelayRejection{}
	case "relay_tempfail":
		return &RelayTempfail{}
	case "sms_status":
		return &SMSStatus{}
	}
	return &Unknown{}
}

func ParseRawJSONEvents(rawEvents []json.RawMessage) ([]Event, error) {
	events := []Event{}

	// Each item is event data in raw JSON.
	for _, rawEvent := range rawEvents {
		var typeLookup EventCommon
		if err := json.Unmarshal(rawEvent, &typeLookup); err != nil {
			typeLookup.Type = "unknown"
		}

		event := EventForName(typeLookup.EventType())
		if e, ok := event.(*Unknown); ok {
			e.EventCommon.Type = typeLookup.EventType()
			e.RawJSON = rawEvent
			e.Error = ErrNotImplemented
			events = append(events, e)
			continue
		}

		// Unmarshal into specic event object.
		if err := json.Unmarshal(rawEvent, &event); err != nil {
			event = &Unknown{
				EventCommon: EventCommon{Type: typeLookup.EventType()},
				RawJSON:     rawEvent,
				Error:       err,
			}
		}
		events = append(events, event)
	}

	return events, nil
}

func (events *Events) UnmarshalJSON(data []byte) error {
	// Parse raw events from Event Webhook ("msys"-wrapped array of events).
	rawEvents, err := parseRawJSONEventsFromWebhook(data)
	if err != nil {
		// Parse raw events from Event Samples ("results" object with array of events).
		rawEvents, err = parseRawJSONEventsFromSamples(data)
		if err != nil {
			return err
		}
	}

	*events, err = ParseRawJSONEvents(rawEvents)
	if err != nil {
		return err
	}
	return nil
}

func parseRawJSONEventsFromWebhook(data []byte) ([]json.RawMessage, error) {
	var rawEvents []json.RawMessage

	// These "msys"-wrapped events are being sent on Webhooks.
	var msysEventWrappers []struct {
		MsysEventWrapper map[string]json.RawMessage `json:"msys"`
	}
	if err := json.Unmarshal(data, &msysEventWrappers); err != nil {
		return nil, err
	}

	for _, wrapper := range msysEventWrappers {
		for _, rawEvent := range wrapper.MsysEventWrapper {
			rawEvents = append(rawEvents, rawEvent)
		}
	}

	return rawEvents, nil
}

func parseRawJSONEventsFromSamples(data []byte) ([]json.RawMessage, error) {
	// Object with array of events is being sent on Events Samples.
	var resultsWrapper struct {
		RawEvents []json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(data, &resultsWrapper); err != nil {
		return nil, err
	}

	return resultsWrapper.RawEvents, nil
}

func ECLog(e Event) string {
	// XXX: this feels like the wrong way; can't figure out the right way
	switch e.(type) {
	case *Bounce:
		b := e.(*Bounce)
		return b.ECLog()
	case *Delay:
		d := e.(*Delay)
		return d.ECLog()
	case *Delivery:
		d := e.(*Delivery)
		return d.ECLog()
	case *Injection:
		i := e.(*Injection)
		return i.ECLog()
	case *OutOfBand:
		o := e.(*OutOfBand)
		return o.ECLog()
	}
	return ""
}

type ECLogger interface {
	ECLog() string
}

// EventCommon contains fields common to all types of Event objects
type EventCommon struct {
	Type string `json:"type"`
}

func (e EventCommon) EventType() string { return e.Type }

type Unknown struct {
	EventCommon
	RawJSON json.RawMessage
	Error   error
}

func (e *Unknown) EventType() string { return "unknown" }

func (e *Unknown) String() string {
	return fmt.Sprintf("Unknown event (type %q): %v\n%s", e.EventCommon.EventType(), e.Error, e.RawJSON)
}

func (e *Unknown) UnmarshalJSON(data []byte) error {
	return nil
}

type Timestamp time.Time

func (t Timestamp) String() string {
	return time.Time(t).String()
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(time.Time(*t).Unix())), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Trim quotes.
	data = bytes.Trim(data, `"`)

	// Timestamps coming from Webhook Events are Unix timestamps.
	unix, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		*t = Timestamp(time.Unix(unix, 0))
		return nil
	}

	// Timestamps coming from Event Samples are in this RFC 3339-like format.
	customTime, err := time.Parse("2006-01-02T15:04:05.000-07:00", string(data))
	if err != nil {
		return err
	}

	*t = Timestamp(customTime)
	return nil
}

type GeoIP struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  LatLong `json:"latitude"`
	Longitude LatLong `json:"longitude"`
}

// The API inconsistently returns float or string. We need a custom unmarshaller.
type LatLong float32

func (v *LatLong) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", *v)), nil
}

func (v *LatLong) UnmarshalJSON(data []byte) error {
	// Trim quotes if the API returns string.
	data = bytes.Trim(data, `"`)

	// Parse the actual value.
	value, err := strconv.ParseFloat(string(data), 32)
	if err != nil {
		return err
	}

	*v = LatLong(value)
	return nil
}

type Creation struct {
	EventCommon
	Accepted        string      `json:"accepted_rcpts"`
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	InjectionMethod string      `json:"inj_method"`
	NodeName        string      `json:"node_name"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            []string    `json:"rcpt_tags"`
	Submitted       string      `json:"submitted_rcpts"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       Timestamp   `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserID          string      `json:"user_id"`
}

func (c *Creation) String() string {
	return fmt.Sprintf("%s CT %s (%s, %s)",
		c.Timestamp, c.TransmissionID, c.Submitted, c.Accepted)
}

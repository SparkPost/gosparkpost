// Package events defines a struct for each type of event and provides various other helper functions.
package events

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// eventTypes contains all of the valid event types
var eventTypes = map[string]bool{
	"bounce":               true,
	"click":                true,
	"creation":             false,
	"delay":                true,
	"delivery":             true,
	"generation_failure":   true,
	"generation_rejection": true,
	"injection":            true,
	"list_unsubscribe":     true,
	"link_unsubscribe":     true,
	"open":                 true,
	"out_of_band":          true,
	"policy_rejection":     true,
	"spam_complaint":       true,
	"relay_delivery":       true,
	"relay_injection":      true,
	"relay_message":        true,
	"relay_permfail":       true,
	"relay_rejection":      true,
	"relay_tempfail":       true,
}

// ValidEventType returns true if the event name parameter is valid.
func ValidEventType(eventType string) bool {
	if _, ok := eventTypes[eventType]; ok {
		return true
	}
	return false
}

// EventForName returns a struct matching the passed-in type.
func EventForName(eventType string) Event {
	if !ValidEventType(eventType) {
		return nil
	}

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

	default:
		return nil
	}
}

// Event allows 2+ types of event in a data structure.
type Event interface {
	EventType() string
}

// Events is a list of events. Useful for decoding.
type Events []Event

func (events Events) UnmarshalJSON(data []byte) error {
	events = []Event{}

	var eventWrappers []struct {
		MsysEventWrapper map[string]json.RawMessage `json:"msys"`
	}
	if err := json.Unmarshal(data, &eventWrappers); err != nil {
		return err
	}

	for i, wrapper := range eventWrappers {
		for _, eventData := range wrapper.MsysEventWrapper {
			var typeLookup EventCommon
			if err := json.Unmarshal(eventData, &typeLookup); err != nil {
				log.Printf("lookup failed: %v %v", eventData)
				return err
			}
			log.Printf("lookup: %v %v", typeLookup.EventType(), typeLookup)

			var event Event
			switch typeLookup.EventType() {
			case "bounce":
				event = &Bounce{}
			case "click":
				event = &Click{}
			case "creation":
				event = &Creation{}
			case "delay":
				event = &Delay{}
			case "delivery":
				event = &Delivery{}
			case "generation_failure":
				event = &GenerationFailure{}
			case "generation_rejection":
				event = &GenerationRejection{}
			case "injection":
				event = &Injection{}
			case "list_unsubscribe":
				event = &ListUnsubscribe{}
			case "link_unsubscribe":
				event = &LinkUnsubscribe{}
			case "open":
				event = &Open{}
			case "out_of_band":
				event = &OutOfBand{}
			case "policy_rejection":
				event = &PolicyRejection{}
			case "spam_complaint":
				event = &SpamComplaint{}
			case "relay_delivery":
				event = &RelayDelivery{}
			case "relay_injection":
				event = &RelayInjection{}
			case "relay_message":
				event = &RelayMessage{}
			case "relay_permfail":
				event = &RelayPermfail{}
			case "relay_rejection":
				event = &RelayRejection{}
			case "relay_tempfail":
				event = &RelayTempfail{}
			default:
				event = &Unknown{RawJSON: eventData}
			}
			if err := json.Unmarshal(eventData, &event); err != nil {
				event = &Unknown{RawJSON: eventData, Error: err}
				log.Printf("cannot parse into %T: %v: %s", event, err, eventData)
			}
			log.Printf("item[%v]: %T %v", i, event, event.EventType())
			events = append(events, event)
		}
	}

	return nil
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
	RawJSON json.RawMessage
	Error   error
}

func (e *Unknown) EventType() string { return "unknown" }

func (e *Unknown) String() string {
	return fmt.Sprintf("Unknown event: %v %v", e.Error, e.RawJSON)
}

func (e *Unknown) UnmarshalJSON(data []byte) error {
	return nil
}

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(time.Time(*t).Unix())), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	unix, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Timestamp(time.Unix(unix, 0))
	return nil
}

type GeoIP struct {
	Country   string `json:"country"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserID          string      `json:"user_id"`
}

func (c *Creation) String() string {
	return fmt.Sprintf("%s CT %s (%s, %s)",
		c.Timestamp, c.TransmissionID, c.Submitted, c.Accepted)
}

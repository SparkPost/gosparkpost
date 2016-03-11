// Package events defines a struct for each type of event and provides various other helper functions.
package events

import (
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

func (events *Events) UnmarshalJSON(data []byte) error {
	*events = []Event{}

	var eventWrappers []struct {
		MsysEventWrapper map[string]json.RawMessage `json:"msys"`
	}
	if err := json.Unmarshal(data, &eventWrappers); err != nil {
		return err
	}

	for _, wrapper := range eventWrappers {
		for _, eventData := range wrapper.MsysEventWrapper {
			var typeLookup EventCommon
			if err := json.Unmarshal(eventData, &typeLookup); err != nil {
				return err
			}

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
				event = &Unknown{
					EventCommon: EventCommon{Type: typeLookup.EventType()},
					RawJSON:     eventData,
					Error:       errors.New("not implemented"),
				}
				*events = append(*events, event)
				continue
			}
			if err := json.Unmarshal(eventData, &event); err != nil {
				event = &Unknown{
					EventCommon: EventCommon{Type: typeLookup.EventType()},
					RawJSON:     eventData,
					Error:       err,
				}
			}
			*events = append(*events, event)
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

package events

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/goware/lg"
)

func TestEvents(t *testing.T) {
	file, err := os.Open("events.json")
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
		switch e := event.(type) {
		case *Bounce:
			lg.Debugf("%v to %v was bounced at %v because of %v (%v)", e.TransmissionID, e.Recipient, e.Timestamp, e.Reason, e.RawReason)
		case *Delivery:
			lg.Debugf("%v to %v was delivered at %v", e.TransmissionID, e.Recipient, e.Timestamp)
		case *Injection:
			lg.Debugf("%v to %v was injected at %v", e.TransmissionID, e.Recipient, e.Timestamp)
		case *SpamComplaint:
			lg.Debugf("%v to %v was marked spam by %v (reported to %v) at %v because of %v", e.TransmissionID, e.Recipient, e.ReportedBy, e.ReportedTo, e.Timestamp)
		case *Click:
			lg.Debugf("%v to %v - clicked using %v at %v", e.Recipient, e.TransmissionID, e.GeoIP, e.UserAgent, e.Timestamp)
		case *Open:
			lg.Debugf("%v to %v - opened using %v at %v", e.Recipient, e.TransmissionID, e.GeoIP, e.UserAgent, e.Timestamp)
		case *ListUnsubscribe:
			lg.Debugf("%v to %v - list unsubscribed at %v", e.TransmissionID, e.Recipient, e.Timestamp)
		case *LinkUnsubscribe:
			lg.Debugf("%v to %v - link unsubscribed at %v", e.TransmissionID, e.Recipient, e.Timestamp)
		case *Unknown:
			lg.Debugf("wooohoo %v", e)
		default:
			lg.Errorf("Can't parse SparkPost event of type: %T: not implemented", e)
		}
	}
}

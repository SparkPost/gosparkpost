package gosparkpost_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
)

const eventDocumentationFile = "test/event-docs.json"

var eventDocumentationBytes []byte
var eventGroups = []string{"track", "gen", "unsubscribe", "relay", "message"}
var eventGroupMap = map[string]map[string]int{
	"track_event": {
		"click": 22,
		"open":  20,
	},
	"gen_event": {
		"generation_failure":   21,
		"generation_rejection": 23,
	},
	"unsubscribe_event": {
		"list_unsubscribe": 18,
		"link_unsubscribe": 19,
	},
	"relay_event": {
		"relay_permfail":  15,
		"relay_injection": 12,
		"relay_rejection": 14,
		"relay_delivery":  12,
		"relay_tempfail":  15,
	},
	"message_event": {
		"spam_complaint":   24,
		"out_of_band":      21,
		"policy_rejection": 23,
		"delay":            38,
		"bounce":           37,
		"delivery":         36,
		"injection":        31,
		"sms_status":       22,
	},
}

func init() {
	var err error
	eventDocumentationBytes, err = ioutil.ReadFile(eventDocumentationFile)
	if err != nil {
		panic(err)
	}
}

func TestEventDocs_Get_parse(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	path := fmt.Sprintf(sp.EventDocumentationFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write(eventDocumentationBytes)
	})

	// hit our local handler
	groups, res, err := testClient.EventDocumentation()
	if err != nil {
		t.Errorf("EventDocumentation GET returned error: %v", err)
		for _, e := range res.Verbose {
			t.Error(e)
		}
		return
	}

	// basic content test
	if len(groups) == 0 {
		testFailVerbose(t, res, "EventDocumentation GET returned 0 EventGroups")
	} else {
		// check the top level event data
		eventGroupsSeen := make(map[string]bool, len(eventGroups))
		for _, etype := range eventGroups {
			eventGroupsSeen[etype+"_event"] = false
		}

		for gname, v := range groups {
			eventGroupsSeen[gname] = true
			if _, ok := eventGroupMap[gname]; !ok {
				t.Fatalf("expected group [%s] not present in response", gname)
			}
			for ename, efields := range v.Events {
				if fieldCount, ok := eventGroupMap[gname][ename]; !ok {
					t.Fatalf("expected event [%s] not present in [%s]", ename, gname)
					if fieldCount != len(efields.Fields) {
						t.Fatalf("saw %d fields for %s, expected %d", len(efields.Fields), ename, fieldCount)
					}
				}
			}
		}

		for gname, seen := range eventGroupsSeen {
			if !seen {
				t.Fatalf("expected message type [%s] not returned", gname)
			}
		}
	}
}

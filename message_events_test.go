package gosparkpost_test

import (
	"encoding/json"
	"reflect"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/SparkPost/gosparkpost/events"
	jchk "github.com/juju/testing/checkers"
	"github.com/pkg/errors"
)

type EventsPageResult struct {
	err    error
	status int
	json   string
	out    *sp.EventsPage
}

func TestMessageEventsSearch(t *testing.T) {
	var err error
	var page *sp.EventsPage

	//var res200 = loadTestFile(t, "test/json/message-events_search_200.json")
	var res200_1 = loadTestFile(t, "test/json/message-events_search_200-1.json")
	var res200_2 = loadTestFile(t, "test/json/message-events_search_200-2.json")
	var res200_3 = loadTestFile(t, "test/json/message-events_search_200-3.json")

	var ts1, ts2, ts3 events.Timestamp
	json.Unmarshal([]byte(`"2017-04-26T21:37:32.000+00:00"`), &ts1)
	json.Unmarshal([]byte(`"2017-04-26T21:37:17.000+00:00"`), &ts2)
	json.Unmarshal([]byte(`"2017-04-26T21:37:17.000+00:00"`), &ts3)

	// Each test can return multiple pages of results
	for idx, outer := range []struct {
		in      *sp.EventsPage
		results []EventsPageResult
	}{
		{nil, []EventsPageResult{
			{errors.New("MessageEventsSearch called with nil EventsPage!"), 400, `{}`, nil},
		}},

		//{&sp.EventsPage{
		//	Params: map[string]string{"from": "1970-01-01T00:00"}},
		//	[]EventsPageResult{{nil, 200, res200, nil}},
		//},

		{&sp.EventsPage{
			Params: map[string]string{
				"from":     "1970-01-01T00:00",
				"per_page": "1",
			}},
			[]EventsPageResult{
				{nil, 200, res200_1, &sp.EventsPage{
					TotalCount: 3,
					Events: events.Events{&events.Click{
						EventCommon:     events.EventCommon{Type: "click"},
						CampaignID:      "",
						CustomerID:      "42",
						DeliveryMethod:  "esmtp",
						GeoIP:           &events.GeoIP{Country: "US", Region: "NY", City: "Bronx", Latitude: 40.8499, Longitude: -73.8769},
						IPAddress:       "66.102.8.2",
						MessageID:       "0001441ead58afdeb21d",
						Metadata:        map[string]interface{}{},
						Tags:            []string{},
						Recipient:       "developers@sparkpost.com",
						RecipientType:   "",
						TargetLinkName:  "",
						TargetLinkURL:   "https://sparkpost.com",
						TemplateID:      "best-template-ever",
						TemplateVersion: "10",
						Timestamp:       ts1,
						TransmissionID:  "30494147183576458",
						UserAgent:       "lynx",
					}},
					NextPage: "/api/v1/message-events?page=2&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:09&timezone=UTC",
					LastPage: "/api/v1/message-events?page=3&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:09&timezone=UTC",
				}},

				{nil, 200, res200_2, &sp.EventsPage{
					TotalCount: 3,
					Events: events.Events{&events.Open{
						EventCommon:     events.EventCommon{Type: "open"},
						CampaignID:      "",
						CustomerID:      "42",
						DeliveryMethod:  "esmtp",
						GeoIP:           &events.GeoIP{Country: "US", Region: "NY", City: "Bronx", Latitude: 40.8499, Longitude: -73.8769},
						IPAddress:       "66.102.8.28",
						MessageID:       "0001441ead58afdeb21d",
						Metadata:        map[string]interface{}{},
						Tags:            []string{},
						Recipient:       "developers@sparkpost.com",
						RecipientType:   "",
						TemplateID:      "best-template-ever",
						TemplateVersion: "10",
						Timestamp:       ts2,
						TransmissionID:  "30494147183576458",
						UserAgent:       "Mozilla/5.0 (Windows NT 5.1; rv:11.0) Gecko Firefox/11.0 (via ggpht.com GoogleImageProxy)"},
					},
					NextPage:  "/api/v1/message-events?page=3&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:12&timezone=UTC",
					PrevPage:  "/api/v1/message-events?page=1&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:12&timezone=UTC",
					FirstPage: "/api/v1/message-events?page=1&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:12&timezone=UTC",
					LastPage:  "/api/v1/message-events?page=3&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:12&timezone=UTC",
				}},

				{nil, 200, res200_3, &sp.EventsPage{
					TotalCount: 3,
					Events: events.Events{&events.Open{
						EventCommon:     events.EventCommon{Type: "open"},
						CampaignID:      "",
						CustomerID:      "42",
						DeliveryMethod:  "esmtp",
						GeoIP:           &events.GeoIP{Country: "US", Region: "NY", City: "Bronx", Latitude: 40.8499, Longitude: -73.8769},
						IPAddress:       "66.102.8.2",
						MessageID:       "00042a25ad58fc0145e1",
						Metadata:        map[string]interface{}{},
						Tags:            []string{},
						Recipient:       "sales@sparkpost.com",
						RecipientType:   "",
						TemplateID:      "best-template-ever",
						TemplateVersion: "10",
						Timestamp:       ts2,
						TransmissionID:  "84537445295705713",
						UserAgent:       "Mozilla/5.0 (Windows NT 5.1; rv:11.0) Gecko Firefox/11.0 (via ggpht.com GoogleImageProxy)"},
					},
					PrevPage:  "/api/v1/message-events?page=2&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:13&timezone=UTC",
					FirstPage: "/api/v1/message-events?page=1&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:13&timezone=UTC",
				}},
			},
		},
	} {
		for j, test := range outer.results {
			// Set up a new test in our inner loop since re-registering a handler will panic,
			// and the content we're returning needs to vary.
			testSetup(t)
			mockRestResponseBuilderFormat(t, "GET", test.status, sp.MessageEventsPathFormat, test.json)

			if page == nil {
				_, err = testClient.MessageEventsSearch(outer.in)
				page = outer.in
			} else {
				page, _, err = page.Next()
			}

			if err == nil && test.err != nil || err != nil && test.err == nil {
				t.Errorf("MessageEventsSearch[%d.%d] => err %#v want %#v", idx, j, err, test.err)
			} else if err != nil && err.Error() != test.err.Error() {
				t.Errorf("MessageEventsSearch[%d.%d] => err %#v want %#v", idx, j, err, test.err)
			} else if test.out != nil {
				var cmp *sp.EventsPage
				if page != nil {
					cmp = &sp.EventsPage{
						TotalCount: page.TotalCount,
						NextPage:   page.NextPage,
						PrevPage:   page.PrevPage,
						FirstPage:  page.FirstPage,
						LastPage:   page.LastPage,
					}
					cmp.Events = page.Events

					if eq, err := jchk.DeepEqual(cmp, test.out); !eq {
						t.Errorf("MessageEventsSearch[%d.%d] => %q\n", idx, j, err)
						//t.Errorf("%#v\n", cmp.Events[0].(*events.Click))
					}
				} else {
					t.Errorf("MessageEventsSearch[%d.%d] => page is nil!", idx, j)
				}
			}
			testTeardown()
		}
	}
}

func TestEventSamples(t *testing.T) {
	for idx, test := range []struct {
		in     []string
		err    error
		status int
		json   string
		out    *events.Events
	}{
		{nil, nil, 200, `{}`, nil},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "GET", test.status, sp.MessageEventsSamplesPathFormat, test.json)

		events, _, err := testClient.EventSamples(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("EventSamples[%d] => err %#v want %#v", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("EventSamples[%d.%d] => err %#v want %#v", idx, err, test.err)
		} else if test.out != nil {
			if !reflect.DeepEqual(events, test.out) {
				t.Errorf("EventSamples[%d.%d] => events got/want:\n%#v\n%#v", idx, events, test.out)
			}
		}
	}
}

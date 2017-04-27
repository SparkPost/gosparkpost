package gosparkpost_test

import (
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

	var res200 = loadTestFile(t, "test/json/message-events_search_200.json")
	msgEventsPage1.json = loadTestFile(t, "test/json/message-events_search_200-1.json")
	msgEventsPage2.json = loadTestFile(t, "test/json/message-events_search_200-2.json")
	msgEventsPage3.json = loadTestFile(t, "test/json/message-events_search_200-3.json")

	// Each test can return multiple pages of results
	for idx, outer := range []struct {
		in      *sp.EventsPage
		results []EventsPageResult
	}{
		{nil, []EventsPageResult{
			{errors.New("MessageEventsSearch called with nil EventsPage!"), 400, `{}`, nil},
		}},

		{&sp.EventsPage{
			Params: map[string]string{"from": "1970-01-01T00:00"}},
			[]EventsPageResult{{nil, 200, res200, &sp.EventsPage{}}},
		},

		{&sp.EventsPage{
			Params: map[string]string{
				"from":     "1970-01-01T00:00",
				"per_page": "1",
			}},
			[]EventsPageResult{msgEventsPage1, msgEventsPage2, msgEventsPage3},
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

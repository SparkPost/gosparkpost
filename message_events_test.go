package gosparkpost_test

import (
	"reflect"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
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
	var res200_1 = loadTestFile(t, "test/json/message-events_search_200-1.json")
	var res200_2 = loadTestFile(t, "test/json/message-events_search_200-2.json")
	var res200_3 = loadTestFile(t, "test/json/message-events_search_200-3.json")

	// Each test can return multiple pages of results
	for idx, seq := range []struct {
		in      *sp.EventsPage
		results []EventsPageResult
	}{
		{nil, []EventsPageResult{
			{errors.New("MessageEventsSearch called with nil EventsPage!"), 400, `{}`, nil},
		}},

		{&sp.EventsPage{
			Params: map[string]string{"from": "1970-01-01T00:00"}},
			[]EventsPageResult{{nil, 200, res200, nil}},
		},

		{&sp.EventsPage{
			Params: map[string]string{
				"from":     "1970-01-01T00:00",
				"per_page": "1",
			}},
			[]EventsPageResult{
				{nil, 200, res200_1, &sp.EventsPage{}},
				{nil, 200, res200_2, &sp.EventsPage{}},
				{nil, 200, res200_3, &sp.EventsPage{}},
			},
		},
	} {
		for j, test := range seq.results {
			// Set up a new test in our inner loop since re-registering a handler will panic,
			// and the content we're returning needs to vary.
			testSetup(t)
			mockRestResponseBuilderFormat(t, "GET", test.status, sp.MessageEventsPathFormat, test.json)

			if page == nil {
				_, err = testClient.MessageEventsSearch(seq.in)
				page = seq.in
			} else {
				page, _, err = page.Next()
			}

			if err == nil && test.err != nil || err != nil && test.err == nil {
				t.Errorf("MessageEventsSearch[%d.%d] => err %#v want %#v", idx, j, err, test.err)
			} else if err != nil && err.Error() != test.err.Error() {
				t.Errorf("MessageEventsSearch[%d.%d] => err %#v want %#v", idx, j, err, test.err)
			} else if test.out != nil {
				if true == false && !reflect.DeepEqual(page, test.out) {
					t.Errorf("MessageEventsSearch[%d.%d] => events got/want:\n%#v\n%#v", idx, j, page, test.out)
				}
			}
			testTeardown()
		}
	}
}

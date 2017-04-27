package gosparkpost_test

import (
	"fmt"
	"reflect"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

var msgEventsEmpty string = `{
	"links": [],
	"results": [],
	"total_count": 0
}`

type EventsPageResult struct {
	err    error
	status int
	json   string
	out    *sp.EventsPage
}

func TestMessageEventsSearch(t *testing.T) {
	var err error
	var next *sp.EventsPage

	// Each test can return multiple pages of results
	for idx, _test := range []struct {
		in  *sp.EventsPage
		res []EventsPageResult
	}{
		{nil, []EventsPageResult{
			{errors.New("MessageEventsSearch called with nil EventsPage!"), 400, `{}`, nil},
		}},

		{&sp.EventsPage{
			Params: map[string]string{
				"from":   "1970-01-01T00:00",
				"events": "injection",
			}},
			[]EventsPageResult{
				{nil, 200, msgEventsEmpty, nil},
				{nil, 200, `{`, nil},
			},
		},
	} {
		for j, test := range _test.res {
			testSetup(t)
			defer testTeardown()

			path := sp.MessageEventsPathFormat
			if j > 0 {
				path = path + fmt.Sprintf("?page=%d", j)
			}
			mockRestResponseBuilderFormat(t, "GET", test.status, path, test.json)

			if j == 0 {
				_, err = testClient.MessageEventsSearch(_test.in)
			} else {
				_test.in.NextPage = fmt.Sprintf(path, testClient.Config.ApiVersion)
				next, _, err = _test.in.Next()
			}

			if err == nil && test.err != nil || err != nil && test.err == nil {
				t.Errorf("MessageEventsSearch[%d.%d] => err %#v want %#v", idx, j, err, test.err)
			} else if err != nil && err.Error() != test.err.Error() {
				t.Errorf("MessageEventsSearch[%d.%d] => err %#v want %#v", idx, j, err, test.err)
			} else if test.out != nil {
				if !reflect.DeepEqual(test.out, _test.in) {
					t.Errorf("MessageEventsSearch[%d.%d] => events got/want:\n%#v\n%#v", idx, j, _test.in, test.out)
				}
			}

			if j > 0 {
				_test.in = next
			}
		}
	}
}

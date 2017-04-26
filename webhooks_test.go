package gosparkpost_test

import (
	"reflect"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

func TestWebhookStatus(t *testing.T) {
	var res200 = loadTestFile(t, "test/json/webhook_status_200.json")
	var params = map[string]string{"timezone": "UTC"}

	for idx, test := range []struct {
		in     *sp.WebhookStatusWrapper
		err    error
		status int
		json   string
		out    *sp.WebhookStatusWrapper
	}{
		{nil, errors.New("WebhookStatus called with nil WebhookStatusWrapper"), 400, `{}`, nil},
		{&sp.WebhookStatusWrapper{ID: "id"}, errors.New("parsing api response: unexpected end of JSON input"), 200, res200[:len(res200)-2], nil},

		{&sp.WebhookStatusWrapper{ID: "id", WebhookCommon: sp.WebhookCommon{Params: params}},
			nil, 200, res200, &sp.WebhookStatusWrapper{
				ID: "id",
				Results: []sp.WebhookStatus{
					sp.WebhookStatus{
						BatchID:      "032d330540298f54f0e8bcc1373f3cfd",
						Timestamp:    "2014-07-30T21:38:08.000Z",
						Attempts:     7,
						ResponseCode: "200",
						FailureCode:  "",
					},
					sp.WebhookStatus{
						BatchID:      "13c6764994a8f6b4e29906d5712ca7d",
						Timestamp:    "2014-07-30T20:38:08.000Z",
						Attempts:     2,
						ResponseCode: "400",
						FailureCode:  "400"},
				},

				WebhookCommon: sp.WebhookCommon{Params: params}},
		},
	} {
		testSetup(t)
		defer testTeardown()

		id := ""
		if test.in != nil {
			id = test.in.ID
		}
		mockRestResponseBuilderFormat(t, "GET", test.status, sp.WebhooksPathFormat+"/"+id+"/batch-status", test.json)

		_, err := testClient.WebhookStatus(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("WebhookStatus[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("WebhookStatus[%d] => err %q want %q", idx, err, test.err)
		} else if test.out != nil {
			if !reflect.DeepEqual(test.out, test.in) {
				t.Errorf("WebhookStatus[%d] => webhook got/want:\n%#v\n%#v", idx, test.in, test.out)
			}
		}
	}
}

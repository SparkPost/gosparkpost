package gosparkpost_test

import (
	"reflect"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

func TestWebhookBadHost(t *testing.T) {
	testSetup(t)
	defer testTeardown()
	// Get the request to fail by mangling the HTTP host
	testClient.Config.BaseUrl = "%zz"

	_, err := testClient.Webhooks(&sp.WebhookListWrapper{})
	if err == nil || err.Error() != `building request: parse %zz/api/v1/webhooks: invalid URL escape "%zz"` {
		t.Errorf("error: %#v", err)
	}

	_, err = testClient.WebhookStatus(&sp.WebhookStatusWrapper{ID: "id"})
	if err == nil || err.Error() != `building request: parse %zz/api/v1/webhooks/id/batch-status: invalid URL escape "%zz"` {
		t.Errorf("error: %#v", err)
	}

	_, err = testClient.WebhookDetail(&sp.WebhookDetailWrapper{ID: "id"})
	if err == nil || err.Error() != `building request: parse %zz/api/v1/webhooks/id: invalid URL escape "%zz"` {
		t.Errorf("error: %#v", err)
	}

}

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
			}},
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
			if !reflect.DeepEqual(test.out.Results, test.in.Results) {
				t.Errorf("WebhookStatus[%d] => webhook got/want:\n%#v\n%#v", idx, test.in.Results, test.out.Results)
			}
		}
	}
}

func TestWebhookDetail(t *testing.T) {
	var res200 = loadTestFile(t, "test/json/webhook_detail_200.json")

	for idx, test := range []struct {
		in     *sp.WebhookDetailWrapper
		err    error
		status int
		json   string
		out    *sp.WebhookDetailWrapper
	}{
		{nil, errors.New("WebhookDetail called with nil WebhookDetailWrapper"), 400, `{}`, nil},
		{&sp.WebhookDetailWrapper{ID: "id"}, errors.New("parsing api response: unexpected end of JSON input"),
			200, res200[:len(res200)-2], nil},

		{&sp.WebhookDetailWrapper{ID: "id"}, nil, 200, res200, &sp.WebhookDetailWrapper{
			Results: &sp.WebhookItem{
				ID: "", Name: "Example webhook", Target: "http://client.example.com/example-webhook", Events: []string{"delivery", "injection", "open", "click"}, AuthType: "oauth2", AuthRequestDetails: struct {
					URL  string "json:\"url,omitempty\""
					Body struct {
						ClientID     string "json:\"client_id,omitempty\""
						ClientSecret string "json:\"client_secret,omitempty\""
					} "json:\"body,omitempty\""
				}{URL: "https://oauth.myurl.com/tokens", Body: struct {
					ClientID     string "json:\"client_id,omitempty\""
					ClientSecret string "json:\"client_secret,omitempty\""
				}{ClientID: "<oauth client id>", ClientSecret: "<oauth client secret>"}}, AuthCredentials: struct {
					Username    string "json:\"username,omitempty\""
					Password    string "json:\"password,omitempty\""
					AccessToken string "json:\"access_token,omitempty\""
					ExpiresIn   int    "json:\"expires_in,omitempty\""
				}{Username: "", Password: "", AccessToken: "<oauth token>", ExpiresIn: 3600}, AuthToken: "", LastSuccessful: "", LastFailure: "", Links: []struct {
					Href   string   "json:\"href,omitempty\""
					Rel    string   "json:\"rel,omitempty\""
					Method []string "json:\"method,omitempty\""
				}{struct {
					Href   string   "json:\"href,omitempty\""
					Rel    string   "json:\"rel,omitempty\""
					Method []string "json:\"method,omitempty\""
				}{Href: "http://www.messagesystems-api-url.com/api/v1/webhooks/12affc24-f183-11e3-9234-3c15c2c818c2/validate", Rel: "urn.msys.webhooks.validate", Method: []string{"POST"}}, struct {
					Href   string   "json:\"href,omitempty\""
					Rel    string   "json:\"rel,omitempty\""
					Method []string "json:\"method,omitempty\""
				}{Href: "http://www.messagesystems-api-url.com/api/v1/webhooks/12affc24-f183-11e3-9234-3c15c2c818c2/batch-status", Rel: "urn.msys.webhooks.batches", Method: []string{"GET"}}}}}},
	} {
		testSetup(t)
		defer testTeardown()

		id := "foo"
		if test.in != nil {
			id = test.in.ID
		}
		mockRestResponseBuilderFormat(t, "GET", test.status, sp.WebhooksPathFormat+"/"+id, test.json)

		_, err := testClient.WebhookDetail(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("WebhookDetail[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("WebhookDetail[%d] => err %q want %q", idx, err, test.err)
		} else if test.out != nil {
			if !reflect.DeepEqual(test.out.Results, test.in.Results) {
				t.Errorf("WebhookDetail[%d] => webhook got/want:\n%#v\n%#v", idx, test.in.Results, test.out.Results)
			}
		}
	}
}

func TestWebhooks(t *testing.T) {
	var res200 = loadTestFile(t, "test/json/webhooks_200.json")

	for idx, test := range []struct {
		in     *sp.WebhookListWrapper
		err    error
		status int
		json   string
		out    *sp.WebhookListWrapper
	}{
		{nil, errors.New("Webhooks called with nil WebhookListWrapper"), 400, `{}`, nil},
		{&sp.WebhookListWrapper{}, errors.New("parsing api response: unexpected end of JSON input"), 200, res200[:len(res200)-2], nil},

		{&sp.WebhookListWrapper{}, nil, 200, res200, &sp.WebhookListWrapper{Results: []sp.WebhookItem{
			sp.WebhookItem{ID: "a2b83490-10df-11e4-b670-c1ffa86371ff", Name: "Some webhook", Target: "http://client.example.com/some-webhook", Events: []string{"delivery", "injection", "open", "click"}, AuthType: "basic", AuthRequestDetails: struct {
				URL  string "json:\"url,omitempty\""
				Body struct {
					ClientID     string "json:\"client_id,omitempty\""
					ClientSecret string "json:\"client_secret,omitempty\""
				} "json:\"body,omitempty\""
			}{URL: "", Body: struct {
				ClientID     string "json:\"client_id,omitempty\""
				ClientSecret string "json:\"client_secret,omitempty\""
			}{ClientID: "", ClientSecret: ""}}, AuthCredentials: struct {
				Username    string "json:\"username,omitempty\""
				Password    string "json:\"password,omitempty\""
				AccessToken string "json:\"access_token,omitempty\""
				ExpiresIn   int    "json:\"expires_in,omitempty\""
			}{Username: "basicuser", Password: "somepass", AccessToken: "", ExpiresIn: 0}, AuthToken: "", LastSuccessful: "2014-08-01 16:09:15", LastFailure: "2014-06-01 15:15:45", Links: []struct {
				Href   string   "json:\"href,omitempty\""
				Rel    string   "json:\"rel,omitempty\""
				Method []string "json:\"method,omitempty\""
			}{struct {
				Href   string   "json:\"href,omitempty\""
				Rel    string   "json:\"rel,omitempty\""
				Method []string "json:\"method,omitempty\""
			}{Href: "http://www.messagesystems-api-url.com/api/v1/webhooks/a2b83490-10df-11e4-b670-c1ffa86371ff", Rel: "urn.msys.webhooks.webhook", Method: []string{"GET", "PUT"}}}},
			sp.WebhookItem{ID: "12affc24-f183-11e3-9234-3c15c2c818c2", Name: "Example webhook", Target: "http://client.example.com/example-webhook", Events: []string{"delivery", "injection", "open", "click"}, AuthType: "oauth2", AuthRequestDetails: struct {
				URL  string "json:\"url,omitempty\""
				Body struct {
					ClientID     string "json:\"client_id,omitempty\""
					ClientSecret string "json:\"client_secret,omitempty\""
				} "json:\"body,omitempty\""
			}{URL: "https://oauth.myurl.com/tokens", Body: struct {
				ClientID     string "json:\"client_id,omitempty\""
				ClientSecret string "json:\"client_secret,omitempty\""
			}{ClientID: "<oauth client id>", ClientSecret: "<oauth client secret>"}}, AuthCredentials: struct {
				Username    string "json:\"username,omitempty\""
				Password    string "json:\"password,omitempty\""
				AccessToken string "json:\"access_token,omitempty\""
				ExpiresIn   int    "json:\"expires_in,omitempty\""
			}{Username: "", Password: "", AccessToken: "<oauth token>", ExpiresIn: 3600}, AuthToken: "", LastSuccessful: "2014-07-01 16:09:15", LastFailure: "2014-08-01 15:15:45", Links: []struct {
				Href   string   "json:\"href,omitempty\""
				Rel    string   "json:\"rel,omitempty\""
				Method []string "json:\"method,omitempty\""
			}{struct {
				Href   string   "json:\"href,omitempty\""
				Rel    string   "json:\"rel,omitempty\""
				Method []string "json:\"method,omitempty\""
			}{Href: "http://www.messagesystems-api-url.com/api/v1/webhooks/12affc24-f183-11e3-9234-3c15c2c818c2", Rel: "urn.msys.webhooks.webhook", Method: []string{"GET", "PUT"}}}},
			sp.WebhookItem{ID: "123456-abcd-efgh-7890-123445566778", Name: "Another webhook", Target: "http://client.example.com/another-example", Events: []string{"generation_rejection", "generation_failure"}, AuthType: "none", AuthRequestDetails: struct {
				URL  string "json:\"url,omitempty\""
				Body struct {
					ClientID     string "json:\"client_id,omitempty\""
					ClientSecret string "json:\"client_secret,omitempty\""
				} "json:\"body,omitempty\""
			}{URL: "", Body: struct {
				ClientID     string "json:\"client_id,omitempty\""
				ClientSecret string "json:\"client_secret,omitempty\""
			}{ClientID: "", ClientSecret: ""}}, AuthCredentials: struct {
				Username    string "json:\"username,omitempty\""
				Password    string "json:\"password,omitempty\""
				AccessToken string "json:\"access_token,omitempty\""
				ExpiresIn   int    "json:\"expires_in,omitempty\""
			}{Username: "", Password: "", AccessToken: "", ExpiresIn: 0}, AuthToken: "5ebe2294ecd0e0f08eab7690d2a6ee69", LastSuccessful: "", LastFailure: "", Links: []struct {
				Href   string   "json:\"href,omitempty\""
				Rel    string   "json:\"rel,omitempty\""
				Method []string "json:\"method,omitempty\""
			}{struct {
				Href   string   "json:\"href,omitempty\""
				Rel    string   "json:\"rel,omitempty\""
				Method []string "json:\"method,omitempty\""
			}{Href: "http://www.messagesystems-api-url.com/api/v1/webhooks/123456-abcd-efgh-7890-123445566778", Rel: "urn.msys.webhooks.webhook", Method: []string{"GET", "PUT"}}}}}}},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "GET", test.status, sp.WebhooksPathFormat, test.json)

		_, err := testClient.Webhooks(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("Webhooks[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("Webhooks[%d] => err %q want %q", idx, err, test.err)
		} else if test.out != nil {
			if !reflect.DeepEqual(test.out.Results, test.in.Results) {
				t.Errorf("Webhooks[%d] => webhook got/want:\n%#v\n%#v", idx, test.in.Results, test.out.Results)
			}
		}
	}
}

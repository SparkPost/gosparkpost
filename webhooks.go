package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// WebhooksPathFormat is the path prefix used for webhook-related requests, with a format string for the API version.
var WebhooksPathFormat = "/api/v%d/webhooks"

// WebhookItem defines how webhook objects will be returned, as well as how they must be sent.
type WebhookItem struct {
	ID       string   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	Target   string   `json:"target,omitempty"`
	Events   []string `json:"events,omitempty"`
	AuthType string   `json:"auth_type,omitempty"`

	AuthRequestDetails struct {
		URL  string `json:"url,omitempty"`
		Body struct {
			ClientID     string `json:"client_id,omitempty"`
			ClientSecret string `json:"client_secret,omitempty"`
		} `json:"body,omitempty"`
	} `json:"auth_request_details,omitempty"`

	AuthCredentials struct {
		Username    string `json:"username,omitempty"`
		Password    string `json:"password,omitempty"`
		AccessToken string `json:"access_token,omitempty"`
		ExpiresIn   int    `json:"expires_in,omitempty"`
	} `json:"auth_credentials,omitempty"`

	AuthToken      string `json:"auth_token,omitempty"`
	LastSuccessful string `json:"last_successful,omitempty,omitempty"`
	LastFailure    string `json:"last_failure,omitempty,omitempty"`

	Links []struct {
		Href   string   `json:"href,omitempty"`
		Rel    string   `json:"rel,omitempty"`
		Method []string `json:"method,omitempty"`
	} `json:"links,omitempty"`
}

// WebhookStatus defines how the status of a webhook will be returned.
type WebhookStatus struct {
	BatchID      string `json:"batch_id,omitempty"`
	Timestamp    string `json:"ts,omitempty"`
	Attempts     int    `json:"attempts,omitempty"`
	ResponseCode string `json:"response_code,omitempty"`
	FailureCode  string `json:"failure_code,omitempty"`
}

// WebhookCommon contains fields common to all response types.
type WebhookCommon struct {
	Errors []interface{}     `json:"errors,omitempty"`
	Params map[string]string `json:"-"`
}

// WebhookListWrapper is passed into and updated by the Webhooks method, using results returned from the API.
type WebhookListWrapper struct {
	Results []WebhookItem `json:"results,omitempty"`
	WebhookCommon
}

// WebhookDetailWrapper is passed into and updated by the WebhookDetail method, using results returned from the API.
type WebhookDetailWrapper struct {
	ID      string       `json:"-"`
	Results *WebhookItem `json:"results,omitempty"`
	WebhookCommon
}

// WebhookStatusWrapper is passed into and updated by the WebhookStatus method, using results returned from the API.
type WebhookStatusWrapper struct {
	ID      string          `json:"-"`
	Results []WebhookStatus `json:"results,omitempty"`
	WebhookCommon
}

func buildUrl(c *Client, path string, parameters map[string]string) string {
	if parameters == nil || len(parameters) == 0 {
		path = fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	} else {
		params := url.Values{}
		for k, v := range parameters {
			params.Add(k, v)
		}

		path = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, params.Encode())
	}

	return path
}

// WebhookStatus updates its argument with details of batch delivery for the specified webhook.
// https://developers.sparkpost.com/api/#/reference/webhooks/batch-status/retrieve-status-information
func (c *Client) WebhookStatus(s *WebhookStatusWrapper) (*Response, error) {
	return c.WebhookStatusContext(context.Background(), s)
}

// WebhookStatusContext is the same as WebhookStatus, and allows the caller to specify their own context.
func (c *Client) WebhookStatusContext(ctx context.Context, s *WebhookStatusWrapper) (*Response, error) {
	if s == nil {
		return nil, errors.New("WebhookStatus called with nil WebhookStatusWrapper")
	}

	path := fmt.Sprintf(WebhooksPathFormat, c.Config.ApiVersion)
	finalUrl := buildUrl(c, path+"/"+s.ID+"/batch-status", s.Params)

	bodyBytes, res, err := doRequest(c, finalUrl, ctx)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodyBytes, s)
	if err != nil {
		return res, errors.Wrap(err, "parsing api response")
	}

	return res, err
}

// WebhookDetail updates its argument with details for the specified webhook.
// https://developers.sparkpost.com/api/#/reference/webhooks/retrieve/retrieve-webhook-details
func (c *Client) WebhookDetail(q *WebhookDetailWrapper) (*Response, error) {
	return c.WebhookDetailContext(context.Background(), q)
}

// WebhookDetailContext is the same as WebhookDetail, and allows the caller to specify their own context.
func (c *Client) WebhookDetailContext(ctx context.Context, d *WebhookDetailWrapper) (*Response, error) {
	if d == nil {
		return nil, errors.New("WebhookDetail called with nil WebhookDetailWrapper")
	}

	path := fmt.Sprintf(WebhooksPathFormat, c.Config.ApiVersion)
	finalUrl := buildUrl(c, path+"/"+d.ID, d.Params)

	bodyBytes, res, err := doRequest(c, finalUrl, ctx)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodyBytes, d)
	if err != nil {
		return res, errors.Wrap(err, "parsing api response")
	}

	return res, err
}

// Webhooks updates its argument with a list of all configured webhooks.
// https://developers.sparkpost.com/api/#/reference/webhooks/list/list-all-webhooks
func (c *Client) Webhooks(l *WebhookListWrapper) (*Response, error) {
	return c.WebhooksContext(context.Background(), l)
}

// WebhooksContext is the same as Webhooks, and allows the caller to specify their own context.
func (c *Client) WebhooksContext(ctx context.Context, l *WebhookListWrapper) (*Response, error) {
	if l == nil {
		return nil, errors.New("Webhooks called with nil WebhookListWrapper")
	}

	path := fmt.Sprintf(WebhooksPathFormat, c.Config.ApiVersion)
	finalUrl := buildUrl(c, path, l.Params)

	bodyBytes, res, err := doRequest(c, finalUrl, ctx)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodyBytes, l)
	if err != nil {
		return res, errors.Wrap(err, "parsing api response")
	}

	return res, err
}

func doRequest(c *Client, finalUrl string, ctx context.Context) ([]byte, *Response, error) {
	// Send off our request
	res, err := c.HttpGet(ctx, finalUrl)
	if err != nil {
		return nil, res, err
	}

	var body []byte
	// Assert that we got a JSON Content-Type back
	if body, err = res.AssertJson(); err != nil {
		return nil, res, err
	}

	return body, res, err
}

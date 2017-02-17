package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"

	"net/url"
)

// https://www.sparkpost.com/api#/reference/message-events
var WebhookListPathFormat = "/api/v%d/webhooks"
var WebhookQueryPathFormat = "/api/v%d/webhooks/%s"
var WebhookStatusPathFormat = "/api/v%d/webhooks/%s/batch-status"

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

type WebhookStatus struct {
	BatchID      string `json:"batch_id,omitempty"`
	Timestamp    string `json:"ts,omitempty"`
	Attempts     int    `json:"attempts,omitempty"`
	ResponseCode string `json:"response_code,omitempty"`
}

type WebhookCommon struct {
	Errors []interface{}     `json:"errors,omitempty"`
	Params map[string]string `json:"-"`
}

type WebhookListWrapper struct {
	Results []*WebhookItem `json:"results,omitempty"`
	WebhookCommon
}

type WebhookQueryWrapper struct {
	ID      string       `json:"-"`
	Results *WebhookItem `json:"results,omitempty"`
	WebhookCommon
}

type WebhookStatusWrapper struct {
	ID      string           `json:"-"`
	Results []*WebhookStatus `json:"results,omitempty"`
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

// https://developers.sparkpost.com/api/#/reference/webhooks/batch-status/retrieve-status-information
func (c *Client) WebhookStatus(s *WebhookStatusWrapper) (*Response, error) {
	return c.WebhookStatusContext(context.Background(), s)
}

func (c *Client) WebhookStatusContext(ctx context.Context, s *WebhookStatusWrapper) (*Response, error) {
	path := fmt.Sprintf(WebhookStatusPathFormat, c.Config.ApiVersion, s.ID)
	finalUrl := buildUrl(c, path, s.Params)

	bodyBytes, res, err := doRequest(c, finalUrl, ctx)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodyBytes, s)
	if err != nil {
		return res, err
	}

	return res, err
}

// https://developers.sparkpost.com/api/#/reference/webhooks/retrieve/retrieve-webhook-details
func (c *Client) QueryWebhook(q *WebhookQueryWrapper) (*Response, error) {
	return c.QueryWebhookContext(context.Background(), q)
}

func (c *Client) QueryWebhookContext(ctx context.Context, q *WebhookQueryWrapper) (*Response, error) {
	path := fmt.Sprintf(WebhookQueryPathFormat, c.Config.ApiVersion, q.ID)
	finalUrl := buildUrl(c, path, q.Params)

	bodyBytes, res, err := doRequest(c, finalUrl, ctx)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodyBytes, q)
	if err != nil {
		return res, err
	}

	return res, err
}

// https://developers.sparkpost.com/api/#/reference/webhooks/list/list-all-webhooks
func (c *Client) Webhooks(l *WebhookListWrapper) (*Response, error) {
	return c.WebhooksContext(context.Background(), l)
}

func (c *Client) WebhooksContext(ctx context.Context, l *WebhookListWrapper) (*Response, error) {
	path := fmt.Sprintf(WebhookListPathFormat, c.Config.ApiVersion)
	finalUrl := buildUrl(c, path, l.Params)

	bodyBytes, res, err := doRequest(c, finalUrl, ctx)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodyBytes, l)
	if err != nil {
		return res, err
	}

	return res, err
}

func doRequest(c *Client, finalUrl string, ctx context.Context) ([]byte, *Response, error) {
	// Send off our request
	res, err := c.HttpGet(ctx, finalUrl)
	if err != nil {
		return nil, res, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return nil, res, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return nil, res, err
	}

	return bodyBytes, res, err
}

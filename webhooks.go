package gosparkpost

import (
	"encoding/json"
	"fmt"

	URL "net/url"
)

// https://www.sparkpost.com/api#/reference/message-events
var webhookListPathFormat = "/api/v%d/webhooks"
var webhookQueryPathFormat = "/api/v%d/webhooks/%s"
var webhookStatusPathFormat = "/api/v%d/webhooks/%s/batch-status"

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
	Ts           string `json:"ts,omitempty"`
	Attempts     int    `json:"attempts,omitempty"`
	ResponseCode string `json:"response_code,omitempty"`
}

type WebhookListWrapper struct {
	Results []*WebhookItem `json:"results,omitempty"`
	Errors  []interface{}  `json:"errors,omitempty"`
	//{"errors":[{"param":"from","message":"From must be before to","value":"2014-07-20T09:00"},{"param":"to","message":"To must be in the format YYYY-MM-DDTHH:mm","value":"now"}]}
}

type WebhookQueryWrapper struct {
	Results *WebhookItem  `json:"results,omitempty"`
	Errors  []interface{} `json:"errors,omitempty"`
	//{"errors":[{"param":"from","message":"From must be before to","value":"2014-07-20T09:00"},{"param":"to","message":"To must be in the format YYYY-MM-DDTHH:mm","value":"now"}]}
}

type WebhookStatusWrapper struct {
	Results []*WebhookStatus `json:"results,omitempty"`
	Errors  []interface{}    `json:"errors,omitempty"`
	//{"errors":[{"param":"from","message":"From must be before to","value":"2014-07-20T09:00"},{"param":"to","message":"To must be in the format YYYY-MM-DDTHH:mm","value":"now"}]}
}

func buildUrl(c *Client, url string, parameters map[string]string) string {

	if parameters == nil || len(parameters) == 0 {
		url = fmt.Sprintf("%s%s", c.Config.BaseUrl, url)
	} else {
		params := URL.Values{}
		for k, v := range parameters {
			params.Add(k, v)
		}

		url = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, url, params.Encode())
	}

	return url
}

// https://developers.sparkpost.com/api/#/reference/webhooks/batch-status/retrieve-status-information
func (c *Client) WebhookStatus(id string, parameters map[string]string) (*WebhookStatusWrapper, error) {
	return c.WebhookStatusWithHeaders(id, parameters, nil)
}

func (c *Client) WebhookStatusWithHeaders(id string, parameters, headers map[string]string) (*WebhookStatusWrapper, error) {

	var finalUrl string
	path := fmt.Sprintf(webhookStatusPathFormat, c.Config.ApiVersion, id)

	finalUrl = buildUrl(c, path, parameters)

	return doWebhookStatusRequest(c, finalUrl, headers)
}

// https://developers.sparkpost.com/api/#/reference/webhooks/retrieve/retrieve-webhook-details
func (c *Client) WebhookQuery(id string, parameters map[string]string) (*WebhookQueryWrapper, error) {
	return c.WebhookQueryWithHeaders(id, parameters, nil)
}

func (c *Client) WebhookQueryWithHeaders(id string, parameters, headers map[string]string) (*WebhookQueryWrapper, error) {

	var finalUrl string
	path := fmt.Sprintf(webhookQueryPathFormat, c.Config.ApiVersion, id)

	finalUrl = buildUrl(c, path, parameters)

	return doWebhooksQueryRequest(c, finalUrl, headers)
}

// https://developers.sparkpost.com/api/#/reference/webhooks/list/list-all-webhooks
func (c *Client) WebhooksList(parameters map[string]string) (*WebhookListWrapper, error) {
	return c.WebhooksListWithHeaders(parameters, nil)
}
func (c *Client) WebhooksListWithHeaders(parameters, headers map[string]string) (*WebhookListWrapper, error) {

	var finalUrl string
	path := fmt.Sprintf(webhookListPathFormat, c.Config.ApiVersion)

	finalUrl = buildUrl(c, path, parameters)

	return doWebhooksListRequest(c, finalUrl, headers)
}

func doWebhooksListRequest(c *Client, finalUrl string, headers map[string]string) (*WebhookListWrapper, error) {

	bodyBytes, err := doRequest(c, finalUrl, headers)
	if err != nil {
		return nil, err
	}

	// Parse expected response structure
	var resMap WebhookListWrapper
	err = json.Unmarshal(bodyBytes, &resMap)

	if err != nil {
		return nil, err
	}

	return &resMap, err
}

func doWebhooksQueryRequest(c *Client, finalUrl string, headers map[string]string) (*WebhookQueryWrapper, error) {
	bodyBytes, err := doRequest(c, finalUrl, headers)

	// Parse expected response structure
	var resMap WebhookQueryWrapper
	err = json.Unmarshal(bodyBytes, &resMap)

	if err != nil {
		return nil, err
	}

	return &resMap, err
}

func doWebhookStatusRequest(c *Client, finalUrl string, headers map[string]string) (*WebhookStatusWrapper, error) {
	bodyBytes, err := doRequest(c, finalUrl, headers)

	// Parse expected response structure
	var resMap WebhookStatusWrapper
	err = json.Unmarshal(bodyBytes, &resMap)

	if err != nil {
		return nil, err
	}

	return &resMap, err
}

func doRequest(c *Client, finalUrl string, headers map[string]string) ([]byte, error) {
	// Send off our request
	res, err := c.HttpGet(finalUrl, headers)
	if err != nil {
		return nil, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return nil, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return nil, err
	}

	return bodyBytes, err
}

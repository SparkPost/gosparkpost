package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// https://developers.sparkpost.com/api/#/reference/suppression-list
var SuppressionListsPathFormat = "/api/v%d/suppression-list"

type SuppressionEntry struct {
	// Email is used when list is stored
	Email string `json:"email,omitempty"`

	// Recipient is used when a list is returned
	Recipient string `json:"recipient,omitempty"`

	Transactional    bool   `json:"transactional,omitempty"`
	NonTransactional bool   `json:"non_transactional,omitempty"`
	Source           string `json:"source,omitempty"`
	Type             string `json:type,omitempty`
	Description      string `json:"description,omitempty"`
	Updated          string `json:"updated,omitempty"`
	Created          string `json:"created,omitempty"`
}

type SuppressionListWrapper struct {
	Results    []*SuppressionEntry `json:"results,omitempty"`
	Recipients []SuppressionEntry  `json:"recipients,omitempty"`
}

func (c *Client) SuppressionList() (*SuppressionListWrapper, *Response, error) {
	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)
	return suppressionGet(c, c.Config.BaseUrl+path)
}

func (c *Client) SuppressionRetrieve(recipientEmail string) (*SuppressionListWrapper, *Response, error) {
	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)
	finalUrl := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, recipientEmail)

	return suppressionGet(c, finalUrl)
}

func (c *Client) SuppressionSearch(parameters map[string]string) (*SuppressionListWrapper, *Response, error) {
	var finalUrl string
	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)

	if parameters == nil || len(parameters) == 0 {
		finalUrl = fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	} else {
		params := url.Values{}
		for k, v := range parameters {
			params.Add(k, v)
		}

		finalUrl = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, params.Encode())
	}

	return suppressionGet(c, finalUrl)
}

func (c *Client) SuppressionDelete(recipientEmail string) (res *Response, err error) {
	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)
	finalUrl := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, recipientEmail)

	res, err = c.HttpDelete(context.TODO(), finalUrl)
	if err != nil {
		return res, err
	}

	if res.HTTP.StatusCode >= 200 && res.HTTP.StatusCode <= 299 {
		return res, err

	} else if len(res.Errors) > 0 {
		// handle common errors
		err = res.PrettyError("SuppressionEntry", "delete")
		if err != nil {
			return res, err
		}

		err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return res, err
}

func (c *Client) SuppressionInsertOrUpdate(entries []SuppressionEntry) (*Response, error) {
	if entries == nil {
		return nil, fmt.Errorf("send `entries` cannot be nil here")
	}

	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)
	list := SuppressionListWrapper{nil, entries}

	return suppressionPut(c, c.Config.BaseUrl+path, list)
}

func suppressionPut(c *Client, finalUrl string, recipients SuppressionListWrapper) (*Response, error) {
	jsonBytes, err := json.Marshal(recipients)
	if err != nil {
		return nil, err
	}

	res, err := c.HttpPut(context.TODO(), finalUrl, jsonBytes)
	if err != nil {
		return res, err
	}

	if err = res.AssertJson(); err != nil {
		return res, err
	}

	err = res.ParseResponse()
	if err != nil {
		return res, err
	}

	if res.HTTP.StatusCode == 200 {

	} else if len(res.Errors) > 0 {
		// handle common errors
		err = res.PrettyError("Transmission", "create")
		if err != nil {
			return res, err
		}

		err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return res, err
}

func suppressionGet(c *Client, finalUrl string) (*SuppressionListWrapper, *Response, error) {
	// Send off our request
	res, err := c.HttpGet(context.TODO(), finalUrl)
	if err != nil {
		return nil, res, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return nil, res, err
	}

	err = res.ParseResponse()
	if err != nil {
		return nil, res, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return nil, res, err
	}

	// Parse expected response structure
	var resMap SuppressionListWrapper
	err = json.Unmarshal(bodyBytes, &resMap)

	if err != nil {
		return nil, res, err
	}

	return &resMap, res, err
}

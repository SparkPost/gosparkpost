package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// SuppressionListsPathFormat https://developers.sparkpost.com/api/#/reference/suppression-list
var SuppressionListsPathFormat = "/api/v%d/suppression-list"

// SuppressionEntry stores a recipient’s opt-out preferences. It is a list of recipient email addresses to which you do NOT want to send email.
// https://developers.sparkpost.com/api/suppression-list.html#header-list-entry-attributes
type SuppressionEntry struct {
	// Email is used when list is stored
	Email string `json:"email,omitempty"`

	// Recipient is used when a list is returned
	Recipient string `json:"recipient,omitempty"`

	Transactional    bool   `json:"transactional,omitempty"`
	NonTransactional bool   `json:"non_transactional,omitempty"`
	Source           string `json:"source,omitempty"`
	Type             string `json:"type,omitempty"`
	Description      string `json:"description,omitempty"`
	Updated          string `json:"updated,omitempty"`
	Created          string `json:"created,omitempty"`
}

// WritableSuppressionEntry stores a recipient’s opt-out preferences. It is a list of recipient email addresses to which you do NOT want to send email.
// https://developers.sparkpost.com/api/suppression-list.html#suppression-list-bulk-insert-update-put
type WritableSuppressionEntry struct {
	// Recipient is used when a list is returned
	Recipient   string `json:"recipient,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

// SuppressionPage wraps suppression entries and response meta information
type SuppressionPage struct {
	client *Client

	Results    []*SuppressionEntry `json:"results,omitempty"`
	Recipients []SuppressionEntry  `json:"recipients,omitempty"`
	Errors     []struct {
		Message string `json:"message,omitempty"`
	} `json:"errors,omitempty"`

	TotalCount int `json:"total_count,omitempty"`

	NextPage  string
	PrevPage  string
	FirstPage string
	LastPage  string

	Links []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links,omitempty"`

	Params map[string]string `json:"-"`
}

// SuppressionList retrieves the account's suppression list.
// Suppression lists larger than 10,000 entries will need to use cursor to retrieve more results.
// See https://developers.sparkpost.com/api/suppression-list.html#suppression-list-search-get
func (c *Client) SuppressionList(sp *SuppressionPage) (*Response, error) {
	return c.SuppressionListContext(context.Background(), sp)
}

// SuppressionListContext retrieves the account's suppression list
func (c *Client) SuppressionListContext(ctx context.Context, sp *SuppressionPage) (*Response, error) {
	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)
	return c.suppressionGet(ctx, c.Config.BaseUrl+path, sp)
}

// SuppressionRetrieve retrieves the suppression status for a specific recipient by specifying the recipient’s email address
// // https://developers.sparkpost.com/api/suppression-list.html#suppression-list-retrieve,-delete,-insert-or-update-get
func (c *Client) SuppressionRetrieve(email string, sp *SuppressionPage) (*Response, error) {
	return c.SuppressionRetrieveContext(context.Background(), email, sp)
}

//SuppressionRetrieveContext retrieves the suppression status for a specific recipient by specifying the recipient’s email address
// // https://developers.sparkpost.com/api/suppression-list.html#suppression-list-retrieve,-delete,-insert-or-update-get
func (c *Client) SuppressionRetrieveContext(ctx context.Context, email string, sp *SuppressionPage) (*Response, error) {
	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)
	finalURL := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, email)

	return c.suppressionGet(ctx, finalURL, sp)
}

// SuppressionSearch search for suppression entries. For a list of parameters see
// https://developers.sparkpost.com/api/suppression-list.html#suppression-list-search-get
func (c *Client) SuppressionSearch(sp *SuppressionPage) (*Response, error) {
	return c.SuppressionSearchContext(context.Background(), sp)
}

// SuppressionSearchContext search for suppression entries. For a list of parameters see
// https://developers.sparkpost.com/api/suppression-list.html#suppression-list-search-get
func (c *Client) SuppressionSearchContext(ctx context.Context, sp *SuppressionPage) (*Response, error) {
	var finalURL string
	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)

	if sp.Params == nil || len(sp.Params) == 0 {
		finalURL = fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	} else {
		args := url.Values{}
		for k, v := range sp.Params {
			args.Add(k, v)
		}

		finalURL = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, args.Encode())
	}

	return c.suppressionGet(ctx, finalURL, sp)
}

// Next returns the next page of results from a previous MessageEventsSearch call
func (sp *SuppressionPage) Next() (*SuppressionPage, *Response, error) {
	return sp.NextContext(context.Background())
}

// NextContext is the same as Next, and it accepts a context.Context
func (sp *SuppressionPage) NextContext(ctx context.Context) (*SuppressionPage, *Response, error) {
	if sp.NextPage == "" {
		return nil, nil, nil
	}

	suppressionPage := &SuppressionPage{}
	suppressionPage.client = sp.client
	finalURL := fmt.Sprintf("%s", sp.client.Config.BaseUrl+sp.NextPage)
	res, err := sp.client.suppressionGet(ctx, finalURL, suppressionPage)

	return suppressionPage, res, err
}

// SuppressionDelete deletes an entry from the suppression list
func (c *Client) SuppressionDelete(email string) (res *Response, err error) {
	return c.SuppressionDeleteContext(context.Background(), email)
}

// SuppressionDeleteContext deletes an entry from the suppression list
func (c *Client) SuppressionDeleteContext(ctx context.Context, email string) (res *Response, err error) {
	if email == "" {
		err = fmt.Errorf("Deleting a suppression entry requires an email address")
		return nil, err
	}

	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)
	finalURL := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, email)

	res, err = c.HttpDelete(ctx, finalURL)
	if err != nil {
		return res, err
	}

	// We get an empty response on success. If there are errors we get JSON.
	if _, err = res.AssertJson(); err == nil {
		err = res.ParseResponse()
		if err != nil {
			return res, err
		}
	}

	return res, res.HTTPError()
}

// SuppressionUpsert adds an entry to the suppression, or updates the existing entry
func (c *Client) SuppressionUpsert(entries []WritableSuppressionEntry) (*Response, error) {
	return c.SuppressionUpsertContext(context.Background(), entries)
}

// SuppressionUpsertContext is the same as SuppressionUpsert, and it accepts a context.Context
func (c *Client) SuppressionUpsertContext(ctx context.Context, entries []WritableSuppressionEntry) (*Response, error) {
	if entries == nil {
		return nil, fmt.Errorf("`entries` cannot be nil")
	}

	path := fmt.Sprintf(SuppressionListsPathFormat, c.Config.ApiVersion)

	type EntriesWrapper struct {
		Recipients []WritableSuppressionEntry `json:"recipients,omitempty"`
	}

	entriesWrapper := EntriesWrapper{entries}

	// Marshaling a static type won't fail
	jsonBytes, _ := json.Marshal(entriesWrapper)

	finalURL := c.Config.BaseUrl + path
	return c.HttpPutJson(ctx, finalURL, jsonBytes)
}

// Wraps call to server and unmarshals response
func (c *Client) suppressionGet(ctx context.Context, finalURL string, sp *SuppressionPage) (*Response, error) {

	// Send off our request
	res, err := c.HttpGet(ctx, finalURL)
	if err != nil {
		return res, err
	}

	var body []byte
	// Assert that we got a JSON Content-Type back
	if body, err = res.AssertJson(); err != nil {
		return res, err
	}

	err = res.ParseResponse()
	if err != nil {
		return res, err
	}

	// Parse expected response structure
	err = json.Unmarshal(body, sp)
	if err != nil {
		return res, err
	}

	// For usage convenience parse out common links
	for _, link := range sp.Links {
		switch link.Rel {
		case "next":
			sp.NextPage = link.Href
		case "previous":
			sp.PrevPage = link.Href
		case "first":
			sp.FirstPage = link.Href
		case "last":
			sp.LastPage = link.Href
		}
	}

	if sp.client == nil {
		sp.client = c
	}

	return res, nil
}

package gosparkpost

import (
	"encoding/json"
	"fmt"
	URL "net/url"
)

// https://developers.sparkpost.com/api/#/reference/suppression-list
var supressionListsPathFormat = "/api/v%d/suppression-list"

type SupressionEntry struct {
	Recipient        string `json:"recipient,omitempty"`
	Transactional    bool   `json:"transactional,omitempty"`
	NonTransactional bool   `json:"non_transactional,omitempty"`
	Source           string `json:"source,omitempty"`
	Description      string `json:"description,omitempty"`
	Updated          string `json:"updated,omitempty"`
	Created          string `json:"created,omitempty"`
}

type SupressionListWrapper struct {
	Results []*SupressionEntry `json:"results,omitempty"`
}

func (c *Client) SupressionList() (*SupressionListWrapper, error) {
	path := fmt.Sprintf(supressionListsPathFormat, c.Config.ApiVersion)
	finalUrl := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)

	return doSupressionRequest(c, finalUrl)
}

func (c *Client) SupressionRetrieve(recipientEmail string) (*SupressionListWrapper, error) {
	path := fmt.Sprintf(supressionListsPathFormat, c.Config.ApiVersion)
	finalUrl := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, recipientEmail)

	return doSupressionRequest(c, finalUrl)
}

func (c *Client) SupressionSearch(parameters map[string]string) (*SupressionListWrapper, error) {
	var finalUrl string
	path := fmt.Sprintf(supressionListsPathFormat, c.Config.ApiVersion)

	if parameters == nil || len(parameters) == 0 {
		finalUrl = fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	} else {
		params := URL.Values{}
		for k, v := range parameters {
			params.Add(k, v)
		}

		finalUrl = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, params.Encode())
	}

	return doSupressionRequest(c, finalUrl)
}

func (c *Client) SupressionDelete(recipientEmail string) (res *Response, err error) {
	path := fmt.Sprintf(supressionListsPathFormat, c.Config.ApiVersion)
	finalUrl := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, recipientEmail)

	res, err = c.HttpDelete(finalUrl)
	if err != nil {
		return
	}

	if res.HTTP.StatusCode >= 200 && res.HTTP.StatusCode <= 299 {
		return

	} else if len(res.Errors) > 0 {
		// handle common errors
		err = res.PrettyError("SupressionEntry", "delete")
		if err != nil {
			return
		}

		err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return
}

func doSupressionRequest(c *Client, finalUrl string) (*SupressionListWrapper, error) {
	// Send off our request
	res, err := c.HttpGet(finalUrl)
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

	/*// DEBUG
	err = iou.WriteFile("./supressionlist.json", bodyBytes, 0644)
	if err != nil {
		return nil, err
	}
	*/

	// Parse expected response structure
	var resMap SupressionListWrapper
	err = json.Unmarshal(bodyBytes, &resMap)

	if err != nil {
		return nil, err
	}

	return &resMap, err
}

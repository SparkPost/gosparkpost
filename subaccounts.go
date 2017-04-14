package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// https://www.sparkpost.com/api#/reference/subaccounts
var SubaccountsPathFormat = "/api/v%d/subaccounts"
var SubaccountGrants = []string{
	"smtp/inject",
	"sending_domains/manage",
	"message_events/view",
	"suppression_lists/manage",
	"tracking_domains/view",
	"tracking_domains/manage",
	"transmissions/view",
	"transmissions/modify",
}
var SubaccountStatuses = []string{
	"active",
	"suspended",
	"terminated",
}

// Subaccount is the JSON structure accepted by and returned from the SparkPost Subaccounts API.
type Subaccount struct {
	ID               int      `json:"subaccount_id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Key              string   `json:"key,omitempty"`
	KeyLabel         string   `json:"key_label,omitempty"`
	Grants           []string `json:"key_grants,omitempty"`
	ShortKey         string   `json:"short_key,omitempty"`
	Status           string   `json:"status,omitempty"`
	ComplianceStatus string   `json:"compliance_status,omitempty"`
	IPPool           string   `json:"ip_pool,omitempty"`
}

// SubaccountCreate accepts a populated Subaccount object, validates it,
// and performs an API call against the configured endpoint.
func (c *Client) SubaccountCreate(s *Subaccount) (res *Response, err error) {
	return c.SubaccountCreateContext(context.Background(), s)
}

// SubaccountCreateContext is the same as SubaccountCreate, and it allows the caller to pass in a context
func (c *Client) SubaccountCreateContext(ctx context.Context, s *Subaccount) (res *Response, err error) {
	// enforce required parameters
	if s == nil {
		err = errors.New("Create called with nil Subaccount")
	} else if s.Name == "" {
		err = errors.New("Subaccount requires a non-empty Name")
	} else if s.KeyLabel == "" {
		err = errors.New("Subaccount requires a non-empty Key Label")
	} else
	// enforce max lengths
	if len(s.Name) > 1024 {
		err = errors.New("Subaccount name may not be longer than 1024 bytes")
	} else if len(s.KeyLabel) > 1024 {
		err = errors.New("Subaccount key label may not be longer than 1024 bytes")
	} else if s.IPPool != "" && len(s.IPPool) > 20 {
		err = errors.New("Subaccount ip pool may not be longer than 20 bytes")
	}
	if err != nil {
		return
	}

	if len(s.Grants) == 0 {
		s.Grants = SubaccountGrants
	}

	// Marshaling a static type won't fail
	jsonBytes, _ := json.Marshal(s)

	path := fmt.Sprintf(SubaccountsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	res, err = c.HttpPost(ctx, url, jsonBytes)
	if err != nil {
		return
	}

	if err = res.AssertJson(); err != nil {
		return
	}

	err = res.ParseResponse()
	if err != nil {
		return
	}

	if res.HTTP.StatusCode == 200 {
		var ok bool
		var results map[string]interface{}
		if results, ok = res.Results.(map[string]interface{}); !ok {
			return res, errors.New("Unexpected response to Subaccount creation (results)")
		}
		f, ok := results["subaccount_id"].(float64)
		if !ok {
			err = errors.New("Unexpected response to Subaccount creation (subaccount_id)")
		}
		s.ID = int(f)
		s.ShortKey, ok = results["short_key"].(string)
		if !ok {
			err = errors.New("Unexpected response to Subaccount creation (short_key)")
		}

	} else if len(res.Errors) > 0 {
		err = res.Errors
	}

	return
}

// SubaccountUpdate updates a subaccount with the specified id.
// Actually it will marshal and send all the subaccount fields, but that must not be a problem,
// as fields not supposed for update will be omitted
func (c *Client) SubaccountUpdate(s *Subaccount) (res *Response, err error) {
	return c.SubaccountUpdateContext(context.Background(), s)
}

// SubaccountUpdateContext is the same as SubaccountUpdate, and it allows the caller to provide a context
func (c *Client) SubaccountUpdateContext(ctx context.Context, s *Subaccount) (res *Response, err error) {
	if s.ID == 0 {
		err = errors.New("Subaccount Update called with zero id")
	} else if len(s.Name) > 1024 {
		err = errors.New("Subaccount name may not be longer than 1024 bytes")
	} else if s.Status != "" {
		found := false
		for _, v := range SubaccountStatuses {
			if s.Status == v {
				found = true
			}
		}
		if !found {
			err = errors.New("Not a valid subaccount status")
		}
	}

	if err != nil {
		return
	}

	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return
	}

	path := fmt.Sprintf(SubaccountsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, s.ID)

	res, err = c.HttpPut(ctx, url, jsonBytes)
	if err != nil {
		return
	}

	if err = res.AssertJson(); err != nil {
		return
	}

	err = res.ParseResponse()
	if err != nil {
		return
	}

	if res.HTTP.StatusCode == 200 {
		return

	} else if len(res.Errors) > 0 {
		err = res.Errors
	}

	return
}

// Subaccounts returns metadata for all Templates in the system.
func (c *Client) Subaccounts() (subaccounts []Subaccount, res *Response, err error) {
	return c.SubaccountsContext(context.Background())
}

// SubaccountsContext is the same as Subaccounts, and it allows the caller to provide a context
func (c *Client) SubaccountsContext(ctx context.Context) (subaccounts []Subaccount, res *Response, err error) {
	path := fmt.Sprintf(SubaccountsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	res, err = c.HttpGet(ctx, url)
	if err != nil {
		return
	}

	err = res.AssertJson()
	if err != nil {
		return
	}

	if res.HTTP.StatusCode == 200 {
		var body []byte
		body, err = res.ReadBody()
		if err != nil {
			return
		}
		slist := map[string][]Subaccount{}
		err = json.Unmarshal(body, &slist)
		if err != nil {
			return
		} else if list, ok := slist["results"]; ok {
			subaccounts = list
			return
		}
		err = errors.New("Unexpected response to Subaccount list")
		return

	} else {
		err = res.Errors
	}

	return
}

// Subaccount looks up a subaccount by its id
func (c *Client) Subaccount(id int) (subaccount *Subaccount, res *Response, err error) {
	return c.SubaccountContext(context.Background(), id)
}

// SubaccountContext is the same as Subaccount, and it accepts a context.Context
func (c *Client) SubaccountContext(ctx context.Context, id int) (subaccount *Subaccount, res *Response, err error) {
	path := fmt.Sprintf(SubaccountsPathFormat, c.Config.ApiVersion)
	u := fmt.Sprintf("%s%s/%d", c.Config.BaseUrl, path, id)
	res, err = c.HttpGet(ctx, u)
	if err != nil {
		return
	}

	err = res.AssertJson()
	if err != nil {
		return
	}

	if res.HTTP.StatusCode == 200 {
		if res.HTTP.StatusCode == 200 {
			var body []byte
			body, err = res.ReadBody()
			if err != nil {
				return
			}
			slist := map[string]Subaccount{}
			err = json.Unmarshal(body, &slist)
			if err != nil {
				return
			} else if s, ok := slist["results"]; ok {
				subaccount = &s
				return
			}
			err = errors.New("Unexpected response to Subaccount")
			return
		}
	} else {
		err = res.Errors
	}

	return
}

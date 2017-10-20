package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// SubaccountsPathFormat provides an easy way to fill out the path including the version.
var SubaccountsPathFormat = "/api/v%d/subaccounts"

// SubaccountGrants contains the grants that will be given to new subaccounts by default.
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

// SubaccountStatuses contains valid subaccount statuses.
var SubaccountStatuses = []string{
	"active",
	"suspended",
	"terminated",
}

// Subaccount is the JSON structure accepted by and returned from the SparkPost Subaccounts API.
type Subaccount struct {
	ID               int      `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Key              string   `json:"key,omitempty"`
	KeyLabel         string   `json:"key_label,omitempty"`
	Grants           []string `json:"key_grants,omitempty"`
	ShortKey         string   `json:"short_key,omitempty"`
	Status           string   `json:"status,omitempty"`
	ComplianceStatus string   `json:"compliance_status,omitempty"`
	IPPool           string   `json:"ip_pool,omitempty"`
}

// SubaccountCreate attempts to create a subaccount using the provided object
func (c *Client) SubaccountCreate(s *Subaccount) (res *Response, err error) {
	return c.SubaccountCreateContext(context.Background(), s)
}

// SubaccountCreateContext is the same as SubaccountCreate, and it allows the caller to pass in a context
// New subaccounts will have all grants in SubaccountGrants, unless s.Grants is non-nil.
func (c *Client) SubaccountCreateContext(ctx context.Context, s *Subaccount) (res *Response, err error) {
	// enforce required parameters
	if s == nil {
		err = errors.New("Create called with nil Subaccount")
		return
	}

	if len(s.Grants) == 0 {
		s.Grants = SubaccountGrants
	}

	// Marshaling a static type won't fail
	jsonBytes, _ := json.Marshal(s)
	path := fmt.Sprintf(SubaccountsPathFormat, c.Config.ApiVersion)

	res, err = c.HttpPostJson(ctx, c.Config.BaseUrl+path, jsonBytes, nil)
	if err != nil {
		return
	}

	if results, ok := res.Results.(map[string]interface{}); !ok {
		err = errors.New("Unexpected response to Subaccount creation (results)")
	} else if f, ok := results["subaccount_id"].(float64); !ok {
		err = errors.New("Unexpected response to Subaccount creation (subaccount_id)")
	} else {
		s.ID = int(f)
		if s.ShortKey, ok = results["short_key"].(string); !ok {
			err = errors.New("Unexpected response to Subaccount creation (short_key)")
		}
	}

	return
}

// SubaccountUpdate updates a subaccount with the specified id.
// It marshals and sends all the subaccount fields, ignoring the read-only ones.
func (c *Client) SubaccountUpdate(s *Subaccount) (res *Response, err error) {
	return c.SubaccountUpdateContext(context.Background(), s)
}

// SubaccountUpdateContext is the same as SubaccountUpdate, and it allows the caller to provide a context
func (c *Client) SubaccountUpdateContext(ctx context.Context, s *Subaccount) (res *Response, err error) {
	if s == nil {
		err = errors.New("Subaccount Update called with nil Subaccount")
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

	// Marshaling a static type won't fail
	jsonBytes, _ := json.Marshal(s)

	path := fmt.Sprintf(SubaccountsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s/%d", c.Config.BaseUrl, path, s.ID)

	return c.HttpPutJson(ctx, url, jsonBytes, nil)
}

// Subaccounts returns metadata for all Subaccounts in the system.
func (c *Client) Subaccounts() (subaccounts []Subaccount, res *Response, err error) {
	return c.SubaccountsContext(context.Background())
}

// SubaccountsContext is the same as Subaccounts, and it allows the caller to provide a context
func (c *Client) SubaccountsContext(ctx context.Context) (subaccounts []Subaccount, res *Response, err error) {
	path := fmt.Sprintf(SubaccountsPathFormat, c.Config.ApiVersion)
	slist := map[string][]Subaccount{}

	res, err = c.HttpGetJson(ctx, c.Config.BaseUrl+path, &slist)
	if err != nil {
		return
	}

	if list, ok := slist["results"]; ok {
		subaccounts = list
	} else {
		err = errors.New("Unexpected response to Subaccount list")
	}

	return
}

// Subaccount looks up a subaccount using the provided id
func (c *Client) Subaccount(id int) (subaccount *Subaccount, res *Response, err error) {
	return c.SubaccountContext(context.Background(), id)
}

// SubaccountContext is the same as Subaccount, and it accepts a context.Context
func (c *Client) SubaccountContext(ctx context.Context, id int) (subaccount *Subaccount, res *Response, err error) {
	path := fmt.Sprintf(SubaccountsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s/%d", c.Config.BaseUrl, path, id)

	// pass nil to skip auto-unmarshalling
	res, err = c.HttpGetJson(ctx, url, nil)
	if err != nil {
		return nil, res, err
	}

	if len(res.Errors) > 0 {
		return nil, res, res.Errors
	}

	out := map[string]Subaccount{}
	if err = json.Unmarshal(res.Body, &out); err != nil {
		err = errors.New("Unexpected response to Subaccount fetch")
	} else if sub, ok := out["results"]; !ok {
		err = errors.New("Unexpected response to Subaccount fetch (results)")
	} else {
		subaccount = &sub
	}

	return subaccount, res, err
}

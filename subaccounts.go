package gosparkpost

import (
	"encoding/json"
	"fmt"
)

// https://www.sparkpost.com/api#/reference/subaccounts
var subaccountsPathFormat = "/api/v%d/subaccounts"
var availableGrants = []string{
	"smtp/inject",
	"sending_domains/manage",
	"message_events/view",
	"suppression_lists/manage",
	"transmissions/view",
	"transmissions/modify",
}
var validStatuses = []string{
	"active",
	"suspended",
	"terminated",
}

// Subaccount is the JSON structure accepted by and returned from the SparkPost Subaccounts API.
type Subaccount struct {
	ID               string   `json:"subaccount_id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Key              string   `json:"key,omitempty"`
	KeyLabel         string   `json:"key_label,omitempty"`
	Grants           []string `json:"key_grants,omitempty"`
	ShortKey         string   `json:"short_key,omitempty"`
	Status           string   `json:"status,omitempty"`
	ComplianceStatus string   `json:"compliance_status,omitempty"`
}

// Create accepts a populated Subaccount object, validates it,
// and performs an API call against the configured endpoint.
func (c *Client) SubaccountCreate(s *Subaccount) (res *Response, err error) {
	// enforce required parameters
	if s == nil {
		err = fmt.Errorf("Create called with nil Subaccount")
	} else if s.Name == "" {
		err = fmt.Errorf("Subaccount requires a non-empty Name")
	} else if s.KeyLabel == "" {
		err = fmt.Errorf("Subaccount requires a non-empty Key Label")
	} else
	// enforce max lengths
	if len(s.Name) > 1024 {
		err = fmt.Errorf("Subaccount name may not be longer than 1024 bytes")
	} else if len(s.KeyLabel) > 1024 {
		err = fmt.Errorf("Subaccount key label may not be longer than 1024 bytes")
	}
	if err != nil {
		return
	}

	if len(s.Grants) == 0 {
		s.Grants = availableGrants
	}

	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return
	}

	path := fmt.Sprintf(subaccountsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	res, err = c.HttpPost(url, jsonBytes)
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
		s.ID, ok = res.Results["subaccount_id"].(string)
		s.ShortKey, ok = res.Results["short_key"].(string)
		if !ok {
			err = fmt.Errorf("Unexpected response to Subaccount creation")
		}

	} else if len(res.Errors) > 0 {
		// handle common errors
		err = res.PrettyError("Subaccount", "create")
		if err != nil {
			return
		}

		if res.HTTP.StatusCode == 422 { // subaccount syntax error
			eobj := res.Errors[0]
			err = fmt.Errorf("%s: %s\n%s", eobj.Code, eobj.Message, eobj.Description)
		} else { // everything else
			err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
		}
	}

	return
}

// Update updates a draft/published template with the specified id
func (c *Client) SubaccountUpdate(s *Subaccount) (res *Response, err error) {
	if s.ID == "" {
		err = fmt.Errorf("Subaccount Update called with blank id")
	} else if len(s.Name) > 1024 {
		err = fmt.Errorf("Subaccount name may not be longer than 1024 bytes")
	} else if s.Status != "" {
		found := false
		for _, v := range availableGrants {
			if s.Status == v {
				found = true
			}
		}
		if !found {
			err = fmt.Errorf("Not a valid subaccount status")
		}
	}

	if err != nil {
		return
	}

	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return
	}

	path := fmt.Sprintf(templatesPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, s.ID)

	res, err = c.HttpPut(url, jsonBytes)
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
		// handle common errors
		err = res.PrettyError("Subaccount", "update")
		if err != nil {
			return
		}

		// handle template-specific ones
		if res.HTTP.StatusCode == 409 {
			err = fmt.Errorf("Subaccount with id [%s] is in use by msg generation", s.ID)
		} else { // everything else
			err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
		}
	}

	return
}

// List returns metadata for all Templates in the system.
func (c *Client) Subaccounts() ([]Subaccount, *Response, error) {
	path := fmt.Sprintf(subaccountsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	res, err := c.HttpGet(url)
	if err != nil {
		return nil, nil, err
	}

	if err = res.AssertJson(); err != nil {
		return nil, res, err
	}

	if res.HTTP.StatusCode == 200 {
		var body []byte
		body, err = res.ReadBody()
		if err != nil {
			return nil, res, err
		}
		slist := map[string][]Subaccount{}
		if err = json.Unmarshal(body, &slist); err != nil {
			return nil, res, err
		} else if list, ok := slist["results"]; ok {
			return list, res, nil
		}
		return nil, res, fmt.Errorf("Unexpected response to Subaccount list")

	} else {
		err = res.ParseResponse()
		if err != nil {
			return nil, res, err
		}
		if len(res.Errors) > 0 {
			err = res.PrettyError("Subaccount", "list")
			if err != nil {
				return nil, res, err
			}
		}
		return nil, res, fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return nil, res, err
}

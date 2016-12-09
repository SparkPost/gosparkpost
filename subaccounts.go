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
	ID               int      `json:"subaccount_id,omitempty"`
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
		var results map[string]interface{}
		if results, ok = res.Results.(map[string]interface{}); !ok {
			return res, fmt.Errorf("Unexpected response to Subaccount creation (results)")
		}
		f, ok := results["subaccount_id"].(float64)
		if !ok {
			err = fmt.Errorf("Unexpected response to Subaccount creation")
		}
		s.ID = int(f)
		s.ShortKey, ok = results["short_key"].(string)
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

// Update updates a subaccount with the specified id.
// Actually it will marshal and send all the subaccount fields, but that must not be a problem,
// as fields not supposed for update will be omitted
func (c *Client) SubaccountUpdate(s *Subaccount) (res *Response, err error) {
	if s.ID == 0 {
		err = fmt.Errorf("Subaccount Update called with zero id")
	} else if len(s.Name) > 1024 {
		err = fmt.Errorf("Subaccount name may not be longer than 1024 bytes")
	} else if s.Status != "" {
		found := false
		for _, v := range validStatuses {
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
func (c *Client) Subaccounts() (subaccounts []Subaccount, res *Response, err error) {
	path := fmt.Sprintf(subaccountsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	res, err = c.HttpGet(url)
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
		err = fmt.Errorf("Unexpected response to Subaccount list")
		return

	} else {
		err = res.ParseResponse()
		if err != nil {
			return
		}
		if len(res.Errors) > 0 {
			err = res.PrettyError("Subaccount", "list")
			if err != nil {
				return
			}
		}
		err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
		return
	}

	return
}

func (c *Client) Subaccount(id int) (subaccount *Subaccount, res *Response, err error) {
	path := fmt.Sprintf(subaccountsPathFormat, c.Config.ApiVersion)
	u := fmt.Sprintf("%s%s/%d", c.Config.BaseUrl, path, id)
	res, err = c.HttpGet(u)
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
			err = fmt.Errorf("Unexpected response to Subaccount")
			return
		}
	} else {
		err = res.ParseResponse()
		if err != nil {
			return
		}
		if len(res.Errors) > 0 {
			err = res.PrettyError("Subaccount", "retrieve")
			if err != nil {
				return
			}
		}
		err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
		return
	}

	return
}

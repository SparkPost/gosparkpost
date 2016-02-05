package gosparkpost

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// https://www.sparkpost.com/api#/reference/transmissions
var transmissionsPathFormat = "/api/v%d/transmissions"

// Transmission is the JSON structure accepted by and returned from the SparkPost Transmissions API.
type Transmission struct {
	ID               string      `json:"id,omitempty"`
	State            string      `json:"state,omitempty"`
	Options          *TxOptions  `json:"options,omitempty"`
	Recipients       interface{} `json:"recipients"`
	CampaignID       string      `json:"campaign_id,omitempty"`
	Description      string      `json:"description,omitempty"`
	Metadata         interface{} `json:"metadata,omitempty"`
	SubstitutionData interface{} `json:"substitution_data,omitempty"`
	ReturnPath       string      `json:"return_path,omitempty"`
	Content          interface{} `json:"content"`

	TotalRecipients      *int `json:"total_recipients,omitempty"`
	NumGenerated         *int `json:"num_generated,omitempty"`
	NumFailedGeneration  *int `json:"num_failed_generation,omitempty"`
	NumInvalidRecipients *int `json:"num_invalid_recipients,omitempty"`
}

// Options specifies settings to apply to this Transmission.
// If not specified, and present in TmplOptions, those values will be used.
type TxOptions struct {
	TmplOptions

	StartTime       time.Time `json:"start_time,omitempty"`
	Sandbox         string    `json:"sandbox,omitempty"`
	SkipSuppression string    `json:"skip_suppression,omitempty"`
}

// ParseRecipients asserts that Transmission.Recipients is valid.
func ParseRecipients(recips interface{}) (ra *[]Recipient, err error) {
	switch rVal := recips.(type) {
	case map[string]interface{}:
		for k, v := range rVal {
			switch vVal := v.(type) {
			case string:
				if strings.EqualFold(k, "list_id") {
					return
				}
			default:
				err = fmt.Errorf("Transmission.Recipient objects must contain string values, not [%s]",
					reflect.TypeOf(vVal))
				return
			}
		}
		err = fmt.Errorf("Transmission.Recipient objects must contain a key `list_id`")
		return

	case map[string]string:
		for k, _ := range rVal {
			if strings.EqualFold(k, "list_id") {
				return
			}
		}
		err = fmt.Errorf("Transmission.Recipient objects must contain a key `list_id`")
		return

	case []string:
		raObj := make([]Recipient, len(rVal))
		for i, r := range rVal {
			// Make a full Recipient object from each string
			raObj[i] = Recipient{Address: map[string]string{"email": r}}
		}
		ra := &raObj
		return ra, nil

	case []interface{}:
		for _, v := range rVal {
			switch r := v.(type) {
			case Recipient:
				err = r.Validate()
				if err != nil {
					return
				}

			default:
				err = fmt.Errorf("Failed to parse inline Transmission.Recipient list")
				return
			}
		}

	case []Recipient:
		for _, v := range rVal {
			err = v.Validate()
			if err != nil {
				return
			}
		}

	default:
		err = fmt.Errorf("Unsupported Transmission.Recipient type [%s]", reflect.TypeOf(rVal))
		return
	}

	return
}

// ParseContent asserts that Transmission.Content is valid.
func ParseContent(content interface{}) (err error) {
	switch rVal := content.(type) {
	case map[string]interface{}:
		for k, v := range rVal {
			switch vVal := v.(type) {
			case string:
				if strings.EqualFold(k, "template_id") {
					return nil
				}
			default:
				return fmt.Errorf("Transmission.Content objects must contain string values, not [%s]", reflect.TypeOf(vVal))
			}
		}
		return fmt.Errorf("Transmission.Content objects must contain a key `template_id`")

	case map[string]string:
		for k, _ := range rVal {
			if strings.EqualFold(k, "template_id") {
				return nil
			}
		}
		return fmt.Errorf("Transmission.Content objects must contain a key `template_id`")

	case Content:
		te := &Template{Name: "tmp", Content: rVal}
		return te.Validate()

	default:
		return fmt.Errorf("Unsupported Transmission.Recipient type [%s]", reflect.TypeOf(rVal))
	}

	return
}

// Validate runs sanity checks of a Transmission struct.
// This should catch most errors before attempting a doomed API call.
func (t *Transmission) Validate() error {
	if t == nil {
		return fmt.Errorf("Can't Validate a nil Transmission")
	}

	// enforce required parameters
	if t.Recipients == nil {
		return fmt.Errorf("Transmission requires Recipients")
	} else if t.Content == nil {
		return fmt.Errorf("Transmission requires Content")
	}

	// enforce max lengths
	if len(t.CampaignID) > 64 {
		return fmt.Errorf("Campaign id may not be longer than 64 bytes")
	} else if len(t.Description) > 1024 {
		return fmt.Errorf("Transmission description may not be longer than 1024 bytes")
	}

	// validate members from other packages
	recips, err := ParseRecipients(t.Recipients)
	if err != nil {
		return err
	}
	// Use the updated Recipients object optionally returned from ParseRecipients
	if recips != nil {
		t.Recipients = *recips
	}

	err = ParseContent(t.Content)
	if err != nil {
		return err
	}

	// Metadata must be an object, not an array or bool etc.
	if t.Metadata != nil {
		err := AssertObject(t.Metadata, "metadata")
		if err != nil {
			return err
		}
	}

	// SubstitutionData must be an object, not an array or bool etc.
	if t.SubstitutionData != nil {
		err := AssertObject(t.SubstitutionData, "substitution_data")
		if err != nil {
			return err
		}
	}

	return nil
}

// Create accepts a populated Transmission object, performs basic sanity
// checks on it, and performs an API call against the configured endpoint.
// Calling this function can cause email to be sent, if used correctly.
func (c *Client) Send(t *Transmission) (id string, res *Response, err error) {
	if t == nil {
		err = fmt.Errorf("Create called with nil Transmission")
		return
	}

	err = t.Validate()
	if err != nil {
		return
	}

	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return
	}

	path := fmt.Sprintf(transmissionsPathFormat, c.Config.ApiVersion)
	u := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	res, err = c.HttpPost(u, jsonBytes)
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
		id, ok = res.Results["id"].(string)
		if !ok {
			err = fmt.Errorf("Unexpected response to Transmission creation")
		}

	} else if len(res.Errors) > 0 {
		// handle common errors
		err = res.PrettyError("Transmission", "create")
		if err != nil {
			return
		}

		err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return
}

// Retrieve accepts a Transmission.ID and retrieves the corresponding object.
func (c *Client) Transmission(id string) (*Transmission, *Response, error) {
	if nonDigit.MatchString(id) {
		return nil, nil, fmt.Errorf("id may only contain digits")
	}
	path := fmt.Sprintf(transmissionsPathFormat, c.Config.ApiVersion)
	u := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, id)
	res, err := c.HttpGet(u)
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

		// Unwrap the returned Transmission
		tmp := map[string]map[string]Transmission{}
		if err = json.Unmarshal(body, &tmp); err != nil {
			return nil, res, err
		} else if results, ok := tmp["results"]; ok {
			if tr, ok := results["transmission"]; ok {
				return &tr, res, nil
			} else {
				return nil, res, fmt.Errorf("Unexpected results structure in response")
			}
		}
		return nil, res, fmt.Errorf("Unexpected response to Transmission.Retrieve")

	} else {
		err = res.ParseResponse()
		if err != nil {
			return nil, res, err
		}
		if len(res.Errors) > 0 {
			err = res.PrettyError("Transmission", "retrieve")
			if err != nil {
				return nil, res, err
			}
		}
		return nil, res, fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return nil, res, err
}

// Delete attempts to remove the Transmission with the specified id.
// Only Transmissions which are scheduled for future generation may be deleted.
func (c *Client) TransmissionDelete(id string) (*Response, error) {
	if id == "" {
		return nil, fmt.Errorf("Delete called with blank id")
	}
	if nonDigit.MatchString(id) {
		return nil, fmt.Errorf("Transmissions.Delete: id may only contain digits")
	}

	path := fmt.Sprintf(transmissionsPathFormat, c.Config.ApiVersion)
	u := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, id)
	res, err := c.HttpDelete(u)
	if err != nil {
		return nil, err
	}

	if err = res.AssertJson(); err != nil {
		return res, err
	}

	if err = res.ParseResponse(); err != nil {
		return res, err
	}

	if res.HTTP.StatusCode == 200 {
		return res, nil

	} else if len(res.Errors) > 0 {
		// handle common errors
		err = res.PrettyError("Transmission", "delete")
		if err != nil {
			return res, err
		}

		return res, fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return res, nil
}

// List returns Transmission summary information for matching Transmissions.
// To skip filtering by campaign or template id, use a nil param.
func (c *Client) Transmissions(campaignID, templateID *string) ([]Transmission, *Response, error) {
	// If a query parameter is present and empty, that searches for blank IDs, as opposed
	// to when it is omitted entirely, which returns everything.
	qp := make([]string, 0, 2)
	if campaignID != nil {
		qp = append(qp, fmt.Sprintf("campaign_id=%s", url.QueryEscape(*campaignID)))
	}
	if templateID != nil {
		qp = append(qp, fmt.Sprintf("template_id=%s", url.QueryEscape(*templateID)))
	}

	qstr := ""
	if len(qp) > 0 {
		qstr = strings.Join(qp, "&")
	}
	path := fmt.Sprintf(transmissionsPathFormat, c.Config.ApiVersion)
	u := fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, qstr)

	res, err := c.HttpGet(u)
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
		tlist := map[string][]Transmission{}
		if err = json.Unmarshal(body, &tlist); err != nil {
			return nil, res, err
		} else if list, ok := tlist["results"]; ok {
			return list, res, nil
		}
		return nil, res, fmt.Errorf("Unexpected response to Transmission list")

	} else {
		err = res.ParseResponse()
		if err != nil {
			return nil, res, err
		}
		if len(res.Errors) > 0 {
			err = res.PrettyError("Transmission", "list")
			if err != nil {
				return nil, res, err
			}
		}
		return nil, res, fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}
}

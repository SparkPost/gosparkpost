package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// https://www.sparkpost.com/api#/reference/recipient-lists
var RecipientListsPathFormat = "/api/v%d/recipient-lists"

// RecipientList is the JSON structure accepted by and returned from the SparkPost Recipient Lists API.
// It's mostly metadata at this level - see Recipients for more detail.
type RecipientList struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Attributes  interface{} `json:"attributes,omitempty"`
	Recipients  []Recipient `json:"recipients"`

	Accepted *int `json:"total_accepted_recipients,omitempty"`
}

// Recipient represents one email (you guessed it) recipient.
type Recipient struct {
	Address          interface{} `json:"address"`
	ReturnPath       string      `json:"return_path,omitempty"`
	Tags             []string    `json:"tags,omitempty"`
	Metadata         interface{} `json:"metadata,omitempty"`
	SubstitutionData interface{} `json:"substitution_data,omitempty"`
}

// Address describes the nested object way of specifying the Recipient's email address.
// Recipient.Address can also be a plain string.
type Address struct {
	Email    string `json:"email"`
	Name     string `json:"name,omitempty"`
	HeaderTo string `json:"header_to,omitempty"`
}

// ParseAddress parses the various allowable Content.From values.
func ParseAddress(addr interface{}) (a Address, err error) {
	// handle the allowed types
	switch addrVal := addr.(type) {
	case string: // simple string value
		if addrVal == "" {
			err = errors.New("Recipient.Address may not be empty")
		} else {
			a.Email = addrVal
		}

	case Address:
		a = addr.(Address)

	case map[string]interface{}:
		// auto-parsed nested json object
		for k, v := range addrVal {
			switch vVal := v.(type) {
			case string:
				if strings.EqualFold(k, "name") {
					a.Name = vVal
				} else if strings.EqualFold(k, "email") {
					a.Email = vVal
				} else if strings.EqualFold(k, "header_to") {
					a.HeaderTo = vVal
				}
			default:
				err = errors.New("strings are required for all Recipient.Address values")
				break
			}
		}

	case map[string]string:
		// user-provided json literal (convenience)
		for k, v := range addrVal {
			if strings.EqualFold(k, "name") {
				a.Name = v
			} else if strings.EqualFold(k, "email") {
				a.Email = v
			} else if strings.EqualFold(k, "header_to") {
				a.HeaderTo = v
			}
		}

	default:
		err = errors.Errorf("unsupported Recipient.Address value type [%T]", addrVal)
	}

	return
}

// Validate runs sanity checks on a RecipientList struct. This should
// catch most errors before attempting a doomed API call.
func (rl *RecipientList) Validate() error {
	if rl == nil {
		return errors.New("Can't validate a nil RecipientList")
	}

	// enforce required parameters
	if rl.Recipients == nil || len(rl.Recipients) <= 0 {
		return errors.New("RecipientList requires at least one Recipient")
	}

	// enforce max lengths
	if len(rl.ID) > 64 {
		return errors.New("RecipientList id may not be longer than 64 bytes")
	} else if len(rl.Name) > 64 {
		return errors.New("RecipientList name may not be longer than 64 bytes")
	} else if len(rl.Description) > 1024 {
		return errors.New("RecipientList description may not be longer than 1024 bytes")
	}

	var err error
	for _, r := range rl.Recipients {
		err = r.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate runs sanity checks on a Recipient struct. This should
// catch most errors before attempting a doomed API call.
func (r Recipient) Validate() error {
	_, err := ParseAddress(r.Address)
	if err != nil {
		return err
	}
	return nil
}

// RecipientListCreate accepts a populated RecipientList object, validates it,
// and performs an API call against the configured endpoint.
func (c *Client) RecipientListCreate(rl *RecipientList) (id string, res *Response, err error) {
	return c.RecipientListCreateContext(context.Background(), rl)
}

// RecipientListCreateContext is the same as RecipientListCreate, and it accepts a context.Context
func (c *Client) RecipientListCreateContext(ctx context.Context, rl *RecipientList) (id string, res *Response, err error) {
	if rl == nil {
		err = errors.New("Create called with nil RecipientList")
		return
	}

	err = rl.Validate()
	if err != nil {
		return
	}

	jsonBytes, err := json.Marshal(rl)
	if err != nil {
		return
	}

	path := fmt.Sprintf(RecipientListsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	res, err = c.HttpPost(ctx, url, jsonBytes)
	if err != nil {
		return
	}

	if _, err = res.AssertJson(); err != nil {
		return
	}

	if err = res.ParseResponse(); err != nil {
		return
	}

	if Is2XX(res.HTTP.StatusCode) {
		var ok bool
		var results map[string]interface{}
		if results, ok = res.Results.(map[string]interface{}); !ok {
			err = errors.New("Unexpected response to Recipient List creation (results)")
		} else if id, ok = results["id"].(string); !ok {
			err = errors.New("Unexpected response to Recipient List creation (id)")
		}
	} else {
		err = res.HTTPError()
	}

	return
}

// RecipientLists returns all recipient lists
func (c *Client) RecipientLists() ([]RecipientList, *Response, error) {
	return c.RecipientListsContext(context.Background())
}

// RecipientListsContext is the same as RecipientLists, and it accepts a context.Context
func (c *Client) RecipientListsContext(ctx context.Context) ([]RecipientList, *Response, error) {
	path := fmt.Sprintf(RecipientListsPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	res, err := c.HttpGet(ctx, url)
	if err != nil {
		return nil, nil, err
	}

	var body []byte
	if body, err = res.AssertJson(); err != nil {
		return nil, res, err
	}

	if Is2XX(res.HTTP.StatusCode) {
		rllist := map[string][]RecipientList{}
		if err = json.Unmarshal(body, &rllist); err != nil {
		} else if list, ok := rllist["results"]; ok {
			return list, res, nil
		} else {
			err = errors.New("Unexpected response to RecipientList list")
		}

	} else {
		if err = res.ParseResponse(); err == nil {
			err = res.HTTPError()
		}
	}

	return nil, res, err
}

package gosparkpost

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// https://www.sparkpost.com/api#/reference/recipient-lists
var recipListsPathFormat = "/api/v%d/recipient-lists"

// RecipientList is the JSON structure accepted by and returned from the SparkPost Recipient Lists API.
// It's mostly metadata at this level - see Recipients for more detail.
type RecipientList struct {
	ID          string       `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Attributes  interface{}  `json:"attributes,omitempty"`
	Recipients  *[]Recipient `json:"recipients"`

	Accepted *int `json:"total_accepted_recipients,omitempty"`
}

func (rl *RecipientList) String() string {
	n := 0
	if rl.Recipients != nil {
		n = len(*rl.Recipients)
	} else if rl.Accepted != nil {
		n = *rl.Accepted
	}
	return fmt.Sprintf("ID:\t%s\nName:\t%s\nDesc:\t%s\nCount:\t%d\n",
		rl.ID, rl.Name, rl.Description, n)
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
			err = fmt.Errorf("Recipient.Address may not be empty")
		} else {
			a.Email = addrVal
		}

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
				err = fmt.Errorf("strings are required for all Recipient.Address values")
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
		err = fmt.Errorf("unsupported Recipient.Address value type [%s]", reflect.TypeOf(addrVal))
	}

	return
}

// Validate runs sanity checks on a RecipientList struct. This should
// catch most errors before attempting a doomed API call.
func (rl *RecipientList) Validate() error {
	// enforce required parameters
	if rl.Recipients == nil || len(*rl.Recipients) <= 0 {
		return fmt.Errorf("RecipientList requires at least one Recipient")
	}

	// enforce max lengths
	if len(rl.ID) > 64 {
		return fmt.Errorf("RecipientList id may not be longer than 64 bytes")
	} else if len(rl.Name) > 64 {
		return fmt.Errorf("RecipientList name may not be longer than 64 bytes")
	} else if len(rl.Description) > 1024 {
		return fmt.Errorf("RecipientList description may not be longer than 1024 bytes")
	}

	var err error
	for _, r := range *rl.Recipients {
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

	// Metadata must be an object, not an array or bool etc.
	if r.Metadata != nil {
		err := AssertObject(r.Metadata, "metadata")
		if err != nil {
			return err
		}
	}

	// SubstitutionData must be an object, not an array or bool etc.
	if r.SubstitutionData != nil {
		err := AssertObject(r.SubstitutionData, "substitution_data")
		if err != nil {
			return err
		}
	}

	return nil
}

// Create accepts a populated RecipientList object, validates it,
// and performs an API call against the configured endpoint.
func (c *Client) RecipientListCreate(rl *RecipientList) (id string, res *Response, err error) {
	if rl == nil {
		err = fmt.Errorf("Create called with nil RecipientList")
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

	path := fmt.Sprintf(recipListsPathFormat, c.Config.ApiVersion)
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
		id, ok = res.Results["id"].(string)
		if !ok {
			err = fmt.Errorf("Unexpected response to Recipient List creation")
		}

	} else if len(res.Errors) > 0 {
		// handle common errors
		err = res.PrettyError("RecipientList", "create")
		if err != nil {
			return
		}

		code := res.HTTP.StatusCode
		if code == 400 || code == 422 {
			eobj := res.Errors[0]
			err = fmt.Errorf("%s: %s\n%s", eobj.Code, eobj.Message, eobj.Description)
		} else { // everything else
			err = fmt.Errorf("%d: %s", code, string(res.Body))
		}
	}

	return
}

func (c *Client) RecipientLists() (*[]RecipientList, *Response, error) {
	path := fmt.Sprintf(recipListsPathFormat, c.Config.ApiVersion)
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
		rllist := map[string][]RecipientList{}
		if err = json.Unmarshal(body, &rllist); err != nil {
			return nil, res, err
		} else if list, ok := rllist["results"]; ok {
			return &list, res, nil
		}
		return nil, res, fmt.Errorf("Unexpected response to RecipientList list")

	} else {
		err = res.ParseResponse()
		if err != nil {
			return nil, res, err
		}
		if len(res.Errors) > 0 {
			err = res.PrettyError("RecipientList", "list")
			if err != nil {
				return nil, res, err
			}
		}
		return nil, res, fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return nil, res, err
}

// Package recipients interacts with the SparkPost Recipient Lists API.
// https://www.sparkpost.com/api#/reference/recipient-lists
package recipients

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/SparkPost/go-sparkpost/api"
)

// RecipientLists is your handle for the Recipient Lists API.
type RecipientLists struct{ api.API }

// New gets a RecipientLists object ready to use with the specified config.
func New(cfg *api.Config) (*RecipientLists, error) {
	// FIXME: allow caller to set api version
	rl := &RecipientLists{}
	err := rl.Init(cfg, "/api/v1/recipient-lists")
	if err != nil {
		return nil, err
	}
	return rl, nil
}

// RecipientList is the JSON structure accepted by and returned from the SparkPost Recipient Lists API.
// It's mostly metadata at this level - see Recipients for more detail.
type RecipientList struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Attributes  interface{} `json:"attributes,omitempty"`
	Recipients  []Recipient `json:"recipients"`
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
// It can also be a plain string.
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

// Validate runs sanity checks on a RecipientList struct.
// This should catch most errors before attempting a doomed API call.
func (rl RecipientList) Validate() error {
	// enforce required parameters
	if len(rl.Recipients) <= 0 {
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
	for _, r := range rl.Recipients {
		err = r.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate runs sanity checks on a Recipient struct
// This should catch most errors before attempting a doomed API call.
func (r Recipient) Validate() error {
	_, err := ParseAddress(r.Address)
	if err != nil {
		return err
	}

	// Metadata must be an object, not an array or bool etc.
	if r.Metadata != nil {
		err := api.AssertObject(r.Metadata, "metadata")
		if err != nil {
			return err
		}
	}

	// SubstitutionData must be an object, not an array or bool etc.
	if r.SubstitutionData != nil {
		err := api.AssertObject(r.SubstitutionData, "substitution_data")
		if err != nil {
			return err
		}
	}

	return nil
}

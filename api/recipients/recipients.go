// Package recipients interacts with the SparkPost Recipient Lists API.
// https://www.sparkpost.com/api#/reference/recipient-lists
package recipients

import (
	"fmt"
	"reflect"
	"strings"

	"bitbucket.org/yargevad/go-sparkpost/api"
)

// RecipientLists is your handle for the Recipient Lists API.
type RecipientLists struct {
	api.API
	Path string
}

// New gets a RecipientLists object ready to use with the specified config.
func New(cfg *api.Config) (*RecipientLists, error) {
	rl := &RecipientLists{}
	//err := rl.Init(cfg)
	//if err != nil {
	//	return nil, err
	//}
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
	MetaData         interface{} `json:"metadata,omitempty"`
	SubstitutionData interface{} `json:"substitution_data,omitempty"`
}

type Address struct {
	Email    string `json:"email"`
	Name     string `json:"name,omitempty"`
	HeaderTo string `json:"header_to,omitempty"`
}

func (rl *RecipientLists) Init(cfg *api.Config) error {
	// FIXME: allow caller to set api version
	rl.Path = "/api/v1/recipient-lists"
	return rl.API.Init(cfg)
}

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

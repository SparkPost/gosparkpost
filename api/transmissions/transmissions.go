// Package transmissions interacts with the SparkPost Transmissions API.
// https://www.sparkpost.com/api#/reference/transmissions
package transmissions

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/SparkPost/go-sparkpost/api"
	recipients "github.com/SparkPost/go-sparkpost/api/recipient_lists"
	"github.com/SparkPost/go-sparkpost/api/templates"
)

// Transmissions is your handle for the Transmissions API.
type Transmissions struct{ api.API }

func New(cfg api.Config) (*Transmissions, error) {
	t := &Transmissions{}
	path := fmt.Sprintf("/api/v%d/transmissions", cfg.ApiVersion)
	err := t.Init(cfg, path)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Transmission is the JSON structure accepted by and returned from the SparkPost Transmissions API.
type Transmission struct {
	ID               string      `json:"id,omitempty"`
	State            string      `json:"state,omitempty"`
	Options          *Options    `json:"options,omitempty"`
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
// If not specified, and present in templates.Options, those values will be used.
type Options struct {
	templates.Options

	StartTime       time.Time `json:"start_time,omitempty"`
	Sandbox         bool      `json:"sandbox,omitempty"`
	SkipSuppression bool      `json:"skip_suppression,omitempty"`
}

// ParseRecipients asserts that Transmission.Recipients is valid.
func ParseRecipients(recips interface{}) (err error) {
	switch rVal := recips.(type) {
	case map[string]interface{}:
		for k, v := range rVal {
			switch vVal := v.(type) {
			case string:
				if strings.EqualFold(k, "list_id") {
					return nil
				}
			default:
				return fmt.Errorf("Transmission.Recipient objects must contain string values, not [%s]", reflect.TypeOf(vVal))
			}
		}
		return fmt.Errorf("Transmission.Recipient objects must contain a key `list_id`")

	case map[string]string:
		for k, _ := range rVal {
			if strings.EqualFold(k, "list_id") {
				return nil
			}
		}
		return fmt.Errorf("Transmission.Recipient objects must contain a key `list_id`")

	case []interface{}:
		for _, v := range rVal {
			switch r := v.(type) {
			case recipients.Recipient:
				err = r.Validate()
				if err != nil {
					return err
				}

			default:
				return fmt.Errorf("Failed to parse inline Transmission.Recipient list")
			}
		}

	default:
		return fmt.Errorf("Unsupported Transmission.Recipient type [%s]", reflect.TypeOf(rVal))
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

	case templates.Content:
		te := &templates.Template{Name: "tmp", Content: rVal}
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
	err := ParseRecipients(t.Recipients)
	if err != nil {
		return err
	}
	err = ParseContent(t.Content)
	if err != nil {
		return err
	}

	// Metadata must be an object, not an array or bool etc.
	if t.Metadata != nil {
		err := api.AssertObject(t.Metadata, "metadata")
		if err != nil {
			return err
		}
	}

	// SubstitutionData must be an object, not an array or bool etc.
	if t.SubstitutionData != nil {
		err := api.AssertObject(t.SubstitutionData, "substitution_data")
		if err != nil {
			return err
		}
	}

	return nil
}

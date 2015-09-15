// Package templates interacts with the SparkPost Templates API.
// https://www.sparkpost.com/api#/reference/templates
package templates

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"bitbucket.org/yargevad/go-sparkpost/api"
)

// Templates is your handle for the Templates API.
type Templates struct {
	api.API

	Path     string
	Response *api.Response
	// TODO? move Response into api.API
}

// New gets a Templates object ready to use with the specified config.
func New(cfg *api.Config) (*Templates, error) {
	t := &Templates{}
	err := t.Init(cfg)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Template is the JSON structure understood by the SparkPost API.
// It's mostly metadata at this level - see Content and Options for more detail.
type Template struct {
	ID          string    `json:"id,omitempty"`
	Content     Content   `json:"content,omitempty"`
	Published   bool      `json:"published,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	LastUse     time.Time `json:"last_use,omitempty"`
	LastUpdate  time.Time `json:"last_update_time,omitempty"`
	Options     *Options  `json:"options,omitempty"`
}

// Content is what you'll send to your Recipients.
// Knowledge of SparkPost's substitution/templating capabilities will come in handy here.
// https://www.sparkpost.com/api#/introduction/substitutions-reference
type Content struct {
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	Subject     string            `json:"subject,omitempty"`
	From        interface{}       `json:"from,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	EmailRFC822 string            `json:"email_rfc822,omitempty"`
}

// From describes the nested object way of specifying the From header.
// It can also be a plain string.
type From struct {
	Email string
	Name  string
}

// Options specifies settings to apply to this template.
// These settings may be overridden in the Transmission API call.
type Options struct {
	OpenTracking  bool `json:"open_tracking,omitempty"`
	ClickTracking bool `json:"click_tracking,omitempty"`
	Transactional bool `json:"transactional,omitempty"`
}

// Init sets the path part of the API URL and initializes the embedded API object.
func (t *Templates) Init(cfg *api.Config) (err error) {
	// FIXME: allow specifying api version
	t.Path = "/api/v1/templates"
	return t.API.Init(cfg)
}

// ParseFrom parses the various allowable Content.From values.
func ParseFrom(from interface{}) (f From, err error) {
	// handle both of the allowed types
	switch fromVal := from.(type) {
	case string: // simple string value
		if fromVal == "" {
			err = fmt.Errorf("Content.From may not be empty")
		} else {
			f.Email = fromVal
		}

	case map[string]interface{}:
		// auto-parsed nested json object
		for k, v := range fromVal {
			switch vVal := v.(type) {
			case string:
				if strings.EqualFold(k, "name") {
					f.Name = vVal
				} else if strings.EqualFold(k, "email") {
					f.Email = vVal
				}
			default:
				err = fmt.Errorf("strings are required for all Content.From values")
				break
			}
		}

	case map[string]string:
		// user-provided json literal (convenience)
		for k, v := range fromVal {
			if strings.EqualFold(k, "name") {
				f.Name = v
			} else if strings.EqualFold(k, "email") {
				f.Email = v
			}
		}

	default:
		err = fmt.Errorf("unsupported Content.From value type [%s]", reflect.TypeOf(fromVal))
	}

	return
}

// Validate runs sanity checks on a Template struct.
// This should catch most errors before attempting a doomed API call.
func (t Template) Validate() error {
	if t.Content.EmailRFC822 != "" {
		// TODO: optionally validate MIME structure
		// if MIME content is present, clobber all other Content options
		t.Content = Content{EmailRFC822: t.Content.EmailRFC822}
		return nil
	}

	// enforce required parameters
	if t.Name == "" && t.ID == "" {
		return fmt.Errorf("Template requires a name or id")
	} else if t.Content.Subject == "" {
		return fmt.Errorf("Template requires a non-empty Content.Subject")
	} else if t.Content.HTML == "" && t.Content.Text == "" {
		return fmt.Errorf("Template requires either Content.HTML or Content.Text")
	}
	_, err := ParseFrom(t.Content.From)
	if err != nil {
		return err
	}

	// enforce max lengths
	// TODO: enforce 15MB Content limit
	if len(t.ID) > 64 {
		return fmt.Errorf("Template id may not be longer than 64 bytes")
	} else if len(t.Name) > 1024 {
		return fmt.Errorf("Template name may not be longer than 1024 bytes")
	} else if len(t.Description) > 1024 {
		return fmt.Errorf("Template description may not be longer than 1024 bytes")
	}

	return nil
}

// Create is a convenience function implementing the "at a minimum" case mentioned in the
// SparkPost API docs.
// https://www.sparkpost.com/api#/reference/templates/create-and-list/create-a-template
func (t Templates) Create(name string, content Content) (string, error) {
	template := &Template{
		Name:    name,
		Content: content,
	}
	return t.CreateWithTemplate(template)
}

// CreateWithJSON accepts JSON and calls CreateWithTemplate if parsing succeeds.
func (t Templates) CreateWithJSON(j string) (string, error) {
	template := &Template{}
	err := json.Unmarshal([]byte(j), template)
	if err != nil {
		return "", err
	}
	return t.CreateWithTemplate(template)
}

// CreateWithTemplate accepts a populated Template object, validates its Contents,
// and performs an API call against the configured endpoint.
// Helper functions call into this function after building a Template object.
func (t Templates) CreateWithTemplate(template *Template) (id string, err error) {
	if template == nil {
		err = fmt.Errorf("CreateWithTemplate called with nil Template")
		return
	}

	err = template.Validate()
	if err != nil {
		return
	}

	jsonBytes, err := json.Marshal(template)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s%s", t.API.Config.BaseUrl, t.Path)
	res, err := t.HttpPost(url, jsonBytes)
	if err != nil {
		return
	}

	if err = api.AssertJson(res); err != nil {
		return
	}

	t.Response, err = api.ParseApiResponse(res)
	if err != nil {
		return
	}

	if res.StatusCode == 200 {
		var ok bool
		id, ok = t.Response.Results["id"]
		if !ok {
			err = fmt.Errorf("Unexpected response to Template creation")
		}

	} else if len(t.Response.Errors) > 0 {
		// handle common errors
		err = api.PrettyError("Template", "create", res)
		if err != nil {
			return
		}

		if res.StatusCode == 422 { // template syntax error
			err = fmt.Errorf(t.Response.Errors[0].Description)
		} else { // everything else
			err = fmt.Errorf("%d: %s", res.StatusCode, t.Response.Body)
		}
	}

	return
}

// List returns metadata for all Templates in the system.
func (t Templates) List() ([]Template, error) {
	url := fmt.Sprintf("%s%s", t.API.Config.BaseUrl, t.Path)
	res, err := t.HttpGet(url)
	if err != nil {
		return nil, err
	}

	if err = api.AssertJson(res); err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var body []byte
		body, err = api.ReadBody(res)
		if err != nil {
			return nil, err
		}
		tlist := map[string][]Template{}
		if err = json.Unmarshal(body, &tlist); err != nil {
			return nil, err
		} else if list, ok := tlist["results"]; ok {
			return list, nil
		}
		return nil, fmt.Errorf("Unexpected response to Template list")

	} else {
		t.Response, err = api.ParseApiResponse(res)
		if err != nil {
			return nil, err
		}
	}

	return nil, err
}

// Delete removes the Template with the specified id.
func (t Templates) Delete(id string) (err error) {
	if id == "" {
		err = fmt.Errorf("Delete called with blank id")
		return
	}

	url := fmt.Sprintf("%s%s/%s", t.API.Config.BaseUrl, t.Path, id)
	res, err := t.HttpDelete(url)
	if err != nil {
		return
	}

	if err = api.AssertJson(res); err != nil {
		return
	}

	t.Response, err = api.ParseApiResponse(res)
	if err != nil {
		return
	}

	if res.StatusCode == 200 {
		return nil

	} else if len(t.Response.Errors) > 0 {
		// handle common errors
		err = api.PrettyError("Template", "delete", res)
		if err != nil {
			return
		}

		// handle template-specific ones
		if res.StatusCode == 409 {
			err = fmt.Errorf("Template with id [%s] is in use by msg generation", id)
		} else { // everything else
			err = fmt.Errorf("%d: %s", res.StatusCode, t.Response.Body)
		}
	}

	return
}

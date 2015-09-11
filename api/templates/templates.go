package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"bitbucket.org/yargevad/go-sparkpost/api"
)

type Templates struct {
	api.API

	Path     string
	Response *api.Response
}

func New(cfg *api.Config) (*Templates, error) {
	t := &Templates{}
	err := t.Init(cfg)
	if err != nil {
		return nil, err
	}
	return t, nil
}

type Template struct {
	ID          string   `json:"id,omitempty"`
	Content     Content  `json:"content,omitempty"`
	Published   bool     `json:"published,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Options     *Options `json:"options,omitempty"`
}

type Content struct {
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	Subject     string            `json:"subject,omitempty"`
	From        interface{}       `json:"from,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	EmailRFC822 string            `json:"email_rfc822,omitempty"`
}

type From struct {
	Email string
	Name  string
}

type Options struct {
	OpenTracking  bool `json:"open_tracking,omitempty"`
	ClickTracking bool `json:"click_tracking,omitempty"`
	Transactional bool `json:"transactional,omitempty"`
}

func (t *Templates) Init(cfg *api.Config) (err error) {
	// FIXME: allow specifying api version
	t.Path = "/api/v1/templates"
	return t.API.Init(cfg)
}

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

// Helper function for the "at a minimum" case mentioned in the SparkPost API docs:
// https://www.sparkpost.com/api#/reference/templates/create-and-list/create-a-template
func (t Templates) Create(name string, content Content) (id string, err error) {
	template := &Template{
		Name:    name,
		Content: content,
	}
	err = template.Validate()
	if err != nil {
		return
	}

	id, err = t.CreateWithTemplate(template)
	return
}

// Support for all Template API options.
// Helper functions call into this function after building a Template object.
// Validates input before making request.
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
	res, err := t.Post(url, jsonBytes)
	if err != nil {
		return
	}

	if err = api.AssertJson(res); err != nil {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	apiResponse := &api.Response{}
	apiResponse.HTTP = res
	// pull out either "results" or "errors"
	err = json.Unmarshal(body, apiResponse)
	if err != nil {
		return
	}
	// keep this around so callers can see failure details
	t.Response = apiResponse

	if len(apiResponse.Errors) > 0 {
		err = fmt.Errorf("Error creating Template:\n%s", string(body))
	} else if len(apiResponse.Results) > 0 {
		var ok bool
		id, ok = apiResponse.Results["id"]
		if !ok {
			err = fmt.Errorf("Unexpected response to Template creation:\n%s", string(body))
		}
	}

	return
}

func (t Templates) CreateWithJSON(j string) (id string, err error) {
	template := &Template{}
	err = json.Unmarshal([]byte(j), template)
	if err != nil {
		return
	}

	id, err = t.CreateWithTemplate(template)
	return
}

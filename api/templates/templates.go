package templates

import (
	"fmt"
	"reflect"
	"strings"

	"bitbucket.org/yargevad/go-sparkpost/api"
)

type Templates struct {
	*api.API
	Path   string
	Errors []api.Error
}

type Template struct {
	ID          string  `json:"id,omitempty"`
	Content     Content `json:"content,omitempty"`
	Published   bool    `json:"published,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Options     Options `json:"options,omitempty"`
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

func ParseFrom(from interface{}) (f From, err error) {
	// handle both of the allowed types
	switch fromVal := from.(type) {
	case string: // simple string value
		if fromVal == "" {
			err = fmt.Errorf("Content.From may not be empty")
		} else {
			f.Email = fromVal
		}

	case map[string]interface{}: // nested json object
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

	default:
		err = fmt.Errorf("unsupported Content.From value type [%s]", reflect.TypeOf(fromVal))
	}

	return
}

// Helper function for the "at a minimum" case mentioned in the SparkPost API docs:
// https://www.sparkpost.com/api#/reference/templates/create-and-list/create-a-template
func (t Templates) Create(name string, content Content) (id string, err error) {
	if name == "" {
		err = fmt.Errorf("templates.Create requires a name")
		return
	} else if content.Subject == "" {
		err = fmt.Errorf("templates.Create requires a non-empty Content.Subject")
		return
	} else if content.HTML == "" && content.Text == "" {
		err = fmt.Errorf("templates.Create requires either Content.HTML or Content.Text")
		return
	}
	_, err = ParseFrom(content.From)
	if err != nil {
		return
	}

	template := Template{
		Name:    name,
		Content: content,
	}
	id, err = t.CreateWithTemplate(template)

	return
}

// Support for all Template API options.
// Helper functions call into this function after building a Template object.
// Validates input before making request.
func (t Templates) CreateWithTemplate(template Template) (id string, err error) {
	return
}

func (t Templates) CreateWithJSON(j string) (id string, err error) {
	return
}

// Package templates interacts with the SparkPost Templates API.
// https://www.sparkpost.com/api#/reference/templates
package templates

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/SparkPost/go-sparkpost/api"
)

// Templates is your handle for the Templates API.
type Templates struct{ api.API }

// New gets a Templates object ready to use with the specified config.
func New(cfg api.Config) (*Templates, error) {
	t := &Templates{}
	path := fmt.Sprintf("/api/v%d/templates", cfg.ApiVersion)
	err := t.Init(cfg, path)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Template is the JSON structure accepted by and returned from the SparkPost Templates API.
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

// Options specifies settings to apply to this Template.
// These settings may be overridden in the Transmission API call.
type Options struct {
	OpenTracking  *bool `json:"open_tracking,omitempty"`
	ClickTracking *bool `json:"click_tracking,omitempty"`
	Transactional *bool `json:"transactional,omitempty"`
}

// ParseFrom parses the various allowable Content.From values.
func ParseFrom(from interface{}) (f From, err error) {
	// handle the allowed types
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
func (t *Template) Validate() error {
	if t == nil {
		return fmt.Errorf("Can't Validate a nil Template")
	}

	if t.Content.EmailRFC822 != "" {
		// TODO: optionally validate MIME structure
		// if MIME content is present, clobber all other Content options
		t.Content = Content{EmailRFC822: t.Content.EmailRFC822}
		return nil
	}

	// enforce required parameters
	if t.Content.Subject == "" {
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

// SetHeaders is a convenience method which sets Template.Content.Headers to the provided map.
func (t *Template) SetHeaders(headers map[string]string) {
	t.Content.Headers = headers
}

// Build accepts a map of key/value pairs, builds, and returns a Template object
// suitable for use with Create. Providing an unsupported map key will result in
// an error. Custom headers must be added to Template.Content.Headers by the caller.
//
// The expected map keys are:
//   id: ID used to reference the template
//   name: Editable display name
//   description: Detailed description of the template
//   published: Defaults to False. Whether the template is a published or draft version
//   track_opens: Defaults to transmission level setting. Used to track opens of transmission
//   track_clicks: Defaults to transmission level setting. Used to track clicks of transmission
//   is_transactional: Defaults to transmission level setting.
//              Distinguishes between transactional and non-transactional messages
//              for unsubscribe and suppression purposes
//   html: HTML part of template
//   text: Text part of template
//   subject: Subject of template
//   from_email: Email portion of From header
//   from_name: Name portion of From header
//   reply_to: Reply to of template
func (t Templates) Build(p map[string]string) (*Template, error) {
	T := &Template{}

	// Look up expected keys in the map, deleting as we find them.
	if id, ok := p["id"]; ok {
		T.ID = id
		delete(p, "id")
	}
	if name, ok := p["name"]; ok {
		T.Name = name
		delete(p, "name")
	}
	if desc, ok := p["description"]; ok {
		T.Description = desc
		delete(p, "description")
	}
	if pub, ok := p["published"]; ok {
		if strings.EqualFold(pub, "true") {
			T.Published = true
		} else {
			T.Published = false
		}
		delete(p, "published")
	}

	if opens, ok := p["track_opens"]; ok {
		if T.Options == nil {
			T.Options = new(Options)
		}
		if strings.EqualFold(opens, "true") {
			*T.Options.OpenTracking = true
		} else {
			*T.Options.OpenTracking = false
		}
		delete(p, "track_opens")
	}

	if clicks, ok := p["track_clicks"]; ok {
		if T.Options == nil {
			T.Options = new(Options)
		}
		if strings.EqualFold(clicks, "true") {
			*T.Options.ClickTracking = true
		} else {
			*T.Options.ClickTracking = false
		}
		delete(p, "track_clicks")
	}

	if isTransactional, ok := p["is_transactional"]; ok {
		if T.Options == nil {
			T.Options = new(Options)
		}
		if strings.EqualFold(isTransactional, "true") {
			*T.Options.ClickTracking = true
		} else {
			*T.Options.ClickTracking = false
		}
		delete(p, "is_transactional")
	}

	if html, ok := p["html"]; ok {
		T.Content.HTML = html
		delete(p, "html")
	}
	if text, ok := p["text"]; ok {
		T.Content.Text = text
		delete(p, "text")
	}
	if subject, ok := p["subject"]; ok {
		T.Content.Subject = subject
		delete(p, "subject")
	}
	if replyTo, ok := p["reply_to"]; ok {
		T.Content.ReplyTo = replyTo
		delete(p, "reply_to")
	}

	if email, ok := p["from_email"]; ok {
		if T.Content.From == nil {
			T.Content.From = From{}
		}
		switch from := T.Content.From.(type) {
		case From:
			from.Email = email
			delete(p, "from_email")
		default:
			return nil, fmt.Errorf("Expected type `From`, got [%s].", reflect.TypeOf(from))
		}
	}

	if name, ok := p["from_name"]; ok {
		if T.Content.From == nil {
			T.Content.From = From{}
		}
		switch from := T.Content.From.(type) {
		case From:
			from.Name = name
			delete(p, "from_name")
		default:
			return nil, fmt.Errorf("Expected type `From`, got [%s].", reflect.TypeOf(from))
		}
	}

	// If there are any keys left, they are unsupported.
	if len(p) > 0 {
		return nil, fmt.Errorf("Build received unsupported keys")
	}
	return T, nil
}

// Create accepts a populated Template object, validates its Contents,
// and performs an API call against the configured endpoint.
func (t Templates) Create(template *Template) (id string, err error) {
	if template == nil {
		err = fmt.Errorf("Create called with nil Template")
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

	url := fmt.Sprintf("%s%s", t.Config.BaseUrl, t.Path)
	res, err := t.HttpPost(url, jsonBytes)
	if err != nil {
		return
	}

	if err = api.AssertJson(res); err != nil {
		return
	}

	err = t.ParseResponse(res)
	if err != nil {
		return
	}

	if res.StatusCode == 200 {
		var ok bool
		id, ok = t.Response.Results["id"].(string)
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
			eobj := t.Response.Errors[0]
			err = fmt.Errorf("%s: %s\n%s", eobj.Code, eobj.Message, eobj.Description)
		} else { // everything else
			err = fmt.Errorf("%d: %s", res.StatusCode, string(t.Response.Body))
		}
	}

	return
}

// List returns metadata for all Templates in the system.
func (t Templates) List() ([]Template, error) {
	url := fmt.Sprintf("%s%s", t.Config.BaseUrl, t.Path)
	res, err := t.HttpGet(url)
	if err != nil {
		return nil, err
	}

	if err = api.AssertJson(res); err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var body []byte
		body, err = t.ReadBody(res)
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
		err = t.ParseResponse(res)
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

	url := fmt.Sprintf("%s%s/%s", t.Config.BaseUrl, t.Path, id)
	res, err := t.HttpDelete(url)
	if err != nil {
		return
	}

	if err = api.AssertJson(res); err != nil {
		return
	}

	err = t.ParseResponse(res)
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
			err = fmt.Errorf("%d: %s", res.StatusCode, string(t.Response.Body))
		}
	}

	return
}

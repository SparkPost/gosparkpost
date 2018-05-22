package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// https://www.sparkpost.com/api#/reference/templates
var TemplatesPathFormat = "/api/v%d/templates"

// Template is the JSON structure accepted by and returned from the SparkPost Templates API.
// It's mostly metadata at this level - see Content and Options for more detail.
type Template struct {
	ID          string       `json:"id,omitempty"`
	Content     Content      `json:"content,omitempty"`
	Published   bool         `json:"published,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	LastUse     time.Time    `json:"last_use,omitempty"`
	LastUpdate  time.Time    `json:"last_update_time,omitempty"`
	Options     *TmplOptions `json:"options,omitempty"`
}

// Content is what you'll send to your Recipients.
// Knowledge of SparkPost's substitution/templating capabilities will come in handy here.
// https://www.sparkpost.com/api#/introduction/substitutions-reference
type Content struct {
	HTML         string            `json:"html,omitempty"`
	Text         string            `json:"text,omitempty"`
	Subject      string            `json:"subject,omitempty"`
	From         interface{}       `json:"from,omitempty"`
	ReplyTo      string            `json:"reply_to,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	EmailRFC822  string            `json:"email_rfc822,omitempty"`
	Attachments  []Attachment      `json:"attachments,omitempty"`
	InlineImages []InlineImage     `json:"inline_images,omitempty"`
}

// Attachment contains metadata and the contents of the file to attach.
type Attachment struct {
	MIMEType string `json:"type"`
	Filename string `json:"name"`
	B64Data  string `json:"data"`
}

// InlineImage contains metadata and the contents of the image to make available for inline use.
type InlineImage Attachment

// From describes the nested object way of specifying the From header.
// Content.From can be specified this way, or as a plain string.
type From struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// TmplOptions specifies settings to apply to this Template.
// These settings may be overridden in the Transmission API call.
type TmplOptions struct {
	OpenTracking  *bool `json:"open_tracking,omitempty"`
	ClickTracking *bool `json:"click_tracking,omitempty"`
	Transactional *bool `json:"transactional,omitempty"`
}

// PreviewOptions contains the required subsitution_data object to
// preview a template
type PreviewOptions struct {
	SubstitutionData map[string]interface{} `json:"substitution_data"`
}

// ParseFrom parses the various allowable Content.From values.
func ParseFrom(from interface{}) (f From, err error) {
	// handle the allowed types
	switch fromVal := from.(type) {
	case From:
		f = fromVal

	case Address:
		f.Email = fromVal.Email
		f.Name = fromVal.Name

	case string: // simple string value
		if fromVal == "" {
			err = errors.New("Content.From may not be empty")
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
				err = errors.New("strings are required for all Content.From values")
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
		err = errors.Errorf("unsupported Content.From value type [%T]", fromVal)
	}

	return
}

// Validate runs sanity checks on a Template struct.
// This should catch most errors before attempting a doomed API call.
func (t *Template) Validate() error {
	if t == nil {
		return errors.New("Can't Validate a nil Template")
	}

	if t.Content.EmailRFC822 != "" {
		// TODO: optionally validate MIME structure
		// if MIME content is present, clobber all other Content options
		t.Content = Content{EmailRFC822: t.Content.EmailRFC822}
		return nil
	}

	// enforce required parameters
	if t.Content.Subject == "" {
		return errors.New("Template requires a non-empty Content.Subject")
	} else if t.Content.HTML == "" && t.Content.Text == "" {
		return errors.New("Template requires either Content.HTML or Content.Text")
	}
	_, err := ParseFrom(t.Content.From)
	if err != nil {
		return err
	}

	if len(t.Content.Attachments) > 0 {
		for _, att := range t.Content.Attachments {
			if len(att.Filename) > 255 {
				return errors.Errorf("Attachment name length must be <= 255: [%s]", att.Filename)
			} else if strings.ContainsAny(att.B64Data, "\r\n") {
				return errors.New("Attachment data may not contain line breaks [\\r\\n]")
			}
		}
	}

	if len(t.Content.InlineImages) > 0 {
		for _, img := range t.Content.InlineImages {
			if len(img.Filename) > 255 {
				return errors.Errorf("InlineImage name length must be <= 255: [%s]", img.Filename)
			} else if strings.ContainsAny(img.B64Data, "\r\n") {
				return errors.New("InlineImage data may not contain line breaks [\\r\\n]")
			}
		}
	}

	// enforce max lengths
	if len(t.ID) > 64 {
		return errors.New("Template id may not be longer than 64 bytes")
	} else if len(t.Name) > 1024 {
		return errors.New("Template name may not be longer than 1024 bytes")
	} else if len(t.Description) > 1024 {
		return errors.New("Template description may not be longer than 1024 bytes")
	}

	return nil
}

// TemplateCreate accepts a populated Template object, validates its Contents,
// and performs an API call against the configured endpoint.
func (c *Client) TemplateCreate(t *Template) (id string, res *Response, err error) {
	return c.TemplateCreateContext(context.Background(), t)
}

// TemplateCreateContext is the same as TemplateCreate, and it allows the caller to provide a context.
func (c *Client) TemplateCreateContext(ctx context.Context, t *Template) (id string, res *Response, err error) {
	if t == nil {
		err = errors.New("Create called with nil Template")
		return
	}

	err = t.Validate()
	if err != nil {
		return
	}

	// A Template that makes it past Validate() will always Marshal
	jsonBytes, _ := json.Marshal(t)

	path := fmt.Sprintf(TemplatesPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	res, err = c.HttpPost(ctx, url, jsonBytes)
	if err != nil {
		return
	}

	if err = res.ParseResponse(); err != nil {
		return
	}

	if Is2XX(res.HTTP.StatusCode) {
		var ok bool
		var results map[string]interface{}
		if results, ok = res.Results.(map[string]interface{}); !ok {
			err = errors.New("Unexpected response to Template creation (results)")
		} else if id, ok = results["id"].(string); !ok {
			err = errors.New("Unexpected response to Template creation (id)")
		}
	} else {
		err = res.HTTPError()
	}
	return
}

// TemplateGet fills out the provided template, using the specified id.
func (c *Client) TemplateGet(t *Template, draft bool) (*Response, error) {
	return c.TemplateGetContext(context.Background(), t, draft)
}

// TemplateGetContext is the same as TemplateGet, and it allows the caller to provide a context
func (c *Client) TemplateGetContext(ctx context.Context, t *Template, draft bool) (*Response, error) {
	if t == nil {
		return nil, errors.New("TemplateGet called with nil Template")
	}

	if t.ID == "" {
		return nil, errors.New("TemplateGet called with blank id")
	}

	path := fmt.Sprintf(TemplatesPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s/%s?draft=%t", c.Config.BaseUrl, path, t.ID, draft)

	res, err := c.HttpGet(ctx, url)
	if err != nil {
		return nil, err
	}

	var body []byte
	body, err = res.ReadBody()
	if err != nil {
		return res, err
	}

	if err = res.ParseResponse(); err != nil {
		return res, err
	}

	if Is2XX(res.HTTP.StatusCode) {
		// Unwrap the returned Template
		tmp := map[string]*json.RawMessage{}
		if err = json.Unmarshal(body, &tmp); err != nil {
		} else if results, ok := tmp["results"]; ok {
			err = json.Unmarshal(*results, t)
		} else {
			err = errors.New("Unexpected response to TemplateGet")
		}
	} else {
		err = res.HTTPError()
	}

	return res, err
}

// TemplateUpdate updates a draft/published template with the specified id
// The `updatePublished` parameter controls which version (draft/false vs published/true) of the template will be updated.
func (c *Client) TemplateUpdate(t *Template, updatePublished bool) (res *Response, err error) {
	return c.TemplateUpdateContext(context.Background(), t, updatePublished)
}

// TemplateUpdateContext is the same as TemplateUpdate, and it allows the caller to provide a context
func (c *Client) TemplateUpdateContext(ctx context.Context, t *Template, updatePublished bool) (res *Response, err error) {
	if t == nil {
		err = errors.New("Update called with nil Template")
		return
	}

	if t.ID == "" {
		err = errors.New("Update called with blank id")
		return
	}

	err = t.Validate()
	if err != nil {
		return
	}

	// A Template that makes it past Validate() will always Marshal
	jsonBytes, _ := json.Marshal(t)

	path := fmt.Sprintf(TemplatesPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s/%s?update_published=%t", c.Config.BaseUrl, path, t.ID, updatePublished)

	return c.HttpPutJson(ctx, url, jsonBytes)
}

// Templates returns metadata for all Templates in the system.
func (c *Client) Templates() ([]Template, *Response, error) {
	return c.TemplatesContext(context.Background())
}

// TemplatesContext is the same as Templates, and it allows the caller to provide a context
func (c *Client) TemplatesContext(ctx context.Context) (tl []Template, res *Response, err error) {
	path := fmt.Sprintf(TemplatesPathFormat, c.Config.ApiVersion)
	url := c.Config.BaseUrl + path
	res, err = c.HttpGet(ctx, url)
	if err != nil {
		return
	}

	var body []byte
	if body, err = res.AssertJson(); err != nil {
		return
	}

	if Is2XX(res.HTTP.StatusCode) {
		tlist := map[string][]Template{}
		if err = json.Unmarshal(body, &tlist); err != nil {
			return
		}
		return tlist["results"], res, nil
	}

	if err = res.ParseResponse(); err == nil {
		err = res.HTTPError()
	}

	return
}

// TemplateDelete removes the Template with the specified id.
func (c *Client) TemplateDelete(id string) (res *Response, err error) {
	return c.TemplateDeleteContext(context.Background(), id)
}

// TemplateDeleteContext is the same as TemplateDelete, and it allows the caller to provide a context
func (c *Client) TemplateDeleteContext(ctx context.Context, id string) (res *Response, err error) {
	if id == "" {
		err = errors.New("Delete called with blank id")
		return
	}

	path := fmt.Sprintf(TemplatesPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, id)
	res, err = c.HttpDelete(ctx, url)
	if err != nil {
		return
	}

	if _, err = res.AssertJson(); err != nil {
		return
	}

	if err = res.ParseResponse(); err == nil {
		err = res.HTTPError()
	}

	return
}

// TemplatePreview renders and returns the output of a template using the provided substitution data.
func (c *Client) TemplatePreview(id string, payload *PreviewOptions) (res *Response, err error) {
	return c.TemplatePreviewContext(context.Background(), id, payload)
}

// TemplatePreviewContext is the same as TemplatePreview, and it allows the caller to provide a context
func (c *Client) TemplatePreviewContext(ctx context.Context, id string, payload *PreviewOptions) (res *Response, err error) {
	if id == "" {
		err = errors.New("Preview called with blank id")
		return
	}

	if payload == nil {
		payload = &PreviewOptions{}
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	path := fmt.Sprintf(TemplatesPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s/%s/preview", c.Config.BaseUrl, path, id)
	res, err = c.HttpPost(ctx, url, jsonBytes)
	if err != nil {
		return
	}

	if _, err = res.AssertJson(); err != nil {
		return
	}

	if err = res.ParseResponse(); err == nil {
		err = res.HTTPError()
	}

	return
}

// TemplatePublish publishes a draft template
func (c *Client) TemplatePublish(id string) (res *Response, err error) {
	return c.TemplatePublishContext(context.Background(), id)
}

// TemplatePublishContext is the same as TemplatePublish, and it allows the caller to provide a context
func (c *Client) TemplatePublishContext(ctx context.Context, id string) (res *Response, err error) {
	if id == "" {
		err = errors.New("Publish called with blank id")
		return
	}

	// A Template that makes it past Validate() will always Marshal
	jsonBytes, _ := json.Marshal(map[string]bool{
		"published": true,
	})

	path := fmt.Sprintf(TemplatesPathFormat, c.Config.ApiVersion)
	url := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, id)

	return c.HttpPutJson(ctx, url, jsonBytes)
}

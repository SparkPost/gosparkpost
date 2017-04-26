package gosparkpost_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

// ExampleTemplate builds a native Go Template structure from a JSON string
func ExampleTemplate() {
	template := &sp.Template{}
	jsonStr := `{
		"name": "testy template",
		"content": {
			"html": "this is a <b>test</b> email!",
			"subject": "test email",
			"from": {
				"name": "tester",
				"email": "tester@example.com"
			},
			"reply_to": "tester@example.com"
		}
	}`
	err := json.Unmarshal([]byte(jsonStr), template)
	if err != nil {
		panic(err)
	}
}

func TestTemplateFromValidation(t *testing.T) {
	for idx, test := range []struct {
		in  interface{}
		err error
		out sp.From
	}{
		{sp.From{"a@b.com", "A B"}, nil, sp.From{"a@b.com", "A B"}},
		{sp.Address{"a@b.com", "A B", "c@d.com"}, nil, sp.From{"a@b.com", "A B"}},
		{"a@b.com", nil, sp.From{"a@b.com", ""}},
		{nil, errors.New("unsupported Content.From value type [%!s(<nil>)]"), sp.From{"", ""}},
		{[]byte("a@b.com"), errors.New("unsupported Content.From value type [[]uint8]"), sp.From{"", ""}},
		{"", errors.New("Content.From may not be empty"), sp.From{"", ""}},
		{map[string]interface{}{"name": "A B", "email": "a@b.com"}, nil, sp.From{"a@b.com", "A B"}},
		{map[string]interface{}{"name": 1, "email": "a@b.com"}, errors.New("strings are required for all Content.From values"),
			sp.From{"a@b.com", ""}},
		{map[string]string{"name": "A B", "email": "a@b.com"}, nil, sp.From{"a@b.com", "A B"}},
	} {
		f, err := sp.ParseFrom(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("ParseFrom[%d] => err %q, want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("ParseFrom[%d] => err %q, want %q", idx, err, test.err)
		} else if f.Email != test.out.Email {
			t.Errorf("ParseFrom[%d] => Email %q, want %q", idx, f.Email, test.out.Email)
		} else if f.Name != test.out.Name {
			t.Errorf("ParseFrom[%d] => Name %q, want %q", idx, f.Name, test.out.Name)
		}
	}
}

// Assert that options are actually ... optional,
// and that unspecified options don't default to their zero values.
func TestTemplateOptions(t *testing.T) {
	var jsonb []byte
	var err error
	var opt bool

	te := &sp.Template{}
	to := &sp.TmplOptions{Transactional: &opt}
	te.Options = to

	jsonb, err = json.Marshal(te)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(jsonb, []byte(`"options":{"transactional":false}`)) {
		t.Fatal("expected transactional option to be false")
	}

	opt = true
	jsonb, err = json.Marshal(te)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(jsonb, []byte(`"options":{"transactional":true}`)) {
		t.Fatalf("expected transactional option to be true:\n%s", string(jsonb))
	}
}

func TestTemplateValidation(t *testing.T) {
	for idx, test := range []struct {
		in  *sp.Template
		err error
		out *sp.Template
	}{
		{nil, errors.New("Can't Validate a nil Template"), nil},
		{&sp.Template{}, errors.New("Template requires a non-empty Content.Subject"), nil},
		{&sp.Template{Content: sp.Content{Subject: "s"}}, errors.New("Template requires either Content.HTML or Content.Text"), nil},
		{&sp.Template{Content: sp.Content{Subject: "s", HTML: "h", From: ""}},
			errors.New("Content.From may not be empty"), nil},

		{&sp.Template{ID: strings.Repeat("id", 33), Content: sp.Content{Subject: "s", HTML: "h", From: "f"}},
			errors.New("Template id may not be longer than 64 bytes"), nil},
		{&sp.Template{Name: strings.Repeat("name", 257), Content: sp.Content{Subject: "s", HTML: "h", From: "f"}},
			errors.New("Template name may not be longer than 1024 bytes"), nil},
		{&sp.Template{Description: strings.Repeat("desc", 257), Content: sp.Content{Subject: "s", HTML: "h", From: "f"}},
			errors.New("Template description may not be longer than 1024 bytes"), nil},

		{&sp.Template{
			Content: sp.Content{
				Subject: "s", HTML: "h", From: "f",
				Attachments: []sp.Attachment{{Filename: strings.Repeat("f", 256)}},
			}},
			errors.Errorf("Attachment name length must be <= 255: [%s]", strings.Repeat("f", 256)), nil},
		{&sp.Template{
			Content: sp.Content{
				Subject: "s", HTML: "h", From: "f",
				Attachments: []sp.Attachment{{B64Data: "\r\n"}},
			}},
			errors.New("Attachment data may not contain line breaks [\\r\\n]"), nil},

		{&sp.Template{
			Content: sp.Content{
				Subject: "s", HTML: "h", From: "f",
				InlineImages: []sp.InlineImage{{Filename: strings.Repeat("f", 256)}},
			}},
			errors.Errorf("InlineImage name length must be <= 255: [%s]", strings.Repeat("f", 256)), nil},
		{&sp.Template{
			Content: sp.Content{
				Subject: "s", HTML: "h", From: "f",
				InlineImages: []sp.InlineImage{{B64Data: "\r\n"}},
			}},
			errors.New("InlineImage data may not contain line breaks [\\r\\n]"), nil},

		{&sp.Template{Content: sp.Content{EmailRFC822: "From:foo@example.com\r\n", Subject: "removeme"}},
			nil, &sp.Template{Content: sp.Content{EmailRFC822: "From:foo@example.com\r\n"}}},
	} {
		err := test.in.Validate()
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("Template.Validate[%d] => err %q, want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("Template.Validate[%d] => err %q, want %q", idx, err, test.err)
		} else if test.out != nil && !reflect.DeepEqual(test.in, test.out) {
			t.Errorf("Template.Validate[%d] => failed post-condition check for %q", test.in)
		}
	}
}

func TestTemplateCreate(t *testing.T) {
	for idx, test := range []struct {
		in     *sp.Template
		err    error
		status int
		json   string
		id     string
	}{
		{nil, errors.New("Create called with nil Template"), 0, "", ""},
		{&sp.Template{}, errors.New("Template requires a non-empty Content.Subject"), 0, "", ""},
		{&sp.Template{Content: sp.Content{Subject: "s", HTML: "h", From: "f"}},
			errors.New("Unexpected response to Template creation (results)"),
			200, `{"foo":{"id":"new-template"}}`, ""},
		{&sp.Template{Content: sp.Content{Subject: "s", HTML: "h", From: "f"}},
			errors.New("Unexpected response to Template creation (id)"),
			200, `{"results":{"ID":"new-template"}}`, ""},
		{&sp.Template{Content: sp.Content{Subject: "s", HTML: "h", From: "f"}},
			errors.New("parsing api response: unexpected end of JSON input"),
			200, `{"truncated":{}`, ""},

		{&sp.Template{Content: sp.Content{Subject: "s{{", HTML: "h", From: "f"}},
			sp.SPErrors([]sp.SPError{{
				Message:     "substitution language syntax error in template content",
				Description: "Error while compiling header Subject: substitution statement missing ending '}}'",
				Code:        "3000",
				Part:        "Header:Subject",
			}}),
			422, `{ "errors": [ { "message": "substitution language syntax error in template content", "description": "Error while compiling header Subject: substitution statement missing ending '}}'", "code": "3000", "part": "Header:Subject" } ] }`, ""},

		{&sp.Template{Content: sp.Content{Subject: "s", HTML: "h", From: "f"}},
			errors.New(`parsing api response: invalid character 'B' looking for beginning of value`),
			503, `Bad Gateway`, ""},

		{&sp.Template{Content: sp.Content{Subject: "s", HTML: "h", From: "f"}}, nil,
			200, `{"results":{"id":"new-template"}}`, "new-template"},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "POST", test.status, sp.TemplatesPathFormat, test.json)

		id, _, err := testClient.TemplateCreate(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("TemplateCreate[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("TemplateCreate[%d] => err %q want %q", idx, err, test.err)
		} else if id != test.id {
			t.Errorf("TemplateCreate[%d] => id %q want %q", idx, id, test.id)
		}
	}
}

func TestTemplateGet(t *testing.T) {
	for idx, test := range []struct {
		in     *sp.Template
		draft  bool
		err    error
		status int
		json   string
		out    *sp.Template
	}{
		{nil, false, errors.New("TemplateGet called with nil Template"), 200, "", nil},
		{&sp.Template{ID: ""}, false, errors.New("TemplateGet called with blank id"), 200, "", nil},
		{&sp.Template{ID: "nope"}, false, errors.New(`[{"message":"Resource could not be found","code":"","description":""}]`), 404, `{ "errors": [ { "message": "Resource could not be found" } ] }`, nil},
		{&sp.Template{ID: "nope"}, false, errors.New(`parsing api response: unexpected end of JSON input`), 400, `{`, nil},
		{&sp.Template{ID: "id"}, false, errors.New("Unexpected response to TemplateGet"), 200, `{"foo":{}}`, nil},
		{&sp.Template{ID: "id"}, false, errors.New("parsing api response: unexpected end of JSON input"), 200, `{"foo":{}`, nil},

		{&sp.Template{ID: "id"}, false, errors.New(`parsing api response: invalid character 'B' looking for beginning of value`), 503, `Bad Gateway`, nil},

		{&sp.Template{ID: "id"}, false, nil, 200, `{"results":{"content":{"from":{"email":"a@b.com","name": "a b"},"html":"<blink>hi!</blink>","subject":"blink","text":"no blink ;_;"},"id":"id"}}`, &sp.Template{ID: "id", Content: sp.Content{From: map[string]interface{}{"email": "a@b.com", "name": "a b"}, HTML: "<blink>hi!</blink>", Text: "no blink ;_;", Subject: "blink"}}},
	} {
		testSetup(t)
		defer testTeardown()

		id := ""
		if test.in != nil {
			id = test.in.ID
		}
		mockRestResponseBuilderFormat(t, "GET", test.status, sp.TemplatesPathFormat+"/"+id, test.json)

		_, err := testClient.TemplateGet(test.in, test.draft)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("TemplateGet[%d] => err %v want %v", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("TemplateGet[%d] => err %v want %v", idx, err, test.err)
		} else if test.out != nil {
			if !reflect.DeepEqual(test.out, test.in) {
				t.Errorf("TemplateGet[%d] => template got/want:\n%q\n%q", idx, test.in, test.out)
			}
		}
	}
}

func TestTemplateUpdate(t *testing.T) {
	for idx, test := range []struct {
		in     *sp.Template
		pub    bool
		err    error
		status int
		json   string
	}{
		{nil, false, errors.New("Update called with nil Template"), 0, ""},
		{&sp.Template{ID: ""}, false, errors.New("Update called with blank id"), 0, ""},
		{&sp.Template{ID: "id", Content: sp.Content{}}, false, errors.New("Template requires a non-empty Content.Subject"), 0, ""},
		{&sp.Template{ID: "id", Content: sp.Content{Subject: "s", HTML: "h", From: "f"}}, false, errors.New("parsing api response: unexpected end of JSON input"), 0, `{ "errors": [ { "message": "truncated json" }`},

		{&sp.Template{ID: "id", Content: sp.Content{Subject: "s{{", HTML: "h", From: "f"}}, false,
			sp.SPErrors([]sp.SPError{{
				Message:     "substitution language syntax error in template content",
				Description: "Error while compiling header Subject: substitution statement missing ending '}}'",
				Code:        "3000",
				Part:        "Header:Subject",
			}}),
			422, `{ "errors": [ { "message": "substitution language syntax error in template content", "description": "Error while compiling header Subject: substitution statement missing ending '}}'", "code": "3000", "part": "Header:Subject" } ] }`},

		{&sp.Template{ID: "id", Content: sp.Content{Subject: "s", HTML: "h", From: "f"}}, false, nil, 200, ""},
	} {
		testSetup(t)
		defer testTeardown()

		id := ""
		if test.in != nil {
			id = test.in.ID
		}
		mockRestResponseBuilderFormat(t, "PUT", test.status, sp.TemplatesPathFormat+"/"+id, test.json)

		_, err := testClient.TemplateUpdate(test.in, test.pub)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("TemplateUpdate[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("TemplateUpdate[%d] => err %q want %q", idx, err, test.err)
		}
	}
}

func TestTemplates(t *testing.T) {
	for idx, test := range []struct {
		err    error
		status int
		json   string
	}{
		{errors.New("parsing api response: unexpected end of JSON input"), 0, `{ "errors": [ { "message": "truncated json" }`},
		{errors.New("[{\"message\":\"truncated json\",\"code\":\"\",\"description\":\"\"}]"), 0, `{ "errors": [ { "message": "truncated json" } ] }`},
		{nil, 200, `{ "results": [ { "description": "A test message from SparkPost.com", "id": "my-first-email", "last_update_time": "2006-01-02T15:04:05+00:00", "name": "My First Email", "published": false } ] }`},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "GET", test.status, sp.TemplatesPathFormat, test.json)

		_, _, err := testClient.Templates()
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("Templates[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("Templates[%d] => err %q want %q", idx, err, test.err)
		}
	}
}

func TestTemplateDelete(t *testing.T) {
	for idx, test := range []struct {
		id     string
		err    error
		status int
		json   string
	}{
		{"", errors.New("Delete called with blank id"), 0, ""},
		{"nope", errors.New(`[{"message":"Resource could not be found","code":"","description":""}]`), 404, `{ "errors": [ { "message": "Resource could not be found" } ] }`},
		{"id", nil, 200, "{}"},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "DELETE", test.status, sp.TemplatesPathFormat+"/"+test.id, test.json)

		_, err := testClient.TemplateDelete(test.id)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("TemplateDelete[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("TemplateDelete[%d] => err %q want %q", idx, err, test.err)
		}
	}
}

func TestTemplatePreview(t *testing.T) {
	for idx, test := range []struct {
		id     string
		opts   *sp.PreviewOptions
		err    error
		status int
		json   string
	}{
		{"", nil, errors.New("Preview called with blank id"), 200, ""},
		{"id", &sp.PreviewOptions{map[string]interface{}{
			"func": func() { return }},
		}, errors.New("json: unsupported type: func()"), 200, ""},
		{"id", nil, errors.New("parsing api response: unexpected end of JSON input"), 200, "{"},
		{"nope", nil, errors.New(`[{"message":"Resource could not be found","code":"","description":""}]`), 404, `{ "errors": [ { "message": "Resource could not be found" } ] }`},

		{"id", nil, nil, 200, ""},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "POST", test.status, sp.TemplatesPathFormat+"/"+test.id+"/preview", test.json)

		_, err := testClient.TemplatePreview(test.id, test.opts)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("TemplatePreview[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("TemplatePreview[%d] => err %q want %q", idx, err, test.err)
		}
	}
}

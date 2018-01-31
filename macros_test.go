package gosparkpost_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

var upperMacro = sp.Macro{Name: "ext_upper", Func: strings.ToUpper}
var lowerMacro = sp.Macro{Name: "ext_lower", Func: strings.ToLower}
var noopMacro = sp.Macro{Name: "ext_noop", Func: NoopMacro}

func NoopMacro(in string) string { return in }

func TestRegisterMacro(t *testing.T) {
	tests := []struct {
		macro *sp.Macro
		err   error
	}{
		{nil, errors.New(`can't add nil Macro`)},
		{&sp.Macro{Name: ""}, errors.New(`Macro names must only contain \w characters`)},
		{&sp.Macro{Name: "b:ar"}, errors.New(`Macro names must only contain \w characters`)},
		{&sp.Macro{Name: "bar"}, errors.New(`Macro must have non-nil Func field`)},
		{&sp.Macro{Name: "bar", Func: strings.ToUpper}, nil},
	}

	for idx, test := range tests {
		testSetup(t)
		defer testTeardown()

		err := testClient.RegisterMacro(test.macro)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("RegisterMacro[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("TemplateCreate[%d] => err %q want %q", idx, err, test.err)
		}
	}
	testClient = nil
}

func TestRecipientApply(t *testing.T) {
	tests := []struct {
		r   *sp.Recipient
		in  string
		out string
		err error
	}{
		// nil recipient returns input string unchanged
		{nil, "in", "in", nil},

		// tokenize catches template error
		// this only happens in this function when used standalone
		// when using ApplyMacros, all tokenize errors are caught at that level
		{&sp.Recipient{}, "{{{ foo }}", "", errors.New(`mismatched curly braces near "{{{ foo }}"`)},

		{&sp.Recipient{
			Address:  "test@example.com",
			Metadata: map[string]interface{}{"foo": "bar"},
		}, "{{ foo }}", "bar", nil},
	}

	for idx, test := range tests {
		out, err := test.r.Apply(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("Apply[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("Apply[%d] => err %q want %q", idx, err, test.err)
		} else if out != test.out {
			t.Errorf("Apply[%d] => got/want:\n%s\n%s\n", idx, out, test.out)
		}
	}
}

func TestApplyMacros(t *testing.T) {
	tests := []struct {
		macros   []sp.Macro
		recip    *sp.Recipient
		template string
		out      string
		err      error
	}{
		// not enough closing curlies
		{[]sp.Macro{upperMacro}, nil, "{{{ext_upper}}", "", errors.New(`mismatched curly braces near "{{{ext_upper}}"`)},
		// balanced triple curlies
		{[]sp.Macro{upperMacro}, nil, "{{{ext_upper}}}", "", nil},
		// extra trailing curlies
		{[]sp.Macro{upperMacro}, nil, "{{ext_upper}}}", "}", nil},

		// no macros defined, pass through
		{[]sp.Macro{}, nil, "{{ext_upper}}", "{{ext_upper}}", nil},
		{[]sp.Macro{}, nil, "{{ext_upper foo}}", "{{ext_upper foo}}", nil},
		// macros defined, pass through ones that don't match
		{[]sp.Macro{lowerMacro}, nil, "{{ext_upper}}", "{{ext_upper}}", nil},

		// multiple macros, preserving space outside blocks
		{[]sp.Macro{upperMacro, lowerMacro},
			nil, " {{ ext_upper bar }} {{ ext_lower FOO }} ", " BAR foo ", nil},

		// invalid recipient address
		{[]sp.Macro{upperMacro},
			&sp.Recipient{
				Address: 42,
			}, "{{ ext_upper {{abc}} }}", "", errors.New("parsing recipient address: unsupported Recipient.Address value type [int]")},

		// bad recipient metadata format
		{[]sp.Macro{upperMacro},
			&sp.Recipient{
				Address:  "test@example.com",
				Metadata: 42,
			}, "{{ ext_upper }}", "", errors.New("unexpected metadata type [int] for recipient test@example.com")},

		// bad recipient substitution data format
		{[]sp.Macro{upperMacro},
			&sp.Recipient{
				Address:          "test@example.com",
				SubstitutionData: 42,
			}, "{{ ext_upper }}", "", errors.New("unexpected substitution data type [int] for recipient test@example.com")},

		// non-string meta/sub data
		{[]sp.Macro{upperMacro, lowerMacro},
			&sp.Recipient{
				Address:          "test@example.com",
				Metadata:         map[string]interface{}{"abc": 42},
				SubstitutionData: map[string]interface{}{"ghi": 42},
			}, "{{ ext_upper {{abc}} }}{{ ext_lower ({{ghi}}) }}", "{{ABC}}({{ghi}})", nil},

		// substitution data overrides metadata
		{[]sp.Macro{noopMacro},
			&sp.Recipient{
				Address:          "test@example.com",
				Metadata:         map[string]interface{}{"abc": "def", "ghi": "jkl"},
				SubstitutionData: map[string]interface{}{"abc": "DEF"},
			}, "{{ ext_noop {{abc}} }}{{ ext_noop ({{ghi}}) }}", "DEF(jkl)", nil},

		// recipient macros nested inside client macros
		{[]sp.Macro{upperMacro, lowerMacro},
			&sp.Recipient{
				Address:  "test@example.com",
				Metadata: map[string]interface{}{"abc": "def", "ghi": "JKL"},
			}, "{{ ext_upper {{abc}} }}{{ ext_lower ({{ghi}}) }}", "DEF(jkl)", nil},

		// not enough trailing curlies for nested macros
		{[]sp.Macro{upperMacro},
			&sp.Recipient{
				Address:  "test@example.com",
				Metadata: map[string]interface{}{"abc": "def", "ghi": "JKL"},
			}, "{{ ext_upper {{{abc}} }}", "", errors.New(`mismatched curly braces near "{{ ext_upper {{{abc}} }}"`)},
	}

	for idx, test := range tests {
		testSetup(t)
		defer testTeardown()

		for _, m := range test.macros {
			err := testClient.RegisterMacro(&m)
			if err != nil {
				t.Fatal(err)
			}
		}

		out, err := testClient.ApplyMacros(test.template, test.recip)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("ApplyMacros[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("ApplyMacros[%d] => err %q want %q", idx, err, test.err)
		} else if out != test.out {
			t.Errorf("ApplyMacros[%d] => got/want:\n%s\n%s\n", idx, out, test.out)
		}
		// discard client (where macros live) so next test gets a fresh start
		testClient = nil
	}
}

func TestInvoice(t *testing.T) {
	s := httptest.NewTLSServer(http.HandlerFunc(testInvoice))
	tx := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tx}
	spClient := &sp.Client{}
	invoiceMacro := &sp.Macro{Name: "sp_invoice", Func: InvoiceMacro(client)}
	spClient.RegisterMacro(invoiceMacro)
	r := &sp.Recipient{Address: "test@example.com", Metadata: map[string]interface{}{"invoice_id": "f00f00"}}

	prefix := `Here's your invoice:`
	tmpl := fmt.Sprintf(`%s {{ sp_invoice %s/{{invoice_id}} }}`, prefix, s.URL)
	out, err := spClient.ApplyMacros(tmpl, r)
	if err != nil {
		t.Errorf("%v", err)
	}
	if out != prefix+" "+string(carlinInvoice) {
		t.Errorf("InvoiceMacro - unexpected output:\n%s", out)
	}
}

func InvoiceMacro(client *http.Client) func(string) string {
	return func(url string) string {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err.Error()
		}
		res, err := client.Do(req)
		if err != nil {
			return err.Error()
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil && err != io.EOF {
			return err.Error()
		}
		return string(body)
	}
}

var carlinInvoice = []byte(`
  George Carlin                                            INVOICE
  carlin@example.org

  To:                                                   Invoice #6
      Stephen Hawking                         Date:   May 13, 2014
      hawking@example.org

  +-----------------------------------------------------------------+
  | Quantity |         Description         | Unit Price |   Total   |
  +-----------------------------------------------------------------+
  | 8        | Awesome Comedy Hour         | $99 CAD    | $792 CAD  |
  | 5        | Encore                      | $200 CAD   | $1000 CAD |
  +-----------------------------------------------------------------+

                                                    TOTAL: $1792 CAD

  Payment Instructions:

    Send to PayPal address carlin@example.org.
    Payment is due within 30 days.

  Thank you for your business! (ref# {{invoice_id}})`)

func testInvoice(w http.ResponseWriter, r *http.Request) {
	w.Write(carlinInvoice)
}

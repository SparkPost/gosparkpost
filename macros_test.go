package gosparkpost_test

import (
	"strings"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

var upperMacro = sp.Macro{Name: "ext_upper", Func: strings.ToUpper}
var lowerMacro = sp.Macro{Name: "ext_lower", Func: strings.ToLower}

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

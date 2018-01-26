package gosparkpost_test

import (
	"strings"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

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

func TestApplyMacros(t *testing.T) {
	tests := []struct {
		macros   []sp.Macro
		recip    *sp.Recipient
		template string
		out      string
		err      error
	}{
		{[]sp.Macro{
			sp.Macro{Name: "ext_foo", Func: strings.ToUpper},
			sp.Macro{Name: "ext_bar", Func: strings.ToLower},
		}, nil, "{{ ext_foo bar }}{{ ext_bar FOO }}", "BARfoo", nil},

		{
			[]sp.Macro{
				sp.Macro{Name: "ext_foo", Func: strings.ToUpper},
				sp.Macro{Name: "ext_bar", Func: strings.ToLower},
			}, &sp.Recipient{Address: "test@example.com", Metadata: map[string]interface{}{"abc": "def", "ghi": "JKL"}},
			"{{ ext_foo {{abc}} }}{{ ext_bar {{ghi}} }}", "DEFjkl", nil},
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
	}
}

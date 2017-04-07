package gosparkpost_test

import (
	"reflect"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

func TestAddressValidation(t *testing.T) {
	for idx, test := range []struct {
		in  interface{}
		err error
		out sp.Address
	}{
		{nil, errors.New("unsupported Recipient.Address value type [%!s(<nil>)]"), sp.Address{}},
		{"", errors.New("Recipient.Address may not be empty"), sp.Address{}},
		{"a@b.com", nil, sp.Address{Email: "a@b.com"}},
		{sp.Address{"a@b.com", "A B", "c@d.com"}, nil, sp.Address{"a@b.com", "A B", "c@d.com"}},
		{map[string]interface{}{"foo": 42}, errors.New("strings are required for all Recipient.Address values"), sp.Address{}},
		{map[string]interface{}{"Name": "A B", "email": "a@b.com", "header_To": "c@d.com"}, nil, sp.Address{"a@b.com", "A B", "c@d.com"}},
		{map[string]string{"Name": "A B", "email": "a@b.com", "header_To": "c@d.com"}, nil, sp.Address{"a@b.com", "A B", "c@d.com"}},
	} {
		a, err := sp.ParseAddress(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("ParseAddress[%d] => err %q, want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("ParseAddress[%d] => err %q, want %q", idx, err, test.err)
		} else if !reflect.DeepEqual(a, test.out) {
			t.Errorf("ParseAddress[%d] => got/want:\n%q\n%q", idx, a, test.out)
		}
	}
}

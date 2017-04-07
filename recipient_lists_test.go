package gosparkpost_test

import (
	"reflect"
	"strings"
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

func TestRecipientValidation(t *testing.T) {
	for idx, test := range []struct {
		in  sp.Recipient
		err error
	}{
		{sp.Recipient{}, errors.New("unsupported Recipient.Address value type [%!s(<nil>)]")},
		{sp.Recipient{Address: "a@b.com"}, nil},
	} {
		err := test.in.Validate()
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("Recipient.Validate[%d] => err %q, want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("Recipient.Validate[%d] => err %q, want %q", idx, err, test.err)
		}
	}
}

func TestRecipientListValidation(t *testing.T) {
	for idx, test := range []struct {
		in  *sp.RecipientList
		err error
	}{
		{nil, errors.New("Can't validate a nil RecipientList")},
		{&sp.RecipientList{}, errors.New("RecipientList requires at least one Recipient")},

		{&sp.RecipientList{ID: strings.Repeat("id", 33),
			Recipients: &[]sp.Recipient{{}}}, errors.New("RecipientList id may not be longer than 64 bytes")},
		{&sp.RecipientList{ID: "id", Name: strings.Repeat("name", 17),
			Recipients: &[]sp.Recipient{{}}}, errors.New("RecipientList name may not be longer than 64 bytes")},
		{&sp.RecipientList{ID: "id", Name: "name", Description: strings.Repeat("desc", 257),
			Recipients: &[]sp.Recipient{{}}}, errors.New("RecipientList description may not be longer than 1024 bytes")},

		{&sp.RecipientList{ID: "id", Name: "name", Description: "desc",
			Recipients: &[]sp.Recipient{{}}}, errors.New("unsupported Recipient.Address value type [%!s(<nil>)]")},
		{&sp.RecipientList{ID: "id", Name: "name", Description: "desc",
			Recipients: &[]sp.Recipient{{Address: "a@b.com"}}}, nil},
	} {
		err := test.in.Validate()
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("RecipientList.Validate[%d] => err %q, want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("RecipientList.Validate[%d] => err %q, want %q", idx, err, test.err)
		}
	}
}

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
		{nil, errors.New("unsupported Recipient.Address value type [<nil>]"), sp.Address{}},
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
		{sp.Recipient{}, errors.New("unsupported Recipient.Address value type [<nil>]")},
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
			Recipients: []sp.Recipient{{}}}, errors.New("RecipientList id may not be longer than 64 bytes")},
		{&sp.RecipientList{ID: "id", Name: strings.Repeat("name", 17),
			Recipients: []sp.Recipient{{}}}, errors.New("RecipientList name may not be longer than 64 bytes")},
		{&sp.RecipientList{ID: "id", Name: "name", Description: strings.Repeat("desc", 257),
			Recipients: []sp.Recipient{{}}}, errors.New("RecipientList description may not be longer than 1024 bytes")},

		{&sp.RecipientList{ID: "id", Name: "name", Description: "desc",
			Recipients: []sp.Recipient{{}}}, errors.New("unsupported Recipient.Address value type [<nil>]")},
		{&sp.RecipientList{ID: "id", Name: "name", Description: "desc",
			Recipients: []sp.Recipient{{Address: "a@b.com"}}}, nil},
	} {
		err := test.in.Validate()
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("RecipientList.Validate[%d] => err %q, want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("RecipientList.Validate[%d] => err %q, want %q", idx, err, test.err)
		}
	}
}

func TestRecipientListCreate(t *testing.T) {
	for idx, test := range []struct {
		in     *sp.RecipientList
		err    error
		status int
		json   string
		id     string
	}{
		{nil, errors.New("Create called with nil RecipientList"), 0, "", ""},
		{&sp.RecipientList{}, errors.New("RecipientList requires at least one Recipient"), 0, "", ""},
		{&sp.RecipientList{ID: "id", Recipients: []sp.Recipient{{Address: "a@b.com"}}},
			errors.New("Unexpected response to Recipient List creation (results)"), 200, `{"foo":{"id":"id"}}`, ""},
		{&sp.RecipientList{ID: "id", Recipients: []sp.Recipient{{Address: "a@b.com"}}},
			errors.New("Unexpected response to Recipient List creation (id)"), 200, `{"results":{"ID":"id"}}`, ""},
		{&sp.RecipientList{ID: "id", Attributes: func() { return }, Recipients: []sp.Recipient{{Address: "a@b.com"}}},
			errors.New("json: unsupported type: func()"), 200, `{"results":{"ID":"id"}}`, ""},
		{&sp.RecipientList{ID: "id", Recipients: []sp.Recipient{{Address: "a@b.com"}}},
			errors.New("parsing api response: unexpected end of JSON input"), 200, `{"results":{"ID":"id"}`, ""},

		{&sp.RecipientList{ID: "id", Recipients: []sp.Recipient{{Address: "a@b.com"}}},
			errors.New(`[{"message":"List already exists","code":"5001","description":"List 'id' already exists"}]`), 400,
			`{"errors":[{"message":"List already exists","code":"5001","description":"List 'id' already exists"}]}`, ""},

		{&sp.RecipientList{ID: "id", Recipients: []sp.Recipient{{Address: "a@b.com"}}}, nil, 200,
			`{"results":{"total_rejected_recipients": 0,"total_accepted_recipients":1,"id":"id"}}`, "id"},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "POST", test.status, sp.RecipientListsPathFormat, test.json)

		id, _, err := testClient.RecipientListCreate(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("RecipientListCreate[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("RecipientListCreate[%d] => err %q want %q", idx, err, test.err)
		} else if id != test.id {
			t.Errorf("RecipientListCreate[%d] => id %q want %q", idx, id, test.id)
		}
	}
}

func TestRecipientLists(t *testing.T) {
	var res200 = loadTestFile(t, "test/json/recipient_lists_200.json")
	var acc200 = []int{3, 8}
	var rl200 = []sp.RecipientList{{
		ID:          "unique_id_4_graduate_students_list",
		Name:        "graduate_students",
		Description: "An email list of graduate students at UMBC",
		Attributes: map[string]interface{}{
			"internal_id":   float64(112),
			"list_group_id": float64(12321),
		},
		Accepted: &acc200[0],
	}, {
		ID:          "unique_id_4_undergraduates",
		Name:        "undergraduate_students",
		Description: "An email list of undergraduate students at UMBC",
		Attributes: map[string]interface{}{
			"internal_id":   float64(111),
			"list_group_id": float64(11321),
		},
		Accepted: &acc200[1],
	}}

	for idx, test := range []struct {
		err    error
		status int
		json   string
		out    []sp.RecipientList
	}{
		{nil, 200, `{"results":[{}]}`, []sp.RecipientList{{}}},
		{errors.New("Unexpected response to RecipientList list"), 200, `{"foo":[{}]}`, nil},
		{errors.New("unexpected end of JSON input"), 200, `{"results":[{}]`, nil},
		{errors.New(`[{"message":"No RecipientList for you!","code":"","description":""}]`), 401, `{"errors":[{"message":"No RecipientList for you!"}]}`, nil},
		{errors.New("parsing api response: unexpected end of JSON input"), 401, `{"errors":[]`, nil},
		{nil, 200, res200, rl200},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "GET", test.status, sp.RecipientListsPathFormat, test.json)

		lists, _, err := testClient.RecipientLists()
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("RecipientLists[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("RecipientLists[%d] => err %q want %q", idx, err, test.err)
		} else if test.out != nil && !reflect.DeepEqual(lists, test.out) {
			t.Errorf("RecipientLists[%d] => got/want:\n%q\n%q", idx, lists, test.out)
		}
	}
}

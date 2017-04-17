package gosparkpost_test

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

func TestSubaccountCreate(t *testing.T) {
	var res200 = loadTestFile(t, "test/json/subaccount_create_200.json")

	for idx, test := range []struct {
		in     *sp.Subaccount
		err    error
		status int
		json   string
		out    *sp.Subaccount
	}{
		{nil, errors.New("Create called with nil Subaccount"), 0, "", nil},

		{&sp.Subaccount{Name: "n", KeyLabel: "kl"},
			errors.New("Unexpected response to Subaccount creation (results)"), 200,
			strings.Replace(res200, `"results"`, `"foo"`, 1), nil},
		{&sp.Subaccount{Name: "n", KeyLabel: "kl"},
			errors.New("parsing api response: unexpected end of JSON input"), 200,
			res200[:len(res200)/2], nil},
		{&sp.Subaccount{Name: "n", KeyLabel: "kl"},
			errors.New("Unexpected response to Subaccount creation (subaccount_id)"), 200,
			strings.Replace(res200, `"subaccount_id"`, `"foo"`, 1), nil},
		{&sp.Subaccount{Name: "n", KeyLabel: "kl"},
			errors.New("Unexpected response to Subaccount creation (short_key)"), 200,
			strings.Replace(res200, `"short_key"`, `"foo"`, 1), nil},

		{&sp.Subaccount{Name: "n", KeyLabel: "kl"},
			errors.New(`[{"message":"error","code":"","description":""}]`), 400,
			`{"errors":[{"message":"error"}]}`, nil},

		{&sp.Subaccount{Name: "n", KeyLabel: "kl"}, nil, 200, res200,
			&sp.Subaccount{
				ID:       888,
				Name:     "n",
				KeyLabel: "kl",
				ShortKey: "cf80",
				Grants:   sp.SubaccountGrants,
			}},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "POST", test.status, sp.SubaccountsPathFormat, test.json)

		_, err := testClient.SubaccountCreate(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("SubaccountCreate[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("SubaccountCreate[%d] => err %q want %q", idx, err, test.err)
		} else if test.out != nil && !reflect.DeepEqual(test.in, test.out) {
			t.Errorf("SubaccountCreate[%d] => got/want:\n%q\n%q", idx, test.in, test.out)
		}
	}
}

func TestSubaccountUpdate(t *testing.T) {
	for idx, test := range []struct {
		in     *sp.Subaccount
		err    error
		status int
		json   string
	}{
		{nil, errors.New("Subaccount Update called with nil Subaccount"), 0, ""},
		{&sp.Subaccount{ID: 42, Name: "n", Status: "super"},
			errors.New("Not a valid subaccount status"), 0, ""},

		{&sp.Subaccount{ID: 42, Name: "n"},
			errors.New("parsing api response: unexpected end of JSON input"), 200,
			`{"foo":{"message":"syntax error"}`},
		{&sp.Subaccount{ID: 42, Name: "n"},
			errors.New(`[{"message":"error","code":"","description":""}]`), 400,
			`{"errors":[{"message":"error"}]}`},

		{&sp.Subaccount{ID: 42, Name: "n", Status: "active"}, nil, 200,
			`{"results":{"message":"Successfully updated subaccount information"}}`},
	} {
		testSetup(t)
		defer testTeardown()

		id := "0"
		if test.in != nil {
			id = strconv.Itoa(test.in.ID)
		}
		mockRestResponseBuilderFormat(t, "PUT", test.status, sp.SubaccountsPathFormat+"/"+id, test.json)

		_, err := testClient.SubaccountUpdate(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("SubaccountUpdate[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("SubaccountUpdate[%d] => err %q want %q", idx, err, test.err)
		}
	}
}

func TestSubaccounts(t *testing.T) {
	var res200 = loadTestFile(t, "test/json/subaccounts_200.json")
	var sa200 = []sp.Subaccount{{
		ID:               123,
		Name:             "Joe's Garage",
		Status:           "active",
		IPPool:           "my_ip_pool",
		ComplianceStatus: "active",
	}, {
		ID:               456,
		Name:             "SharkPost",
		Status:           "active",
		ComplianceStatus: "active",
	}, {
		ID:               789,
		Name:             "Dev Avocado",
		Status:           "suspended",
		ComplianceStatus: "active",
	}}

	for idx, test := range []struct {
		err    error
		status int
		json   string
		out    []sp.Subaccount
	}{
		{errors.New("unexpected end of JSON input"), 200, `{"foo":[]`, nil},
		{errors.New("Unexpected response to Subaccount list"), 200, `{"foo":[]}`, nil},

		{errors.New(`[{"message":"error","code":"","description":""}]`), 400,
			`{"errors":[{"message":"error"}]}`, nil},

		{nil, 200, res200, sa200},
	} {
		testSetup(t)
		defer testTeardown()
		mockRestResponseBuilderFormat(t, "GET", test.status, sp.SubaccountsPathFormat, test.json)

		subs, _, err := testClient.Subaccounts()
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("Subaccounts[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("Subaccounts[%d] => err %q want %q", idx, err, test.err)
		} else if test.out != nil && !reflect.DeepEqual(subs, test.out) {
			t.Errorf("Subaccounts[%d] => got/want:\n%q\n%q", idx, subs, test.out)
		}
	}
}

func TestSubaccount(t *testing.T) {
	var res200 = loadTestFile(t, "test/json/subaccount_200.json")
	var sub200 = &sp.Subaccount{
		ID:               123,
		Name:             "Joe's Garage",
		Status:           "active",
		IPPool:           "assigned_ip_pool",
		ComplianceStatus: "active",
	}

	for idx, test := range []struct {
		in     int
		err    error
		status int
		json   string
		out    *sp.Subaccount
	}{
		{42, errors.New("unexpected end of JSON input"), 200, "{", nil},
		{42, errors.New("Unexpected response to Subaccount"), 200, `{"foo":{}}`, nil},
		{42, errors.New(`[{"message":"error","code":"","description":""}]`), 400,
			`{"errors":[{"message":"error"}]}`, nil},

		{123, nil, 200, res200, sub200},
	} {
		testSetup(t)
		defer testTeardown()

		id := strconv.Itoa(test.in)
		mockRestResponseBuilderFormat(t, "GET", test.status, sp.SubaccountsPathFormat+"/"+id, test.json)

		sub, _, err := testClient.Subaccount(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("SubaccountCreate[%d] => err %q want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("SubaccountCreate[%d] => err %q want %q", idx, err, test.err)
		} else if test.out != nil && !reflect.DeepEqual(sub, test.out) {
			t.Errorf("SubaccountCreate[%d] => got/want:\n%q\n%q", idx, sub, test.out)
		}
	}
}

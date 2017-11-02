package gosparkpost_test

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/pkg/errors"
)

var (
	testMux    *http.ServeMux
	testClient *sp.Client
	testServer *httptest.Server
)

func testSetup(t *testing.T) {
	// spin up a test server
	testMux = http.NewServeMux()
	testServer = httptest.NewTLSServer(testMux)
	// our client configured to hit the https test server with self-signed cert
	tx := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	if testClient == nil {
		testClient = &sp.Client{Client: &http.Client{Transport: tx}}
	}
	testClient.Config = &sp.Config{Verbose: true}
	testUrl, err := url.Parse(testServer.URL)
	if err != nil {
		t.Fatalf("Test server url parsing failed: %v", err)
	}
	testClient.Config.BaseUrl = testUrl.String()
	err = testClient.Init(testClient.Config)
	if err != nil {
		t.Fatalf("Test client init failed: %v", err)
	}
}

func testTeardown() {
	testServer.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Fatalf("Request method: %v, want %v", got, want)
	}
}

func testFailVerbose(t *testing.T, res *sp.Response, fmt string, args ...interface{}) {
	if res != nil {
		for _, e := range res.Verbose {
			t.Error(e)
		}
	}
	t.Fatalf(fmt, args...)
}

func TestNewConfig(t *testing.T) {
	for idx, test := range []struct {
		in  map[string]string
		cfg *sp.Config
		err error
	}{
		{map[string]string{}, nil, errors.New("BaseUrl is required for api config")},
		{map[string]string{"baseurl": "http://example.com"}, nil, errors.New("ApiKey is required for api config")},
		{map[string]string{"baseurl": "http://example.com", "apikey": "foo"}, &sp.Config{BaseUrl: "http://example.com", ApiKey: "foo"}, nil},
	} {
		cfg, err := sp.NewConfig(test.in)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("NewConfig[%d] => err %q, want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("NewConfig[%d] => err %q, want %q", idx, err, test.err)
		} else if cfg == nil && test.cfg != nil || cfg != nil && test.cfg == nil {
			t.Errorf("NewConfig[%d] => cfg %v, want %v", idx, cfg, test.cfg)
		}
	}
}

func TestSPErrors(t *testing.T) {
	for idx, test := range []struct {
		name     string
		in       string
		code     string
		parsed   sp.SPErrors
		expected string
	}{
		{"string code",
			`[{"message":"This is fine.","code":"13","description":""}]`, "13",
			sp.SPErrors([]sp.SPError{{Message: "This is fine.", Code: "13"}}),
			`[{"message":"This is fine.","code":"13","description":""}]`,
		}, {"int code",
			`[{"message":"This is fine.","code":42,"description":""}]`, "42",
			sp.SPErrors([]sp.SPError{{Message: "This is fine.", Code: "42"}}),
			`[{"message":"This is fine.","code":"42","description":""}]`,
		}, {"just message",
			`[{"message":"This is fine.","code":"","description":""}]`, "",
			sp.SPErrors([]sp.SPError{{Message: "This is fine."}}),
			`[{"message":"This is fine.","code":"","description":""}]`,
		},
	} {
		// test the round trip from []byte to sp.SPErrors and back
		errs := sp.SPErrors{}
		err := json.Unmarshal([]byte(test.in), &errs)
		if err != nil {
			t.Fatal(err)
		}
		if test.parsed[0].Code.String() != test.code {
			t.Errorf("SPErrors[%d] (%s) code => %q, want %q", idx, test.name, test.parsed[0].Code, test.code)
		}
		errstr := test.parsed.Error()
		if test.expected != errstr {
			t.Errorf("SPErrors.Stringify[%d] (%s) => %q, want %q", idx, test.name, errstr, test.in)
		}
	}
}

func TestDoRequest(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	for idx, test := range []struct {
		client *sp.Client
		method string
		err    error
	}{
		{nil, "", errors.New("Client must be non-nil!")},
		{&sp.Client{}, "", errors.New("Client.Client (http.Client) must be non-nil!")},
		{&sp.Client{Client: http.DefaultClient}, "", errors.New("Client.Config must be non-nil!")},
		{testClient, "💩", errors.New(`building request: net/http: invalid method "💩"`)},
	} {
		_, err := test.client.DoRequest(nil, test.method, "", nil)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("DoRequest[%d] => err %v, want %v", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("DoRequest[%d] => err %v, want %v", idx, err, test.err)
		}
	}
}

func TestHTTPError(t *testing.T) {
	for idx, test := range []struct {
		name string
		res  *sp.Response
		err  error
	}{
		{"nil response", nil, errors.New("Internal error: Response may not be nil")},
		{"nil http", &sp.Response{}, errors.New("Internal error: Response.HTTP may not be nil")},
		{"got error", &sp.Response{HTTP: &http.Response{Status: "418 I'm a teapot"}},
			sp.SPErrors([]sp.SPError{{Code: sp.ErrorCode("418 I'm a teapot"), Description: "HTTP/JSON Error"}})},
	} {
		err := test.res.HTTPError()
		if test.err != nil && err == nil {
			t.Errorf("HTTPError[%d] (%s) => no error, wanted %q", idx, test.name, test.err)
		} else if test.err == nil && err != nil {
			t.Errorf("HTTPError[%d] (%s) => unexpected error: %q", idx, test.name, err)
		} else if test.err.Error() != err.Error() {
			t.Errorf("HTTPError[%d] (%s) => mismatched errors (got/want):\n%q\n%q", idx, test.name, err, test.err)
		}
	}
}

func TestInit(t *testing.T) {
	for idx, test := range []struct {
		api *sp.Client
		cfg *sp.Config
		out *sp.Config
		err error
	}{
		{&sp.Client{}, &sp.Config{BaseUrl: ""}, &sp.Config{BaseUrl: "https://api.sparkpost.com"}, nil},
		{&sp.Client{}, &sp.Config{BaseUrl: "http://api.sparkpost.com"}, nil, errors.New("API base url must be https!")},
	} {
		err := test.api.Init(test.cfg)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("Init[%d] => err %q, want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("Init[%d] => err %q, want %q", idx, err, test.err)
		} else if test.out != nil && test.api.Config.BaseUrl != test.out.BaseUrl {
			t.Errorf("Init[%d] => BaseUrl %q, want %q", idx, test.api.Config.BaseUrl, test.out.BaseUrl)
		}
	}
}

func loadTestFile(t *testing.T, fileToLoad string) string {
	b, err := ioutil.ReadFile(fileToLoad)

	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	return string(b)
}

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

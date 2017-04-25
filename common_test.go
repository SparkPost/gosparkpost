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
	testClient = &sp.Client{Client: &http.Client{Transport: tx}}
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

func TestJson(t *testing.T) {
	var e = sp.SPErrors([]sp.SPError{{Message: "This is fine."}})
	var exp = `[{"message":"This is fine.","code":"","description":""}]`
	str := e.Error()
	if str != exp {
		t.Errorf("*SPError.Stringify => %q, want %q", str, exp)
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
	var res *sp.Response
	err := res.HTTPError()
	if err == nil {
		t.Error("nil response should fail")
	}

	res = &sp.Response{}
	err = res.HTTPError()
	if err == nil {
		t.Error("nil http should fail")
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

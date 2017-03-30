package gosparkpost_test

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
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

var newConfigTests = []struct {
	in  map[string]string
	cfg *sp.Config
	err error
}{
	{map[string]string{}, nil, errors.New("BaseUrl is required for api config")},
	{map[string]string{"baseurl": "http://example.com"}, nil, errors.New("ApiKey is required for api config")},
	{map[string]string{"baseurl": "http://example.com", "apikey": "foo"}, &sp.Config{BaseUrl: "http://example.com", ApiKey: "foo"}, nil},
}

func TestNewConfig(t *testing.T) {
	for idx, test := range newConfigTests {
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
	var e = &sp.Error{Message: "This is fine."}
	var exp = `{"message":"This is fine.","code":"","description":""}`
	str := e.Json()
	if str != exp {
		t.Errorf("*Error.Json => %q, want %q", str, exp)
	}
}

var initTests = []struct {
	api *sp.Client
	cfg *sp.Config
	out *sp.Config
	err error
}{
	{&sp.Client{}, &sp.Config{BaseUrl: ""}, &sp.Config{BaseUrl: "https://api.sparkpost.com"}, nil},
	{&sp.Client{}, &sp.Config{BaseUrl: "http://api.sparkpost.com"}, nil, errors.New("API base url must be https!")},
}

func TestInit(t *testing.T) {
	for idx, test := range initTests {
		err := test.api.Init(test.cfg)
		if err == nil && test.err != nil || err != nil && test.err == nil {
			t.Errorf("Init[%d] => err %q, want %q", idx, err, test.err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Errorf("NewConfig[%d] => err %q, want %q", idx, err, test.err)
		} else if test.out != nil && test.api.Config.BaseUrl != test.out.BaseUrl {
			t.Errorf("Init[%d] => BaseUrl %q, want %q", idx, test.api.Config.BaseUrl, test.out.BaseUrl)
		}
	}
}

/* // either make headers public or can it entirely in favor of context
func TestSetHeader(t *testing.T) {
	var val string
	var ok bool
	cl := &sp.Client{}
	cl.SetHeader("X-Foo", "Bar")
	if val, ok = cl.headers["X-Foo"]; !ok {
		t.Errorf("SetHeader => nil, want %q", "Bar")
	} else if val != "Bar" {
		t.Errorf("SetHeader => %q, want %q", val, "Bar")
	}
}
*/

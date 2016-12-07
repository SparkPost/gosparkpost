package gosparkpost

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	testMux    *http.ServeMux
	testClient *Client
	testServer *httptest.Server
)

func testSetup(t *testing.T) {
	// spin up a test server
	testMux = http.NewServeMux()
	testServer = httptest.NewTLSServer(testMux)
	// our client configured to hit the https test server with self-signed cert
	tx := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	testClient = &Client{Client: &http.Client{Transport: tx}}
	testClient.Config = &Config{Verbose: true}
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

func testFailVerbose(t *testing.T, res *Response, fmt string, args ...interface{}) {
	if res != nil {
		for _, e := range res.Verbose {
			t.Error(e)
		}
	}
	t.Fatalf(fmt, args...)
}

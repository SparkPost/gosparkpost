// Package api provides structures and functions used by other SparkPost API packages.
package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	certifi "github.com/certifi/gocertifi"
)

// Config includes all information necessary to make an API request.
type Config struct {
	BaseUrl string
	ApiKey  string
}

// Response contains information about the last HTTP response.
// Helpful when an error message doesn't necessarily give the complete picture.
type Response struct {
	HTTP    *http.Response
	Body    string
	Results map[string]string `json:"results,omitempty"`
	Errors  []Error           `json:"errors,omitempty"`
}

// API exists to be embedded in other API objects.
type API struct {
	Config   *Config
	Client   *http.Client
	Response *Response
}

// The error structure returned by SparkPost APIs.
type Error struct {
	Message     string `json:"message"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Part        string `json:"part,omitempty"`
	Line        int    `json:"line,omitempty"`
}

// Init pulls together everything necessary to make an API request.
func (api *API) Init(cfg *Config) (err error) {
	api.Config = cfg

	// load Mozilla cert pool
	pool, err := certifi.CACerts()
	if err != nil {
		return
	}

	// configure transport using Mozilla cert pool
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}

	// configure http client using transport
	api.Client = &http.Client{Transport: transport}

	return
}

// HttpPost sends a Post request with the provided JSON payload to the specified url.
// Query params are supported via net/url - roll your own and stringify it.
// Authenticate using the configured API key.
func (api *API) HttpPost(url string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", api.Config.ApiKey)
	return api.Client.Do(req)
}

// HttpGet sends a Get request to the specified url.
// Query params are supported via net/url - roll your own and stringify it.
// Authenticate using the configured API key.
func (api *API) HttpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", api.Config.ApiKey)
	return api.Client.Do(req)
}

// HttpDelete sends a Delete request to the provided url.
// Query params are supported via net/url - roll your own and stringify it.
// Authenticate using the configured API key.
func (api *API) HttpDelete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", api.Config.ApiKey)
	return api.Client.Do(req)
}

// ReadBody is a convenience wrapper for the response body.
func ReadBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// ParseApiResponse pulls info from JSON http responses into api.Response object.
// It's helpful to call api.AssertJson before calling this function.
func (api *API) ParseResponse(res *http.Response) error {
	body, err := ReadBody(res)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, api.Response)
	if err != nil {
		return err
	}
	api.Response.Body = string(body)
	api.Response.HTTP = res

	return nil
}

// AssertObject asserts that the provided variable is a map[string]something.
func AssertObject(obj interface{}, label string) error {
	// List of handled types from here:
	// http://golang.org/pkg/encoding/json/#Unmarshal
	switch objVal := obj.(type) {
	case map[string]interface{}:
		// auto-parsed nested json object
	case map[string]bool:
		// user-provided json literal (convenience)
	case map[string]float64:
		// user-provided json literal (convenience)
	case map[string]string:
		// user-provided json literal (convenience)
	case map[string][]interface{}:
		// user-provided json literal (convenience)
	case map[string]map[string]interface{}:
		// user-provided json literal (convenience)
	default:
		return fmt.Errorf("expected key/val pairs for %s, got [%s]", label, reflect.TypeOf(objVal))
	}
	return nil
}

// AssertJson returns an error if the provided HTTP response isn't JSON.
func AssertJson(res *http.Response) error {
	if res == nil {
		return fmt.Errorf("AssertJson got nil http.Response")
	}
	contentType := res.Header.Get("Content-Type")
	if !strings.EqualFold(contentType, "application/json") {
		return fmt.Errorf("Expected json, got [%s] with code %d", contentType, res.StatusCode)
	}
	return nil
}

// PrettyError returns a human-readable error message for common http errors returned by the API.
func PrettyError(noun, verb string, res *http.Response) error {
	if res.StatusCode == 404 {
		return fmt.Errorf("%s does not exist, %s failed.", noun, verb)
	} else if res.StatusCode == 401 {
		return fmt.Errorf("%s %s failed, permission denied. Check your API key.", noun, verb)
	} else if res.StatusCode == 403 {
		// This is what happens if an endpoint URL gets typo'd. (dgray 2015-09-14)
		return fmt.Errorf("%s %s failed. Are you using the right API path?", noun, verb)
	}
	return nil
}

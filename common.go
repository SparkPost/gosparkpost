package gosparkpost

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	re "regexp"
	"strings"

	certifi "github.com/certifi/gocertifi"
)

// TODO: define paths statically in each endpoint's file, move version to Config
// TODO: rename e.g. Transmissions.Create to Client.SendTransmission

// Config includes all information necessary to make an API request.
type Config struct {
	BaseUrl    string
	ApiKey     string
	ApiVersion int
}

// Client contains connection and authentication information.
// Specifying your own http.Client gives you lots of control over how connections are made.
type Client struct {
	Config *Config
	Client *http.Client
}

var nonDigit *re.Regexp = re.MustCompile(`\D`)

// NewConfig builds a Config object using the provided map.
func NewConfig(m map[string]string) (*Config, error) {
	c := &Config{}

	if baseurl, ok := m["baseurl"]; ok {
		c.BaseUrl = baseurl
	} else {
		return nil, fmt.Errorf("BaseUrl is required for api config")
	}

	if apikey, ok := m["apikey"]; ok {
		c.ApiKey = apikey
	} else {
		return nil, fmt.Errorf("ApiKey is required for api config")
	}

	return c, nil
}

// Response contains information about the last HTTP response.
// Helpful when an error message doesn't necessarily give the complete picture.
type Response struct {
	HTTP    *http.Response
	Body    []byte
	Results map[string]interface{} `json:"results,omitempty"`
	Errors  []Error                `json:"errors,omitempty"`
}

// Error mirrors the error format returned by SparkPost APIs.
type Error struct {
	Message     string `json:"message"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Part        string `json:"part,omitempty"`
	Line        int    `json:"line,omitempty"`
}

func (e Error) Json() (string, error) {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// Init pulls together everything necessary to make an API request.
// Caller may provide their own http.Client by setting it in the provided API object.
func (api *Client) Init(cfg *Config) error {
	// Set default values
	if cfg.BaseUrl == "" {
		cfg.BaseUrl = "https://api.sparkpost.com"
	}
	if cfg.ApiVersion == 0 {
		cfg.ApiVersion = 1
	}
	api.Config = cfg

	if api.Client == nil {
		// Ran into an issue where USERTrust was not recognized on OSX.
		// The rest of this block was the fix.

		// load Mozilla cert pool
		pool, err := certifi.CACerts()
		if err != nil {
			return err
		}

		// configure transport using Mozilla cert pool
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
		}

		// configure http client using transport
		api.Client = &http.Client{Transport: transport}
	}

	return nil
}

// HttpPost sends a Post request with the provided JSON payload to the specified url.
// Query params are supported via net/url - roll your own and stringify it.
// Authenticate using the configured API key.
func (c *Client) HttpPost(url string, data []byte) (*Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Config.ApiKey)
	res, err := c.Client.Do(req)
	ares := &Response{HTTP: res}
	return ares, err
}

// HttpGet sends a Get request to the specified url.
// Query params are supported via net/url - roll your own and stringify it.
// Authenticate using the configured API key.
func (c *Client) HttpGet(url string) (*Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Config.ApiKey)
	res, err := c.Client.Do(req)
	ares := &Response{HTTP: res}
	return ares, err
}

// HttpDelete sends a Delete request to the provided url.
// Query params are supported via net/url - roll your own and stringify it.
// Authenticate using the configured API key.
func (c *Client) HttpDelete(url string) (*Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Config.ApiKey)
	res, err := c.Client.Do(req)
	ares := &Response{HTTP: res}
	return ares, err
}

// ReadBody is a convenience method that returns the http.Response body.
// The first time this function is called, the body is read from the
// http.Response. For subsequent calls, the cached version in
// Response.Body is returned.
func (r *Response) ReadBody() ([]byte, error) {
	// Calls 2+ to this function for the same http.Response will now DWIM
	if r.Body != nil {
		return r.Body, nil
	}

	defer r.HTTP.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.HTTP.Body)
	r.Body = bodyBytes
	return bodyBytes, err
}

// ParseResponse pulls info from JSON http responses into api.Response object.
// It's helpful to call Response.AssertJson before calling this function.
func (r *Response) ParseResponse() error {
	body, err := r.ReadBody()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return fmt.Errorf("Failed to parse API response: [%s]\n%s", err, string(body))
	}

	return nil
}

// AssertObject asserts that the provided variable is a map[string]something.
// The string parameter is used to customize the generated error message.
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

var jctype string = "application/json"

// AssertJson returns an error if the provided HTTP response isn't JSON.
func (r *Response) AssertJson() error {
	if r.HTTP == nil {
		return fmt.Errorf("AssertJson got nil http.Response")
	}
	ctype := r.HTTP.Header.Get("Content-Type")
	// allow things like "application/json; charset=utf-8" in addition to the bare content type
	if !(strings.EqualFold(ctype, jctype) || strings.HasPrefix(ctype, jctype)) {
		return fmt.Errorf("Expected json, got [%s] with code %d", ctype, r.HTTP.StatusCode)
	}
	return nil
}

// PrettyError returns a human-readable error message for common http errors returned by the API.
// The string parameters are used to customize the generated error message
// (example: noun=template, verb=create).
func (r *Response) PrettyError(noun, verb string) error {
	if r.HTTP == nil {
		return nil
	}
	code := r.HTTP.StatusCode
	if code == 404 {
		return fmt.Errorf("%s does not exist, %s failed.", noun, verb)
	} else if code == 401 {
		return fmt.Errorf("%s %s failed, permission denied. Check your API key.", noun, verb)
	} else if code == 403 {
		// This is what happens if an endpoint URL gets typo'd.
		return fmt.Errorf("%s %s failed. Are you using the right API path?", noun, verb)
	}
	return nil
}

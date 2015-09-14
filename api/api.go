package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	certifi "github.com/certifi/gocertifi"
)

type Config struct {
	BaseUrl string
	ApiKey  string
}

type API struct {
	Config *Config
	Client *http.Client
}

type Response struct {
	HTTP    *http.Response
	Body    string
	Results map[string]string `json:"results,omitempty"`
	Errors  []Error           `json:"errors,omitempty"`
}

type Error struct {
	Message     string `json:"message"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Part        string `json:"part,omitempty"`
	Line        int    `json:"line,omitempty"`
}

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

// Send a Post request with the provided JSON payload to the specified url.
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

// Send a Get request to the specified url.
// Query params are supported via net/url - pass in url.String()
// Authenticate using the configured API key.
func (api *API) HttpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", api.Config.ApiKey)
	return api.Client.Do(req)
}

// Send a Delete request to the provided url.
// Authenticate using the configured API key.
func (api *API) HttpDelete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", api.Config.ApiKey)
	return api.Client.Do(req)
}

func ReadBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func ParseApiResponse(res *http.Response) (*Response, error) {
	body, err := ReadBody(res)
	if err != nil {
		return nil, err
	}

	apiRes := &Response{}
	err = json.Unmarshal(body, apiRes)
	if err != nil {
		return nil, err
	}
	apiRes.Body = string(body)
	apiRes.HTTP = res

	return apiRes, nil
}

// Return an error if the provided HTTP response isn't JSON.
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

func PrettyError(noun, verb string, res *http.Response) error {
	if res.StatusCode == 404 {
		return fmt.Errorf("%s does not exist, %s failed.", noun, verb)
	} else if res.StatusCode == 401 {
		return fmt.Errorf("%s %s failed, permission denied. Check your API key.", noun, verb)
	} else if res.StatusCode == 403 {
		return fmt.Errorf("%s %s failed. Are you using the right API path?", noun, verb)
	}
	return nil
}

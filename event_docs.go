package gosparkpost

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

var eventDocumentationFormat = "/api/v%d/webhooks/events/documentation"

type EventGroup struct {
	Name        string
	Events      map[string]EventMeta `json:"events"`
	Description string               `json:"description"`
	DisplayName string               `json:"display_name"`
}

type EventMeta struct {
	Name        string
	Fields      map[string]EventField `json:"event"`
	Description string                `json:"description"`
	DisplayName string                `json:"display_name"`
}

type EventField struct {
	Description string      `json:"description"`
	SampleValue interface{} `json:"sampleValue"`
}

func (c *Client) EventDocumentation() (g map[string]*EventGroup, res *Response, err error) {
	path := fmt.Sprintf(eventDocumentationFormat, c.Config.ApiVersion)
	res, err = c.HttpGet(c.Config.BaseUrl + path)
	if err != nil {
		return nil, nil, err
	}

	if err = res.AssertJson(); err != nil {
		return nil, res, err
	}

	if res.HTTP.StatusCode == 200 {
		var body []byte
		var ok bool
		body, err = res.ReadBody()
		if err != nil {
			return nil, res, err
		}

		var results map[string]map[string]*EventGroup
		var groups map[string]*EventGroup
		if err = json.Unmarshal(body, &results); err != nil {
			return nil, res, err
		} else if groups, ok = results["results"]; ok {
			return groups, res, err
		}
		return nil, res, errors.New("Unexpected response format")
	} else {
		err = res.ParseResponse()
		if err != nil {
			return nil, res, err
		}
		if len(res.Errors) > 0 {
			err = res.PrettyError("EventDocumentation", "retrieve")
			if err != nil {
				return nil, res, err
			}
		}
		return nil, res, errors.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return nil, res, err
}

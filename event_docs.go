package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

var EventDocumentationFormat = "/api/v%d/webhooks/events/documentation"

type EventGroups struct {
	Groups map[string]*EventGroup `json:"groups"`

	Context context.Context `json:"-"`
}

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

func (c *Client) EventDocumentation(eg *EventGroups) (res *Response, err error) {
	path := fmt.Sprintf(EventDocumentationFormat, c.Config.ApiVersion)
	res, err = c.HttpGet(context.TODO(), c.Config.BaseUrl+path)
	if err != nil {
		return nil, err
	}

	if err = res.AssertJson(); err != nil {
		return res, err
	}

	if res.HTTP.StatusCode == 200 {
		var body []byte
		var ok bool
		body, err = res.ReadBody()
		if err != nil {
			return res, err
		}

		var results map[string]map[string]*EventGroup
		if err = json.Unmarshal(body, &results); err != nil {
			return res, err
		} else if eg.Groups, ok = results["results"]; ok {
			return res, err
		}
		return res, errors.New("Unexpected response format")
	} else {
		err = res.ParseResponse()
		if err != nil {
			return res, err
		}
		if len(res.Errors) > 0 {
			err = res.PrettyError("EventDocumentation", "retrieve")
			if err != nil {
				return res, err
			}
		}
		return res, errors.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return res, err
}

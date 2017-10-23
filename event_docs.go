package gosparkpost

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

var EventDocumentationFormat = "/api/v%d/webhooks/events/documentation"

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
	return c.EventDocumentationContext(context.Background())
}

func (c *Client) EventDocumentationContext(ctx context.Context) (groups map[string]*EventGroup, res *Response, err error) {
	path := fmt.Sprintf(EventDocumentationFormat, c.Config.ApiVersion)
	var results map[string]map[string]*EventGroup
	res, err = c.HttpGetJson(ctx, c.Config.BaseUrl+path, &results)
	if err != nil {
		return
	}

	var ok bool
	if groups, ok = results["results"]; ok {
		// Success!
	} else {
		err = errors.New("Unexpected response format (results)")
	}

	return
}

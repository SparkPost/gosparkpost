package gosparkpost

var eventDocumentationFormat = "/api/v%d/webhooks/events/documentation"

type EventDocumentationResponse struct {
	Results map[string]*EventGroup `json:"results,omitempty"`
}

type EventGroup struct {
	Name        string
	Events      map[string]EventField
	Description string `json:"description"`
	DisplayName string `json:"display_name"`
}

type EventField struct {
	Description string `json:"description"`
	SampleValue string `json:"sampleValue"`
}

func (c *Client) EventDocumentation() (g *EventGroup, res *Response, err error) {
	return nil, nil, nil
}

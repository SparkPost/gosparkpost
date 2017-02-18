package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/SparkPost/gosparkpost/events"
)

// https://www.sparkpost.com/api#/reference/message-events
var (
	MessageEventsPathFormat        = "/api/v%d/message-events"
	MessageEventsSamplesPathFormat = "/api/v%d/message-events/events/samples"
)

type EventsPage struct {
	client *Client

	Events     events.Events
	TotalCount int
	Errors     []interface{}

	NextPage  string
	PrevPage  string
	FirstPage string
	LastPage  string

	Params map[string]string `json:"-"`
}

// https://developers.sparkpost.com/api/#/reference/message-events/events-samples/search-for-message-events
func (c *Client) MessageEventsSearch(ep *EventsPage) (*Response, error) {
	return c.MessageEventsSearchContext(context.Background(), ep)
}

// MessageEventsSearchContext is the same as MessageEventsSearch, and it accepts a context.Context
func (c *Client) MessageEventsSearchContext(ctx context.Context, ep *EventsPage) (*Response, error) {
	path := fmt.Sprintf(MessageEventsPathFormat, c.Config.ApiVersion)
	url, err := url.Parse(c.Config.BaseUrl + path)
	if err != nil {
		return nil, err
	}

	if len(ep.Params) > 0 {
		q := url.Query()
		for k, v := range ep.Params {
			q.Add(k, v)
		}
		url.RawQuery = q.Encode()
	}

	// Send off our request
	res, err := c.HttpGet(ctx, url.String())
	if err != nil {
		return res, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return res, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodyBytes, ep)
	if err != nil {
		return res, err
	}

	ep.client = c

	return res, nil
}

// Next returns the next page of results from a previous MessageEventsSearch call
func (ep *EventsPage) Next() (*EventsPage, *Response, error) {
	return ep.NextContext(context.Background())
}

// NextContext is the same as Next, and it accepts a context.Context
func (ep *EventsPage) NextContext(ctx context.Context) (*EventsPage, *Response, error) {
	if ep.NextPage == "" {
		return nil, nil, nil
	}

	// Send off our request
	res, err := ep.client.HttpGet(ctx, ep.client.Config.BaseUrl+ep.NextPage)
	if err != nil {
		return nil, res, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return nil, res, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return nil, res, err
	}

	var eventsPage EventsPage
	err = json.Unmarshal(bodyBytes, &eventsPage)
	if err != nil {
		return nil, res, err
	}

	eventsPage.client = ep.client

	return &eventsPage, res, nil
}

func (ep *EventsPage) UnmarshalJSON(data []byte) error {
	// Clear object.
	*ep = EventsPage{}

	// Object with array of events and cursors is being sent on Message Events.
	var resultsWrapper struct {
		RawEvents  []json.RawMessage `json:"results"`
		TotalCount int               `json:"total_count,omitempty"`
		Links      []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links,omitempty"`
		Errors []interface{} `json:"errors,omitempty"`
	}
	err := json.Unmarshal(data, &resultsWrapper)
	if err != nil {
		return err
	}

	ep.Events, err = events.ParseRawJSONEvents(resultsWrapper.RawEvents)
	if err != nil {
		return err
	}

	ep.Errors = resultsWrapper.Errors
	ep.TotalCount = resultsWrapper.TotalCount

	for _, link := range resultsWrapper.Links {
		switch link.Rel {
		case "next":
			ep.NextPage = link.Href
		case "previous":
			ep.PrevPage = link.Href
		case "first":
			ep.FirstPage = link.Href
		case "last":
			ep.LastPage = link.Href
		}
	}

	return nil
}

// EventSamples requests a list of example event data.
func (c *Client) EventSamples(types *[]string) (*events.Events, *Response, error) {
	return c.EventSamplesContext(context.Background(), types)
}

// EventSamplesContext is the same as EventSamples, and it accepts a context.Context
func (c *Client) EventSamplesContext(ctx context.Context, types *[]string) (*events.Events, *Response, error) {
	path := fmt.Sprintf(MessageEventsSamplesPathFormat, c.Config.ApiVersion)
	url, err := url.Parse(c.Config.BaseUrl + path)
	if err != nil {
		return nil, nil, err
	}

	// Filter out types.
	if types != nil {
		// validate types
		for _, etype := range *types {
			if !events.ValidEventType(etype) {
				return nil, nil, fmt.Errorf("Invalid event type [%s]", etype)
			}
		}

		// get the query string object so we can modify it
		q := url.Query()
		// add the requested events and re-encode
		q.Set("events", strings.Join(*types, ","))
		url.RawQuery = q.Encode()
	}

	// Send off our request
	res, err := c.HttpGet(ctx, url.String())
	if err != nil {
		return nil, res, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return nil, res, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return nil, res, err
	}

	var events events.Events
	err = json.Unmarshal(bodyBytes, &events)
	if err != nil {
		return nil, res, err
	}

	return &events, res, nil
}

// ParseEvents function is left only for backward-compatibility. Events are parsed by events pkg.
func ParseEvents(rawEventsPtr []*json.RawMessage) (*[]events.Event, error) {
	rawEvents := make([]json.RawMessage, len(rawEventsPtr))
	for i, ptr := range rawEventsPtr {
		rawEvents[i] = *ptr
	}

	events, err := events.ParseRawJSONEvents(rawEvents)
	if err != nil {
		return nil, err
	}
	return &events, nil
}

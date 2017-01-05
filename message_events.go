package gosparkpost

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/SparkPost/gosparkpost/events"
)

// https://www.sparkpost.com/api#/reference/message-events
var (
	ErrEmptyPage                   = errors.New("empty page")
	MessageEventsPathFormat        = "/api/v%d/message-events"
	MessageEventsSamplesPathFormat = "/api/v%d/message-events/events/samples"
)

type EventsPage struct {
	client *Client

	Events     events.Events
	TotalCount int
	nextPage   string
	prevPage   string
	firstPage  string
	lastPage   string
}

// https://developers.sparkpost.com/api/#/reference/message-events/events-samples/search-for-message-events
func (c *Client) MessageEvents(params map[string]string) (*EventsPage, *Response, error) {
	path := fmt.Sprintf(MessageEventsPathFormat, c.Config.ApiVersion)
	url, err := url.Parse(c.Config.BaseUrl + path)
	if err != nil {
		return nil, nil, err
	}

	if len(params) > 0 {
		q := url.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		url.RawQuery = q.Encode()
	}

	// Send off our request
	res, err := c.HttpGet(context.TODO(), url.String())
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

	eventsPage.client = c

	return &eventsPage, res, nil
}

func (events *EventsPage) Next() (*EventsPage, error) {
	if events.nextPage == "" {
		return nil, ErrEmptyPage
	}

	// Send off our request
	res, err := events.client.HttpGet(context.TODO(), events.client.Config.BaseUrl+events.nextPage)
	if err != nil {
		return nil, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return nil, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return nil, err
	}

	var eventsPage EventsPage
	err = json.Unmarshal(bodyBytes, &eventsPage)
	if err != nil {
		return nil, err
	}

	eventsPage.client = events.client

	return &eventsPage, nil
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
	}
	err := json.Unmarshal(data, &resultsWrapper)
	if err != nil {
		return err
	}

	ep.Events, err = events.ParseRawJSONEvents(resultsWrapper.RawEvents)
	if err != nil {
		return err
	}

	ep.TotalCount = resultsWrapper.TotalCount

	for _, link := range resultsWrapper.Links {
		switch link.Rel {
		case "next":
			ep.nextPage = link.Href
		case "previous":
			ep.prevPage = link.Href
		case "first":
			ep.firstPage = link.Href
		case "last":
			ep.lastPage = link.Href
		}
	}

	return nil
}

// Samples requests a list of example event data.
func (c *Client) EventSamples(types *[]string) (*events.Events, error) {
	path := fmt.Sprintf(MessageEventsSamplesPathFormat, c.Config.ApiVersion)
	url, err := url.Parse(c.Config.BaseUrl + path)
	if err != nil {
		return nil, err
	}

	// Filter out types.
	if types != nil {
		// validate types
		for _, etype := range *types {
			if !events.ValidEventType(etype) {
				return nil, fmt.Errorf("Invalid event type [%s]", etype)
			}
		}

		// get the query string object so we can modify it
		q := url.Query()
		// add the requested events and re-encode
		q.Set("events", strings.Join(*types, ","))
		url.RawQuery = q.Encode()
	}

	// Send off our request
	res, err := c.HttpGet(context.TODO(), url.String())
	if err != nil {
		return nil, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return nil, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return nil, err
	}

	var events events.Events
	err = json.Unmarshal(bodyBytes, &events)
	if err != nil {
		return nil, err
	}

	return &events, nil
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

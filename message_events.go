package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/SparkPost/gosparkpost/events"
	"github.com/pkg/errors"
)

// https://www.sparkpost.com/api#/reference/message-events
var (
	MessageEventsPathFormat        = "/api/v%d/message-events"
	MessageEventsSamplesPathFormat = "/api/v%d/message-events/events/samples"
)

type EventsPage struct {
	Events     events.Events
	TotalCount int
	Errors     []interface{}

	NextPage  string
	PrevPage  string
	FirstPage string
	LastPage  string

	Client *Client           `json:"-"`
	Params map[string]string `json:"-"`
}

// https://developers.sparkpost.com/api/#/reference/message-events/events-samples/search-for-message-events
func (c *Client) MessageEventsSearch(page *EventsPage) (*Response, error) {
	return c.MessageEventsSearchContext(context.Background(), page)
}

// MessageEventsSearchContext is the same as MessageEventsSearch, and it accepts a context.Context
func (c *Client) MessageEventsSearchContext(ctx context.Context, page *EventsPage) (*Response, error) {
	path := fmt.Sprintf(MessageEventsPathFormat, c.Config.ApiVersion)
	url, err := url.Parse(c.Config.BaseUrl + path)
	if err != nil {
		return nil, errors.Wrap(err, "parsing url")
	}

	if len(page.Params) > 0 {
		q := url.Query()
		for k, v := range page.Params {
			q.Add(k, v)
		}
		url.RawQuery = q.Encode()
	}

	// Send off our request
	res, err := c.HttpGetJson(ctx, url.String(), page)
	if err != nil {
		return res, err
	}

	page.Client = c

	return res, nil
}

// Next returns the next page of results from a previous MessageEventsSearch call
func (page *EventsPage) Next() (*EventsPage, *Response, error) {
	return page.NextContext(context.Background())
}

// NextContext is the same as Next, and it accepts a context.Context
func (page *EventsPage) NextContext(ctx context.Context) (*EventsPage, *Response, error) {
	var nextPage EventsPage
	if page.NextPage == "" {
		return nil, nil, nil
	}

	// Send off our request
	res, err := page.Client.HttpGetJson(ctx, page.Client.Config.BaseUrl+page.NextPage, &nextPage)
	if err != nil {
		return nil, res, err
	}

	nextPage.Client = page.Client

	return &nextPage, res, nil
}

func (page *EventsPage) UnmarshalJSON(data []byte) error {
	// Clear object.
	*page = EventsPage{}

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

	page.Events, err = events.ParseRawJSONEvents(resultsWrapper.RawEvents)
	if err != nil {
		return err
	}

	page.Errors = resultsWrapper.Errors
	page.TotalCount = resultsWrapper.TotalCount

	for _, link := range resultsWrapper.Links {
		switch link.Rel {
		case "next":
			page.NextPage = link.Href
		case "previous":
			page.PrevPage = link.Href
		case "first":
			page.FirstPage = link.Href
		case "last":
			page.LastPage = link.Href
		}
	}

	return nil
}

// EventSamples requests a list of example event data.
func (c *Client) EventSamples(types []string) (*events.Events, *Response, error) {
	return c.EventSamplesContext(context.Background(), types)
}

// EventSamplesContext is the same as EventSamples, and it accepts a context.Context
func (c *Client) EventSamplesContext(ctx context.Context, types []string) (*events.Events, *Response, error) {
	path := fmt.Sprintf(MessageEventsSamplesPathFormat, c.Config.ApiVersion)
	url, err := url.Parse(c.Config.BaseUrl + path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parsing url")
	}

	// Filter out types.
	if types != nil {
		// validate types
		for _, etype := range types {
			if !events.ValidEventType(etype) {
				return nil, nil, fmt.Errorf("Invalid event type [%s]", etype)
			}
		}

		// get the query string object so we can modify it
		q := url.Query()
		// add the requested events and re-encode
		q.Set("events", strings.Join(types, ","))
		url.RawQuery = q.Encode()
	}

	// Send off our request
	var events events.Events
	res, err := c.HttpGetJson(ctx, url.String(), &events)
	if err != nil {
		return nil, res, err
	}

	return &events, res, nil
}

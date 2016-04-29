package gosparkpost

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/SparkPost/gosparkpost/events"
)

var (
	ErrEmptyPage                   = errors.New("empty page")
	messageEventsPathFormat        = "%s/api/v%d/message-events"
	messageEventsSamplesPathFormat = "%s/api/v%d/message-events/events/samples"
)

// EventsPage contains a list of events, with links to any additional pages of results.
// https://developers.sparkpost.com/api/#/reference/message-events
type EventsPage struct {
	headers map[string]string
	client  *Client

	Events     events.Events
	TotalCount int
	nextPage   string
	prevPage   string
	firstPage  string
	lastPage   string
}

// MessageEvents searches for event data matching the specified params.
func (c *Client) MessageEvents(params map[string]string) (*EventsPage, error) {
	return c.MessageEventsWithHeaders(params, nil)
}

// MessageEventsWithHeaders searches for event data matching the specified params, and allows passing in extra HTTP headers.
func (c *Client) MessageEventsWithHeaders(params, headers map[string]string) (*EventsPage, error) {
	url, err := url.Parse(fmt.Sprintf(messageEventsPathFormat, c.Config.BaseUrl, c.Config.ApiVersion))
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := url.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		url.RawQuery = q.Encode()
	}

	// Send off our request
	res, err := c.HttpGet(url.String(), headers)
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

	eventsPage.client = c
	eventsPage.headers = headers

	return &eventsPage, nil
}

// Next fetches the next page of events for the current query, resending any HTTP headers specified with the original request.
func (events *EventsPage) Next() (*EventsPage, error) {
	if events.nextPage == "" {
		return nil, ErrEmptyPage
	}

	// Send off our request
	res, err := events.client.HttpGet(events.client.Config.BaseUrl+events.nextPage, events.headers)
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
	eventsPage.headers = events.headers

	return &eventsPage, nil
}

type resultsWrapper struct {
	RawEvents  []json.RawMessage `json:"results"`
	TotalCount int               `json:"total_count,omitempty"`
	Links      []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links,omitempty"`
}

// UnmarshalJSON parses the provided []byte, extracting event data into ep.
func (ep *EventsPage) UnmarshalJSON(data []byte) error {
	// Clear object.
	*ep = EventsPage{}

	var resWrapper resultsWrapper
	// Object with array of events and cursors is being sent on Message Events.
	err := json.Unmarshal(data, &resWrapper)
	if err != nil {
		return err
	}

	ep.Events, err = events.ParseRawJSONEvents(resWrapper.RawEvents)
	if err != nil {
		return err
	}

	ep.TotalCount = resWrapper.TotalCount

	for _, link := range resWrapper.Links {
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

// EventSamples requests a list of example event data.
func (c *Client) EventSamples(types *[]string) (*events.Events, error) {
	return c.EventSamplesWithHeaders(types, nil)
}

// EventSamplesWithHeaders requests a list of example event data, and allows passing in extra HTTP headers.
func (c *Client) EventSamplesWithHeaders(types *[]string, headers map[string]string) (*events.Events, error) {
	url, err := url.Parse(fmt.Sprintf(messageEventsSamplesPathFormat, c.Config.BaseUrl, c.Config.ApiVersion))
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
	res, err := c.HttpGet(url.String(), headers)
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

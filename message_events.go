package gosparkpost

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/SparkPost/gosparkpost/events"
)

// https://www.sparkpost.com/api#/reference/message-events
var (
	messageEventsPathFormat        = "%s/api/v%d/message-events"
	messageEventsSamplesPathFormat = "%s/api/v%d/message-events/events/samples"
)

// https://developers.sparkpost.com/api/#/reference/message-events/events-samples/search-for-message-events
func (c *Client) SearchMessageEvents(params map[string]string) (*events.EventsPage, error) {
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
	res, err := c.HttpGet(url.String())
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

	var eventsPage events.EventsPage
	err = json.Unmarshal(bodyBytes, &eventsPage)
	if err != nil {
		return nil, err
	}

	return &eventsPage, nil
}

// Samples requests a list of example event data.
func (c *Client) EventSamples(types *[]string) (*events.Events, error) {
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
	res, err := c.HttpGet(url.String())
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

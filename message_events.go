package gosparkpost

import (
	"encoding/json"
	"fmt"
	URL "net/url"
	"strings"
)

// https://www.sparkpost.com/api#/reference/message-events
var messageEventsPathFormat = "/api/v%d/message-events"

type EventsWrapper struct {
	Results []struct {
		Type string `json:"type"`

		CustomerID      string            `json:"customer_id"`
		DelvMethod      string            `json:"delv_method,omitempty"`
		EventID         string            `json:"event_id"`
		FriendlyFrom    string            `json:"friendly_from,omitempty"`
		IPAddress       string            `json:"ip_address,omitempty"`
		MessageID       string            `json:"message_id,omitempty"`
		MsgFrom         string            `json:"msg_from"`
		MsgSize         string            `json:"msg_size,omitempty"`
		NumRetries      string            `json:"num_retries,omitempty"`
		QueueTime       string            `json:"queue_time,omitempty"`
		RawRcptTo       string            `json:"raw_rcpt_to"`
		RcptMeta        map[string]string `json:"rcpt_meta,omitempty"`
		RcptTags        []interface{}     `json:"rcpt_tags,omitempty"`
		RcptTo          string            `json:"rcpt_to"`
		RoutingDomain   string            `json:"routing_domain,omitempty"`
		Subject         string            `json:"subject,omitempty"`
		Tdate           string            `json:"tdate,omitempty"`
		TemplateID      string            `json:"template_id,omitempty"`
		TemplateVersion string            `json:"template_version,omitempty"`
		TransmissionID  string            `json:"transmission_id,omitempty"`
		Timestamp       string            `json:"timestamp"`
		ErrorCode       string            `json:"error_code,omitempty"`
		RawReason       string            `json:"raw_reason,omitempty"`
		Reason          string            `json:"reason,omitempty"`
		RemoteAddr      string            `json:"remote_addr,omitempty"`
	} `json:"results,omitempty"`
	TotalCount int      `json:"total_count,omitempty"`
	Links      []string `json:"links,omitempty"`
}

// https://developers.sparkpost.com/api/#/reference/message-events/events-samples/search-for-message-events
func (c *Client) SearchMessageEvents(parameters map[string]string) (*EventsWrapper, error) {

	var finalUrl string
	path := fmt.Sprintf(messageEventsPathFormat, c.Config.ApiVersion)
	if parameters == nil || len(parameters) == 0 {
		finalUrl = fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	} else {
		params := URL.Values{}
		for k, v := range parameters {
			params.Add(k, v)
		}

		finalUrl = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, params.Encode())
	}

	fmt.Println("URL: ", finalUrl)

	return DoRequest(c, finalUrl)
}

// Samples requests a list of example event data.
func (c *Client) EventSamples(types *[]string) (*EventsWrapper, error) {
	// append any requested event types to path
	var finalUrl string
	path := fmt.Sprintf(messageEventsPathFormat, c.Config.ApiVersion)
	if types == nil {
		finalUrl = fmt.Sprintf("%s%s/events/samples", c.Config.BaseUrl, path)
	} else {

		// break up the url into a net.URL object
		u, err := URL.Parse(fmt.Sprintf("%s%s/events/samples", c.Config.BaseUrl, path))
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}

		// get the query string object so we can modify it
		q := u.Query()
		// add the requested events and re-encode
		q.Set("events", strings.Join(*types, ","))
		u.RawQuery = q.Encode()
		finalUrl = u.String()
	}

	return DoRequest(c, finalUrl)
}

func DoRequest(c *Client, finalUrl string) (*EventsWrapper, error) {
	// Send off our request
	res, err := c.HttpGet(finalUrl)
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

	/*// DEBUG
	err = iou.WriteFile("./events.json", bodyBytes, 0644)
	if err != nil {
		return nil, err
	}
	*/

	myVal := string(bodyBytes)

	fmt.Println(myVal)

	// Parse expected response structure
	var resMap EventsWrapper //map[string]interface{}
	err = json.Unmarshal(bodyBytes, &resMap)
	if err != nil {
		// FIXME: better error message
		fmt.Println("Failed to unmarshal content")
		return nil, err
	}

	return &resMap, err
}

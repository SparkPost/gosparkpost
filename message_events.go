package gosparkpost

import (
	"encoding/json"
	"fmt"
	URL "net/url"
	"strings"
)

// https://www.sparkpost.com/api#/reference/message-events
var messageEventsPathFormat = "/api/v%d/message-events"

type EventItem struct {
	Type            string            `json:"type"`
	CampaignID      string            `json:"campaign_id"`
	CustomerID      string            `json:"customer_id"`
	DeliveryMethod  string            `json:"delv_method,omitempty"`
	DeviceToken     string            `json:"device_token,omitempty"`
	EventID         string            `json:"event_id"`
	ErrorCode       string            `json:"error_code"`
	FriendlyFrom    string            `json:"friendly_from,omitempty"`
	IPAddress       string            `json:"ip_address,omitempty"`
	MessageID       string            `json:"message_id,omitempty"`
	MessageFrom     string            `json:"msg_from"`
	MessageSize     string            `json:"msg_size,omitempty"`
	QueueTime       string            `json:"queue_time,omitempty"`
	RawReason       string            `json:"raw_reason,omitempty"`
	Reason          string            `json:"reason,omitempty"`
	ReceiveProtocol string            `json:"recv_method,omitempty"`
	RawRcptTo       string            `json:"raw_rcpt_to,omitempty"`
	Metadata        map[string]string `json:"rcpt_meta,omitempty"`
	RcptTags        []interface{}     `json:"rcpt_tags,omitempty"`
	Recipient       string            `json:"rcpt_to,omitempty"`
	RecipientType   string            `json:"rcpt_type,omitempty"`
	Retries         string            `json:"num_retries,omitempty"`
	RoutingDomain   string            `json:"routing_domain,omitempty"`
	Subject         string            `json:"subject,omitempty"`
	Tags            []string          `json:"rcpt_tags,omitempty"`
	Tdate           string            `json:"tdate,omitempty"`
	TemplateID      string            `json:"template_id,omitempty"`
	TemplateVersion string            `json:"template_version,omitempty"`
	TransmissionID  string            `json:"transmission_id,omitempty"`
	Timestamp       string            `json:"timestamp,omitempty"`
	RemoteAddr      string            `json:"remote_addr,omitempty"`
	Binding         string            `json:"binding,omitempty"`
	BindingGroup    string            `json:"binding_group,omitempty"`
	BounceClass     string            `json:"bounce_class,omitempty"`
	TargetLinkName  string            `json:"target_link_name,omitempty"`
	TargetLinkURL   string            `json:"target_link_url,omitempty"`
	UserAgent       string            `json:"user_agent,omitempty"`
	UserID          string            `json:"user_id,omitempty"`
	Submitted       string            `json:"submitted_rcpts,omitempty"`
	InjectionMethod string            `json:"inj_method,omitempty"`
	Accepted        string            `json:"accepted_rcpts,omitempty"`
	RelayID         string            `json:"relay_id,omitempty"`
	Pathway         string            `json:"pathway,omitempty"`
	PathwayGroup    string            `json:"pathway_group,omitempty"`
	From            string            `json:"msg_from,omitempty"`
	To              string            `json:"rcpt_to,omitempty"`
	WebhookID       string            `json:"webhook_id,omitempty"`
	UserString      string            `json:"user_str"`
	ReportedBy      string            `json:"report_by"`
	ReportedTo      string            `json:"report_to"`
	FeedbackType    string            `json:"fbtype"`
}

type EventsWrapper struct {
	Results    []*EventItem  `json:"results,omitempty"`
	TotalCount int           `json:"total_count,omitempty"`
	Links      []string      `json:"links,omitempty"`
	Errors     []interface{} `json:"errors,omitempty"`
	//{"errors":[{"param":"from","message":"From must be before to","value":"2014-07-20T09:00"},{"param":"to","message":"To must be in the format YYYY-MM-DDTHH:mm","value":"now"}]}
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

func (c *Client) EventAsString(e *EventItem) string {
	eventType := e.Type

	switch eventType {
	case "bounce":
		return fmt.Sprintf("%s (%s) %s %s => %s %s: %s",
			e.Timestamp, e.Type, e.TransmissionID, e.Binding, e.Recipient,
			e.BounceClass, e.RawReason)
	case "click":
		return fmt.Sprintf("%s (%s) %s %s => %s",
			e.Timestamp, e.Type, e.TransmissionID, e.Recipient, e.TargetLinkURL)
	case "creation":
		return fmt.Sprintf("%s (%s) %s (%s, %s)",
			e.Timestamp, e.Type, e.TransmissionID, e.Submitted, e.Accepted)
	case "delay":
		return fmt.Sprintf("%s (%s) %s => %s %s: %s",
			e.Timestamp, e.Type, e.MessageFrom, e.Recipient, e.BounceClass, e.RawReason)
	case "delivery":
		return fmt.Sprintf("%s (%s) %s %s => %s",
			e.Timestamp, e.Type, e.TransmissionID, e.Binding, e.Recipient)
	case "generation_failure":
		return fmt.Sprintf("%s (%s) %s %s => %s %s: %s",
			e.Timestamp, e.Type, e.TransmissionID, e.Binding, e.Recipient,
			e.ErrorCode, e.RawReason)
	case "generation_rejection":
		return fmt.Sprintf("%s (%s) %s %s => %s %s: %s",
			e.Timestamp, e.Type, e.TransmissionID, e.Binding, e.Recipient,
			e.ErrorCode, e.RawReason)
	case "injection":
		return fmt.Sprintf("%s (%s) %s %s %s %s %s %s %s",
			e.Timestamp, e.Type, e.MessageID, e.Recipient, e.MessageFrom,
			e.MessageSize, e.ReceiveProtocol,
			e.Binding, e.BindingGroup)
	case "list_unsubscribe":
		return fmt.Sprintf("%s (%s) %s %s: [%s]",
			e.Timestamp, e.Type, e.TransmissionID, e.Recipient, e.CampaignID)
	case "link_unsubscribe":
		return fmt.Sprintf("%s (%s) %s %s: [%s]",
			e.Timestamp, e.Type, e.TransmissionID, e.Recipient, e.CampaignID)
	case "open":
		return fmt.Sprintf("%s (%s) %s %s",
			e.Timestamp, e.Type, e.TransmissionID, e.Recipient)
	case "out_of_band":
		return fmt.Sprintf("%s (%s) [%s] %s => %s %s: %s",
			e.Timestamp, e.Type, e.CampaignID, e.Binding, e.Recipient,
			e.BounceClass, e.RawReason)
	case "policy_rejection":
		return fmt.Sprintf("%s (%s) %s [%s] => %s %s: %s",
			e.Timestamp, e.Type, e.TransmissionID, e.CampaignID, e.Recipient,
			e.ErrorCode, e.RawReason)
	case "spam_complaint":
		return fmt.Sprintf("%s (%s) %s %s %s => %s (%s)",
			e.Timestamp, e.Type, e.TransmissionID, e.Binding, e.ReportedBy, e.ReportedTo, e.Recipient)
	case "relay_delivery":
		return fmt.Sprintf("%s (%s) %s %s <= %s",
			e.Timestamp, e.Type, e.RelayID, e.Binding, e.MessageFrom)
	case "relay_injection":
		return fmt.Sprintf("%s (%s) %s %s %s => %s",
			e.Timestamp, e.Type, e.RelayID, e.Binding, e.MessageFrom, e.Recipient)
	case "relay_message":
		return fmt.Sprintf("%s (%s) %s => %s (%s)", e.Timestamp, e.Type, e.From, e.To, e.WebhookID)
	case "relay_permfail":
		return fmt.Sprintf("%s (%s) %s %s <= %s %s: %s",
			e.Timestamp, e.Type, e.RelayID, e.Binding, e.MessageFrom, e.ErrorCode, e.RawReason)
	case "relay_rejection":
		return fmt.Sprintf("%s (%s) %s %s => %s %s: %s",
			e.Timestamp, e.Type, e.RelayID, e.MessageFrom, e.Recipient, e.ErrorCode, e.RawReason)
	case "relay_tempfail":
		return fmt.Sprintf("%s (%s) %s %s <= %s %s: %s",
			e.Timestamp, e.Type, e.RelayID, e.Binding, e.MessageFrom, e.ErrorCode, e.RawReason)

	default:
		return fmt.Sprintf("%s UNKNOWN(%s)",
			e.Timestamp, e.Type)
	}
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

	// Parse expected response structure
	var resMap EventsWrapper //map[string]interface{}
	err = json.Unmarshal(bodyBytes, &resMap)
	if err != nil {
		return nil, err
	}

	return &resMap, err
}

// Package events defines a struct for each type of event and provides various other helper functions.
package events

// eventTypes contains all of the valid event types
var eventTypes = map[string]bool{
	//"creation":             false,

	"delivery":             true,
	"injection":            true,
	"bounce":               true,
	"delay":                true,
	"out_of_band":          true,
	"open":                 true,
	"click":                true,
	"generation_failure":   true,
	"generation_rejection": true,
	"list_unsubscribe":     true,
	"link_unsubscribe":     true,
	"policy_rejection":     false,
	"spam_complaint":       false,
	"relay_delivery":       false,
	"relay_injection":      false,
	"relay_permfail":       false,
	"relay_rejection":      false,
	"relay_tempfail":       false,
}

// ValidEventType returns true if the event name parameter is valid.
func ValidEventType(eventType string) bool {
	if _, ok := eventTypes[eventType]; ok {
		return true
	}
	return false
}

// EventForName returns a struct matching the passed-in type.
func EventForName(eventType string) Event {
	if !ValidEventType(eventType) {
		return nil
	}

	switch eventType {
	case "bounce":
		return &Bounce{}
	case "click":
		return &Click{}
	case "delay":
		return &Delay{}
	case "delivery":
		return &Delivery{}
	case "generation_failure":
		return &GenerationFailure{}
	case "generation_rejection":
		return &GenerationRejection{}
	case "injection":
		return &Injection{}
	case "list_unsubscribe":
		return &ListUnsubscribe{}
	case "link_unsubscribe":
		return &LinkUnsubscribe{}
	case "open":
		return &Open{}
	case "out_of_band":
		return &OutOfBand{}

	default:
		return nil
	}
}

// All event types must satisfy this interface, so we can have heterogenous Event slices.
type Event interface {
	EventType() string
}

// EventCommon contains fields common to all types of Event objects
type EventCommon struct {
	Type string `json:"type"`
}

func (e EventCommon) EventType() string { return e.Type }

type Bounce struct {
	EventCommon
	Binding         string      `json:"binding"`
	BindingGroup    string      `json:"binding_group"`
	BounceClass     string      `json:"bounce_class"`
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	DeliveryMethod  string      `json:"delv_method"`
	DeviceToken     string      `json:"device_token"`
	ErrorCode       string      `json:"error_code"`
	IPAddress       string      `json:"ip_address"`
	MessageID       string      `json:"message_id"`
	MessageFrom     string      `json:"msg_from"`
	MessageSize     string      `json:"msg_size"`
	Retries         string      `json:"num_retries"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            []string    `json:"rcpt_tags"`
	Recipient       string      `json:"rcpt_to"`
	RecipientType   string      `json:"rcpt_type"`
	RawReason       string      `json:"raw_reason"`
	Reason          string      `json:"reason"`
	ReceiveProtocol string      `json:"recv_method"`
	RoutingDomain   string      `json:"routing_domain"`
	Subject         string      `json:"subject"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

type GeoIP struct {
	Country   string `json:"country"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Click struct {
	EventCommon
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	DeliveryMethod  string      `json:"delv_method"`
	GeoIP           *GeoIP      `json:"geo_ip"`
	IPAddress       string      `json:"ip_address"`
	MessageID       string      `json:"message_id"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            []string    `json:"rcpt_tags"`
	Recipient       string      `json:"rcpt_to"`
	RecipientType   string      `json:"rcpt_type"`
	TargetLinkName  string      `json:"target_link_name"`
	TargetLinkURL   string      `json:"target_link_url"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserAgent       string      `json:"user_agent"`
}

type Delay struct {
	EventCommon
	Binding         string      `json:"binding"`
	BindingGroup    string      `json:"binding_group"`
	BounceClass     string      `json:"bounce_class"`
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	DeliveryMethod  string      `json:"delv_method"`
	DeviceToken     string      `json:"device_token"`
	ErrorCode       string      `json:"error_code"`
	IPAddress       string      `json:"ip_address"`
	MessageID       string      `json:"message_id"`
	MessageFrom     string      `json:"msg_from"`
	MessageSize     string      `json:"msg_size"`
	Retries         string      `json:"num_retries"`
	QueueTime       string      `json:"queue_time"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            []string    `json:"rcpt_tags"`
	Recipient       string      `json:"rcpt_to"`
	RecipientType   string      `json:"rcpt_type"`
	RawReason       string      `json:"raw_reason"`
	Reason          string      `json:"reason"`
	RoutingDomain   string      `json:"routing_domain"`
	Subject         string      `json:"subject"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

type Delivery struct {
	EventCommon
	Binding         string      `json:"binding"`
	BindingGroup    string      `json:"binding_group"`
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	DeliveryMethod  string      `json:"delv_method"`
	DeviceToken     string      `json:"device_token"`
	IPAddress       string      `json:"ip_address"`
	MessageID       string      `json:"message_id"`
	MessageFrom     string      `json:"msg_from"`
	MessageSize     string      `json:"msg_size"`
	Retries         string      `json:"num_retries"`
	QueueTime       string      `json:"queue_time"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            []string    `json:"rcpt_tags"`
	Recipient       string      `json:"rcpt_to"`
	RecipientType   string      `json:"rcpt_type"`
	ReceiveProtocol string      `json:"recv_method"`
	RoutingDomain   string      `json:"routing_domain"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

type GenerationFailure struct {
	EventCommon
	Binding          string      `json:"binding"`
	BindingGroup     string      `json:"binding_group"`
	CampaignID       string      `json:"campaign_id"`
	CustomerID       string      `json:"customer_id"`
	ErrorCode        string      `json:"error_code"`
	Metadata         interface{} `json:"rcpt_meta"`
	SubstitutionData interface{} `json:"rcpt_subs"`
	Tags             []string    `json:"rcpt_tags"`
	Recipient        string      `json:"rcpt_to"`
	RawReason        string      `json:"raw_reason"`
	Reason           string      `json:"reason"`
	ReceiveProtocol  string      `json:"recv_method"`
	RoutingDomain    string      `json:"routing_domain"`
	TemplateID       string      `json:"template_id"`
	TemplateVersion  string      `json:"template_version"`
	Timestamp        int64       `json:"timestamp"`
	TransmissionID   string      `json:"transmission_id"`
}

type GenerationRejection GenerationFailure

type Injection struct {
	EventCommon
	Binding         string      `json:"binding"`
	BindingGroup    string      `json:"binding_group"`
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	MessageID       string      `json:"message_id"`
	MessageFrom     string      `json:"msg_from"`
	MessageSize     string      `json:"msg_size"`
	Metadata        interface{} `json:"rcpt_meta"`
	Pathway         string      `json:"pathway"`
	PathwayGroup    string      `json:"pathway_group"`
	Tags            []string    `json:"rcpt_tags"`
	Recipient       string      `json:"rcpt_to"`
	RecipientType   string      `json:"rcpt_type"`
	ReceiveProtocol string      `json:"recv_method"`
	RoutingDomain   string      `json:"routing_domain"`
	Subject         string      `json:"subject"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

type ListUnsubscribe struct {
	EventCommon
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	MessageFrom     string      `json:"mailfrom"`
	MessageID       string      `json:"message_id"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            []string    `json:"rcpt_tags"`
	Recipient       string      `json:"rcpt_to"`
	RecipientType   string      `json:"rcpt_type"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

type LinkUnsubscribe struct {
	EventCommon
	ListUnsubscribe
	UserAgent string `json:"user_agent"`
}

type Open struct {
	EventCommon
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	DeliveryMethod  string      `json:"delv_method"`
	GeoIP           *GeoIP      `json:"geo_ip"`
	IPAddress       string      `json:"ip_address"`
	MessageID       string      `json:"message_id"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            []string    `json:"rcpt_tags"`
	Recipient       string      `json:"rcpt_to"`
	RecipientType   string      `json:"rcpt_type"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserAgent       string      `json:"user_agent"`
}

type OutOfBand struct {
	EventCommon
	Binding         string `json:"binding"`
	BindingGroup    string `json:"binding_group"`
	BounceClass     string `json:"bounce_class"`
	CampaignID      string `json:"campaign_id"`
	CustomerID      string `json:"customer_id"`
	DeliveryMethod  string `json:"delv_method"`
	DeviceToken     string `json:"device_token"`
	ErrorCode       string `json:"error_code"`
	MessageID       string `json:"message_id"`
	MessageFrom     string `json:"msg_from"`
	Recipient       string `json:"rcpt_to"`
	RawReason       string `json:"raw_reason"`
	Reason          string `json:"reason"`
	ReceiveProtocol string `json:"recv_method"`
	RoutingDomain   string `json:"routing_domain"`
	TemplateID      string `json:"template_id"`
	TemplateVersion string `json:"template_version"`
	Timestamp       int64  `json:"timestamp"`
}

type PolicyRejection struct {
	EventCommon
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	ErrorCode       string      `json:"error_code"`
	MessageID       string      `json:"message_id"`
	MessageFrom     string      `json:"msg_from"`
	Metadata        interface{} `json:"rcpt_meta"`
	Pathway         string      `json:"pathway"`
	PathwayGroup    string      `json:"pathway_group"`
	Tags            []string    `json:"rcpt_tags"`
	RawReason       string      `json:"raw_reason"`
	Reason          string      `json:"reason"`
	Recipient       string      `json:"rcpt_to"`
	RecipientType   string      `json:"rcpt_type"`
	ReceiveProtocol string      `json:"recv_method"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

type RelayInjection struct {
	EventCommon
	Binding         string `json:"binding"`
	BindingGroup    string `json:"binding_group"`
	CustomerID      string `json:"customer_id"`
	MessageFrom     string `json:"msg_from"`
	MessageSize     string `json:"msg_size"`
	Pathway         string `json:"pathway"`
	PathwayGroup    string `json:"pathway_group"`
	Recipient       string `json:"rcpt_to"`
	ReceiveProtocol string `json:"recv_method"`
	RelayID         string `json:"relay_id"`
	RoutingDomain   string `json:"routing_domain"`
	Timestamp       int64  `json:"timestamp"`
}

type RelayRejection struct {
	EventCommon
	CustomerID      string `json:"customer_id"`
	ErrorCode       string `json:"error_code"`
	MessageFrom     string `json:"msg_from"`
	Pathway         string `json:"pathway"`
	PathwayGroup    string `json:"pathway_group"`
	RawReason       string `json:"raw_reason"`
	Reason          string `json:"reason"`
	Recipient       string `json:"rcpt_to"`
	ReceiveProtocol string `json:"recv_method"`
	RelayID         string `json:"relay_id"`
	RemoteAddress   string `json:"remote_addr"`
	Timestamp       int64  `json:"timestamp"`
}

type RelayDelivery struct {
	EventCommon
	Binding         string `json:"binding"`
	BindingGroup    string `json:"binding_group"`
	CustomerID      string `json:"customer_id"`
	DeliveryMethod  string `json:"delv_method"`
	MessageFrom     string `json:"msg_from"`
	Pathway         string `json:"pathway"`
	PathwayGroup    string `json:"pathway_group"`
	QueueTime       string `json:"queue_time"`
	ReceiveProtocol string `json:"recv_method"`
	RelayID         string `json:"relay_id"`
	Retries         string `json:"num_retries"`
	RoutingDomain   string `json:"routing_domain"`
	Timestamp       int64  `json:"timestamp"`
}

type RelayTempfail struct {
	EventCommon
	Binding         string `json:"binding"`
	BindingGroup    string `json:"binding_group"`
	CustomerID      string `json:"customer_id"`
	DeliveryMethod  string `json:"delv_method"`
	ErrorCode       string `json:"error_code"`
	MessageFrom     string `json:"msg_from"`
	Retries         string `json:"num_retries"`
	QueueTime       string `json:"queue_time"`
	Pathway         string `json:"pathway"`
	PathwayGroup    string `json:"pathway_group"`
	RawReason       string `json:"raw_reason"`
	Reason          string `json:"reason"`
	ReceiveProtocol string `json:"recv_method"`
	RelayID         string `json:"relay_id"`
	RoutingDomain   string `json:"routing_domain"`
	Timestamp       int64  `json:"timestamp"`
}

type RelayPermfail RelayTempfail

type SpamComplaint struct {
	EventCommon
	Binding         string      `json:"binding"`
	BindingGroup    string      `json:"binding_group"`
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	DeliveryMethod  string      `json:"delv_method"`
	FeedbackType    string      `json:"fbtype"`
	FriendlyFrom    string      `json:"friendly_from"`
	MessageID       string      `json:"message_id"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            []string    `json:"rcpt_tags"`
	Recipient       string      `json:"rcpt_to"`
	RecipientType   string      `json:"rcpt_type"`
	ReportedBy      string      `json:"report_by"`
	ReportedTo      string      `json:"report_to"`
	Subject         string      `json:"subject"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserString      string      `json:"user_str"`
}

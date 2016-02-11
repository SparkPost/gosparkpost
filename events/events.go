// Package events defines a struct for each type of event and provides various other helper functions.
package events

import "fmt"

// eventTypes contains all of the valid event types
var eventTypes = map[string]bool{
	"bounce":               true,
	"click":                true,
	"creation":             false,
	"delay":                true,
	"delivery":             true,
	"generation_failure":   true,
	"generation_rejection": true,
	"injection":            true,
	"list_unsubscribe":     true,
	"link_unsubscribe":     true,
	"open":                 true,
	"out_of_band":          true,
	"policy_rejection":     true,
	"spam_complaint":       true,
	"relay_delivery":       true,
	"relay_injection":      true,
	"relay_message":        true,
	"relay_permfail":       true,
	"relay_rejection":      true,
	"relay_tempfail":       true,
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
	case "creation":
		return &Creation{}
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
	case "policy_rejection":
		return &PolicyRejection{}
	case "spam_complaint":
		return &SpamComplaint{}
	case "relay_delivery":
		return &RelayDelivery{}
	case "relay_injection":
		return &RelayInjection{}
	case "relay_message":
		return &RelayMessage{}
	case "relay_permfail":
		return &RelayPermfail{}
	case "relay_rejection":
		return &RelayRejection{}
	case "relay_tempfail":
		return &RelayTempfail{}

	default:
		return nil
	}
}

// Event allows 2+ types of event in a data structure.
type Event interface {
	EventType() string
}

func ECLog(e Event) string {
	// XXX: this feels like the wrong way; can't figure out the right way
	switch e.(type) {
	case *Bounce:
		b := e.(*Bounce)
		return b.ECLog()
	case *Delay:
		d := e.(*Delay)
		return d.ECLog()
	case *Delivery:
		d := e.(*Delivery)
		return d.ECLog()
	case *Injection:
		i := e.(*Injection)
		return i.ECLog()
	case *OutOfBand:
		o := e.(*OutOfBand)
		return o.ECLog()
	}
	return ""
}

type ECLogger interface {
	ECLog() string
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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

// String returns a brief summary of a Bounce event
func (b *Bounce) String() string {
	return fmt.Sprintf("%s B %s %s => %s %s: %s",
		b.Timestamp, b.TransmissionID, b.Binding, b.Recipient,
		b.BounceClass, b.RawReason)
}

// ECLog emits a Bounce in the same format that it would be logged to bouncelog.ec:
// https://support.messagesystems.com/docs/web-ref/log_formats.version_3.php
func (b *Bounce) ECLog() string {
	return fmt.Sprintf("%s@%s@@@B@%s@%s@%s@%s@@%s@%s@%s@%s",
		b.Timestamp, b.MessageID, b.Recipient, b.MessageFrom,
		b.Binding, b.BindingGroup, b.BounceClass, b.MessageSize,
		b.IPAddress, b.RawReason)
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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserAgent       string      `json:"user_agent"`
}

// String returns a brief summary of a Click event
func (c *Click) String() string {
	return fmt.Sprintf("%s C %s %s => %s",
		c.Timestamp, c.TransmissionID, c.Recipient, c.TargetLinkURL)
}

type Creation struct {
	EventCommon
	Accepted        string      `json:"accepted_rcpts"`
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	InjectionMethod string      `json:"inj_method"`
	NodeName        string      `json:"node_name"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            []string    `json:"rcpt_tags"`
	Submitted       string      `json:"submitted_rcpts"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserID          string      `json:"user_id"`
}

func (c *Creation) String() string {
	return fmt.Sprintf("%s CT %s (%s, %s)",
		c.Timestamp, c.TransmissionID, c.Submitted, c.Accepted)
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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

// String returns a brief summary of a Delay event
func (d *Delay) String() string {
	return fmt.Sprintf("%s T %s => %s %s: %s",
		d.Timestamp, d.MessageFrom, d.Recipient, d.BounceClass, d.RawReason)
}

// ECLog emits a Delay in the same format that it would be logged to bouncelog.ec:
// https://support.messagesystems.com/docs/web-ref/log_formats.version_3.php
func (d *Delay) ECLog() string {
	return fmt.Sprintf("%s@%s@@@T@%s@%s@%s@%s@@%s@%s@%s@%s",
		d.Timestamp, d.MessageID, d.Recipient, d.MessageFrom,
		d.Binding, d.BindingGroup, d.BounceClass, d.MessageSize,
		d.IPAddress, d.RawReason)
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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

// String returns a brief summary of a Delivery event
func (d *Delivery) String() string {
	return fmt.Sprintf("%s D %s %s => %s",
		d.Timestamp, d.TransmissionID, d.Binding, d.Recipient)
}

// ECLog emits a Delivery in the same format that it would be logged to mainlog.ec:
// https://support.messagesystems.com/docs/web-ref/log_formats.version_3.php
func (d *Delivery) ECLog() string {
	return fmt.Sprintf("%s@%s@@@D@%s@%s@%s@%s@%s@%s@%s",
		d.Timestamp, d.MessageID, d.RoutingDomain, d.MessageSize,
		d.Binding, d.BindingGroup, d.Retries, d.QueueTime, d.IPAddress)
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
	Timestamp        string      `json:"timestamp"`
	TransmissionID   string      `json:"transmission_id"`
}

// String returns a brief summary of a GenerationFailure event
func (g *GenerationFailure) String() string {
	return fmt.Sprintf("%s GF %s %s => %s %s: %s",
		g.Timestamp, g.TransmissionID, g.Binding, g.Recipient,
		g.ErrorCode, g.RawReason)
}

type GenerationRejection GenerationFailure

// String returns a brief summary of a GenerationFailure event
func (g *GenerationRejection) String() string {
	return fmt.Sprintf("%s GR %s %s => %s %s: %s",
		g.Timestamp, g.TransmissionID, g.Binding, g.Recipient,
		g.ErrorCode, g.RawReason)
}

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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

// String returns a brief summary of a GenerationFailure event
func (i *Injection) String() string {
	return fmt.Sprintf("%s R %s %s => %s",
		i.Timestamp, i.TransmissionID, i.Binding, i.Recipient)
}

// ECLog emits an Injection in the same format that it would be logged to mainlog.ec:
// https://support.messagesystems.com/docs/web-ref/log_formats.version_3.php
func (i *Injection) ECLog() string {
	return fmt.Sprintf("%s@%s@@@R@%s@%s@@%s@%s@%s@%s",
		i.Timestamp, i.MessageID, i.Recipient, i.MessageFrom,
		i.MessageSize, i.ReceiveProtocol,
		i.Binding, i.BindingGroup)
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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

// String returns a brief summary of a ListUnsubscribe event
func (l *ListUnsubscribe) String() string {
	return fmt.Sprintf("%s U %s %s: [%s]",
		l.Timestamp, l.TransmissionID, l.Recipient, l.CampaignID)
}

type LinkUnsubscribe struct {
	EventCommon
	ListUnsubscribe
	UserAgent string `json:"user_agent"`
}

// String returns a brief summary of a ListUnsubscribe event
func (l *LinkUnsubscribe) String() string {
	return fmt.Sprintf("%s LU %s %s: [%s]",
		l.Timestamp, l.TransmissionID, l.Recipient, l.CampaignID)
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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserAgent       string      `json:"user_agent"`
}

// String returns a brief summary of an Open event
func (o *Open) String() string {
	return fmt.Sprintf("%s O %s %s",
		o.Timestamp, o.TransmissionID, o.Recipient)
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
	Timestamp       string `json:"timestamp"`
}

// String returns a brief summary of a Bounce event
func (b *OutOfBand) String() string {
	return fmt.Sprintf("%s OOB [%s] %s => %s %s: %s",
		b.Timestamp, b.CampaignID, b.Binding, b.Recipient,
		b.BounceClass, b.RawReason)
}

// ECLog emits an OutOfBand in the same format that it would be logged to bouncelog.ec:
// https://support.messagesystems.com/docs/web-ref/log_formats.version_3.php
func (b *OutOfBand) ECLog() string {
	return fmt.Sprintf("%s@%s@@@B@%s@%s@%s@%s@@%s@@@%s",
		b.Timestamp, b.MessageID, b.Recipient, b.MessageFrom,
		b.Binding, b.BindingGroup, b.BounceClass, b.RawReason)
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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

// String returns a brief summary of a PolicyRejection event
func (p *PolicyRejection) String() string {
	return fmt.Sprintf("%s PR %s [%s] => %s %s: %s",
		p.Timestamp, p.TransmissionID, p.CampaignID, p.Recipient,
		p.ErrorCode, p.RawReason)
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
	Timestamp       string `json:"timestamp"`
}

// String returns a brief summary of a RelayInjection event
func (i *RelayInjection) String() string {
	return fmt.Sprintf("%s RI %s %s %s => %s",
		i.Timestamp, i.RelayID, i.Binding, i.MessageFrom, i.Recipient)
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
	Timestamp       string `json:"timestamp"`
}

// String returns a brief summary of a RelayInjection event
func (r *RelayRejection) String() string {
	return fmt.Sprintf("%s RR %s %s => %s %s: %s",
		r.Timestamp, r.RelayID, r.MessageFrom, r.Recipient, r.ErrorCode, r.RawReason)
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
	Timestamp       string `json:"timestamp"`
}

// String returns a brief summary of a RelayDelivery event
func (d *RelayDelivery) String() string {
	return fmt.Sprintf("%s RD %s %s <= %s",
		d.Timestamp, d.RelayID, d.Binding, d.MessageFrom)
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
	Timestamp       string `json:"timestamp"`
}

// String returns a brief summary of a RelayTempfail event
func (t *RelayTempfail) String() string {
	return fmt.Sprintf("%s RT %s %s <= %s %s: %s",
		t.Timestamp, t.RelayID, t.Binding, t.MessageFrom, t.ErrorCode, t.RawReason)
}

type RelayPermfail RelayTempfail

// String returns a brief summary of a RelayInjection event
func (p *RelayPermfail) String() string {
	return fmt.Sprintf("%s RP %s %s <= %s %s: %s",
		p.Timestamp, p.RelayID, p.Binding, p.MessageFrom, p.ErrorCode, p.RawReason)
}

type RelayContent struct {
	HTML    string              `json:"html"`
	Text    string              `json:"text"`
	Subject string              `json:"subject"`
	To      []string            `json:"to"`
	Cc      []string            `json:"cc"`
	Headers []map[string]string `json:"headers"`
	Email   string              `json:"email_rfc822"`
	Base64  bool                `json:email_rfc822_is_base64"`
}

type RelayMessage struct {
	EventCommon
	Content      RelayContent `json:"content"`
	FriendlyFrom string       `json:"friendly_from"`
	From         string       `json:"msg_from"`
	To           string       `json:"rcpt_to"`
	WebhookID    string       `json:"webhook_id"`
}

func (m *RelayMessage) String() string {
	return fmt.Sprintf("%s => %s (%s)", m.From, m.To, m.WebhookID)
}

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
	Timestamp       string      `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserString      string      `json:"user_str"`
}

// String returns a brief summary of a SpamComplaint event
func (p *SpamComplaint) String() string {
	return fmt.Sprintf("%s S %s %s %s => %s (%s)",
		p.Timestamp, p.TransmissionID, p.Binding, p.ReportedBy, p.ReportedTo, p.Recipient)
}

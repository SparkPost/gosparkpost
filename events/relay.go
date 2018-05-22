package events

import "fmt"

type RelayInjection struct {
	EventCommon
	Binding         string    `json:"binding"`
	BindingGroup    string    `json:"binding_group"`
	CustomerID      string    `json:"customer_id"`
	MessageFrom     string    `json:"msg_from"`
	MessageSize     string    `json:"msg_size"`
	Pathway         string    `json:"pathway"`
	PathwayGroup    string    `json:"pathway_group"`
	Recipient       string    `json:"rcpt_to"`
	ReceiveProtocol string    `json:"recv_method"`
	RelayID         string    `json:"relay_id"`
	RoutingDomain   string    `json:"routing_domain"`
	Timestamp       Timestamp `json:"timestamp"`
}

// String returns a brief summary of a RelayInjection event
func (i *RelayInjection) String() string {
	return fmt.Sprintf("%s RI %s %s %s => %s",
		i.Timestamp, i.RelayID, i.Binding, i.MessageFrom, i.Recipient)
}

type RelayRejection struct {
	EventCommon
	CustomerID      string    `json:"customer_id"`
	ErrorCode       string    `json:"error_code"`
	MessageFrom     string    `json:"msg_from"`
	Pathway         string    `json:"pathway"`
	PathwayGroup    string    `json:"pathway_group"`
	RawReason       string    `json:"raw_reason"`
	Reason          string    `json:"reason"`
	Recipient       string    `json:"rcpt_to"`
	ReceiveProtocol string    `json:"recv_method"`
	RelayID         string    `json:"relay_id"`
	RemoteAddress   string    `json:"remote_addr"`
	Timestamp       Timestamp `json:"timestamp"`
}

// String returns a brief summary of a RelayInjection event
func (r *RelayRejection) String() string {
	return fmt.Sprintf("%s RR %s %s => %s %s: %s",
		r.Timestamp, r.RelayID, r.MessageFrom, r.Recipient, r.ErrorCode, r.RawReason)
}

type RelayDelivery struct {
	EventCommon
	Binding         string    `json:"binding"`
	BindingGroup    string    `json:"binding_group"`
	CustomerID      string    `json:"customer_id"`
	DeliveryMethod  string    `json:"delv_method"`
	MessageFrom     string    `json:"msg_from"`
	Pathway         string    `json:"pathway"`
	PathwayGroup    string    `json:"pathway_group"`
	QueueTime       string    `json:"queue_time"`
	ReceiveProtocol string    `json:"recv_method"`
	RelayID         string    `json:"relay_id"`
	Retries         string    `json:"num_retries"`
	RoutingDomain   string    `json:"routing_domain"`
	Timestamp       Timestamp `json:"timestamp"`
}

// String returns a brief summary of a RelayDelivery event
func (d *RelayDelivery) String() string {
	return fmt.Sprintf("%s RD %s %s <= %s",
		d.Timestamp, d.RelayID, d.Binding, d.MessageFrom)
}

type RelayTempfail struct {
	EventCommon
	Binding         string    `json:"binding"`
	BindingGroup    string    `json:"binding_group"`
	CustomerID      string    `json:"customer_id"`
	DeliveryMethod  string    `json:"delv_method"`
	ErrorCode       string    `json:"error_code"`
	MessageFrom     string    `json:"msg_from"`
	Retries         string    `json:"num_retries"`
	QueueTime       string    `json:"queue_time"`
	Pathway         string    `json:"pathway"`
	PathwayGroup    string    `json:"pathway_group"`
	RawReason       string    `json:"raw_reason"`
	Reason          string    `json:"reason"`
	ReceiveProtocol string    `json:"recv_method"`
	RelayID         string    `json:"relay_id"`
	RoutingDomain   string    `json:"routing_domain"`
	Timestamp       Timestamp `json:"timestamp"`
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
	Base64  bool                `json:"email_rfc822_is_base64"`
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

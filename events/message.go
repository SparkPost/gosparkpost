package events

import "fmt"

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
	Timestamp       Timestamp   `json:"timestamp"`
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
	Timestamp       Timestamp   `json:"timestamp"`
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

type Bounce struct {
	EventCommon
	Binding         string            `json:"binding"`
	BindingGroup    string            `json:"binding_group"`
	BounceClass     string            `json:"bounce_class"`
	CampaignID      string            `json:"campaign_id"`
	CustomerID      string            `json:"customer_id"`
	DeliveryMethod  string            `json:"delv_method"`
	DeviceToken     string            `json:"device_token"`
	ErrorCode       string            `json:"error_code"`
	IPAddress       string            `json:"ip_address"`
	MessageID       string            `json:"message_id"`
	MessageFrom     string            `json:"msg_from"`
	MessageSize     string            `json:"msg_size"`
	Retries         string            `json:"num_retries"`
	Metadata        map[string]string `json:"rcpt_meta"`
	Tags            []string          `json:"rcpt_tags"`
	Recipient       string            `json:"rcpt_to"`
	RecipientType   string            `json:"rcpt_type"`
	RawReason       string            `json:"raw_reason"`
	Reason          string            `json:"reason"`
	ReceiveProtocol string            `json:"recv_method"`
	RoutingDomain   string            `json:"routing_domain"`
	Subject         string            `json:"subject"`
	TemplateID      string            `json:"template_id"`
	TemplateVersion string            `json:"template_version"`
	Timestamp       Timestamp         `json:"timestamp"`
	TransmissionID  string            `json:"transmission_id"`
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

type OutOfBand struct {
	EventCommon
	Binding         string    `json:"binding"`
	BindingGroup    string    `json:"binding_group"`
	BounceClass     string    `json:"bounce_class"`
	CampaignID      string    `json:"campaign_id"`
	CustomerID      string    `json:"customer_id"`
	DeliveryMethod  string    `json:"delv_method"`
	DeviceToken     string    `json:"device_token"`
	ErrorCode       string    `json:"error_code"`
	MessageID       string    `json:"message_id"`
	MessageFrom     string    `json:"msg_from"`
	Recipient       string    `json:"rcpt_to"`
	RawReason       string    `json:"raw_reason"`
	Reason          string    `json:"reason"`
	ReceiveProtocol string    `json:"recv_method"`
	RoutingDomain   string    `json:"routing_domain"`
	TemplateID      string    `json:"template_id"`
	TemplateVersion string    `json:"template_version"`
	Timestamp       Timestamp `json:"timestamp"`
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
	Timestamp       Timestamp   `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserString      string      `json:"user_str"`
}

// String returns a brief summary of a SpamComplaint event
func (p *SpamComplaint) String() string {
	return fmt.Sprintf("%s S %s %s %s => %s (%s)",
		p.Timestamp, p.TransmissionID, p.Binding, p.ReportedBy, p.ReportedTo, p.Recipient)
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
	Timestamp       Timestamp   `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

// String returns a brief summary of a PolicyRejection event
func (p *PolicyRejection) String() string {
	return fmt.Sprintf("%s PR %s [%s] => %s %s: %s",
		p.Timestamp, p.TransmissionID, p.CampaignID, p.Recipient,
		p.ErrorCode, p.RawReason)
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
	Timestamp       Timestamp   `json:"timestamp"`
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

type SMSStatus struct {
	EventCommon
	CustomerID     string `json:"customer_id"`
	DeliveryMethod string `json:"delv_method"`
	// TODO: `json:"dr_latency"`
	IPAddress      string   `json:"ip_address"`
	RawReason      string   `json:"raw_reason"`
	Reason         string   `json:"reason"`
	RoutingDomain  string   `json:"routing_domain"`
	Destination    string   `json:"sms_dst"`
	DestinationNPI string   `json:"sms_dst_npi"`
	DestinationTON string   `json:"sms_dst_ton"`
	RemoteIDs      []string `json:"sms_remoteids"`
	Source         string   `json:"sms_src"`
	SourceNPI      string   `json:"sms_src_npi"`
	SourceTON      string   `json:"sms_src_ton"`
	Text           string   `json:"sms_text"`
	StatusType     string   `json:"stat_type"`
	StatusState    string   `json:"stat_state"`
	// TODO: SubAccountID string `json:"subaccount_id"`
	Timestamp Timestamp `json:"timestamp"`
}

// String returns a brief summary of a Delay event
func (e *SMSStatus) String() string {
	return fmt.Sprintf("%s SMS", e.Timestamp) // TODO: Improve message.
}

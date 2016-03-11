package events

import "fmt"

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
	Timestamp       Timestamp   `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserAgent       string      `json:"user_agent"`
}

// String returns a brief summary of a Click event
func (c *Click) String() string {
	return fmt.Sprintf("%s C %s %s => %s",
		c.Timestamp, c.TransmissionID, c.Recipient, c.TargetLinkURL)
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
	Timestamp       Timestamp   `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
	UserAgent       string      `json:"user_agent"`
}

// String returns a brief summary of an Open event
func (o *Open) String() string {
	return fmt.Sprintf("%s O %s %s",
		o.Timestamp, o.TransmissionID, o.Recipient)
}

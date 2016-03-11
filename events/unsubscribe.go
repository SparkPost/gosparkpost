package events

import "fmt"

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
	Timestamp       Timestamp   `json:"timestamp"`
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

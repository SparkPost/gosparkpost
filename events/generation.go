package events

import "fmt"

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
	Timestamp        Timestamp   `json:"timestamp"`
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

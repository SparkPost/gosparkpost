// Package events defines a struct for each type of event and provides various other helper functions.
package events

// eventTypes contains all of the valid event types
var eventTypes = map[string]bool{
	"creation":             true,
	"delivery":             true,
	"injection":            true,
	"bounce":               true,
	"delay":                true,
	"policy_rejection":     true,
	"out_of_band":          true,
	"open":                 true,
	"click":                true,
	"generation_failure":   true,
	"generation_rejection": true,
	"spam_complaint":       true,
	"list_unsubscribe":     true,
	"link_unsubscribe":     true,
	"relay_delivery":       true,
	"relay_injection":      true,
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

type BounceEvent struct {
	BounceClass     string      `json:"bounce_class"`
	CampaignID      string      `json:"campaign_id"`
	CustomerID      string      `json:"customer_id"`
	ErrorCode       string      `json:"error_code"`
	IPAddress       string      `json:"ip_address"`
	MessageID       string      `json:"message_id"`
	MessageFrom     string      `json:"msg_from"`
	MessageSize     string      `json:"msg_size"`
	NumRetries      int         `json:"num_retries"`
	Metadata        interface{} `json:"rcpt_meta"`
	Tags            interface{} `json:"rcpt_tags"`
	Recipient       string      `json:"rcpt_to"`
	Reason          string      `json:"reason"`
	RoutingDomain   string      `json:"routing_domain"`
	TemplateID      string      `json:"template_id"`
	TemplateVersion string      `json:"template_version"`
	Timestamp       int64       `json:"timestamp"`
	TransmissionID  string      `json:"transmission_id"`
}

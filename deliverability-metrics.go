package gosparkpost

import (
	"encoding/json"
	"fmt"

	URL "net/url"
)

// https://www.sparkpost.com/api#/reference/message-events
var deliverabilityMetricPathFormat = "/api/v%d/metrics/deliverability"

// DeliverabilityMetricItem contains all of the metrics returned from the Metrics endpoint.
type DeliverabilityMetricItem struct {
	CountInjected               int    `json:"count_injected"`
	CountBounce                 int    `json:"count_bounce,omitempty"`
	CountRejected               int    `json:"count_rejected,omitempty"`
	CountDelivered              int    `json:"count_delivered,omitempty"`
	CountDeliveredFirst         int    `json:"count_delivered_first,omitempty"`
	CountDeliveredSubsequent    int    `json:"count_delivered_subsequent,omitempty"`
	TotalDeliveryTimeFirst      int    `json:"total_delivery_time_first,omitempty"`
	TotalDeliveryTimeSubsequent int    `json:"total_delivery_time_subsequent,omitempty"`
	TotalMsgVolume              int    `json:"total_msg_volume,omitempty"`
	CountPolicyRejection        int    `json:"count_policy_rejection,omitempty"`
	CountGenerationRejection    int    `json:"count_generation_rejection,omitempty"`
	CountGenerationFailed       int    `json:"count_generation_failed,omitempty"`
	CountInbandBounce           int    `json:"count_inband_bounce,omitempty"`
	CountOutofbandBounce        int    `json:"count_outofband_bounce,omitempty"`
	CountSoftBounce             int    `json:"count_soft_bounce,omitempty"`
	CountHardBounce             int    `json:"count_hard_bounce,omitempty"`
	CountBlockBounce            int    `json:"count_block_bounce,omitempty"`
	CountAdminBounce            int    `json:"count_admin_bounce,omitempty"`
	CountUndeterminedBounce     int    `json:"count_undetermined_bounce,omitempty"`
	CountDelayed                int    `json:"count_delayed,omitempty"`
	CountDelayedFirst           int    `json:"count_delayed_first,omitempty"`
	CountRendered               int    `json:"count_rendered,omitempty"`
	CountUniqueRendered         int    `json:"count_unique_rendered,omitempty"`
	CountUniqueConfirmedOpened  int    `json:"count_unique_confirmed_opened,omitempty"`
	CountClicked                int    `json:"count_clicked,omitempty"`
	CountUniqueClicked          int    `json:"count_unique_clicked,omitempty"`
	CountTargeted               int    `json:"count_targeted,omitempty"`
	CountSent                   int    `json:"count_sent,omitempty"`
	CountAccepted               int    `json:"count_accepted,omitempty"`
	CountSpamComplaint          int    `json:"count_spam_complaint,omitempty"`
	Domain                      string `json:"domain,omitempty"`
	CampaignId                  string `json:"campaign_id,omitempty"`
	TemplateId                  string `json:"template_id,omitempty"`
	TimeStamp                   string `json:"ts,omitempty"`
	WatchedDomain               string `json:"watched_domain,omitempty"`
	Binding                     string `json:"binding,omitempty"`
	BindingGroup                string `json:"binding_group,omitempty"`
}

// DeliverabilityMetricEventsWrapper is a pagination container for DeliverabilityMetricItem.
type DeliverabilityMetricEventsWrapper struct {
	Results    []*DeliverabilityMetricItem `json:"results,omitempty"`
	TotalCount int                         `json:"total_count,omitempty"`
	Links      []map[string]string         `json:"links,omitempty"`
	Errors     []interface{}               `json:"errors,omitempty"`
	//{"errors":[{"param":"from","message":"From must be before to","value":"2014-07-20T09:00"},{"param":"to","message":"To must be in the format YYYY-MM-DDTHH:mm","value":"now"}]}
}

// https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-domain
func (c *Client) QueryDeliverabilityMetrics(extraPath string, parameters map[string]string) (*DeliverabilityMetricEventsWrapper, error) {
	return c.QueryDeliverabilityMetricsWithHeaders(extraPath, parameters, nil)
}

// https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-domain
func (c *Client) QueryDeliverabilityMetricsWithHeaders(extraPath string, parameters, headers map[string]string) (*DeliverabilityMetricEventsWrapper, error) {

	var finalUrl string
	path := fmt.Sprintf(deliverabilityMetricPathFormat, c.Config.ApiVersion)

	if extraPath != "" {
		path = fmt.Sprintf("%s/%s", path, extraPath)
	}

	//log.Printf("Path: %s", path)

	if parameters == nil || len(parameters) == 0 {
		finalUrl = fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	} else {
		params := URL.Values{}
		for k, v := range parameters {
			params.Add(k, v)
		}

		finalUrl = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, params.Encode())
	}

	return doMetricsRequest(c, finalUrl, headers)
}

func (c *Client) MetricEventAsString(e *DeliverabilityMetricItem) string {

	return fmt.Sprintf("domain: %s, [%v]", e.Domain, e)
}

func doMetricsRequest(c *Client, finalUrl string, headers map[string]string) (*DeliverabilityMetricEventsWrapper, error) {
	// Send off our request
	res, err := c.HttpGet(finalUrl, headers)
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
	var resMap DeliverabilityMetricEventsWrapper
	err = json.Unmarshal(bodyBytes, &resMap)

	if err != nil {
		return nil, err
	}

	return &resMap, err
}

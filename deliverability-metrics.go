package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"

	URL "net/url"
)

// https://www.sparkpost.com/api#/reference/message-events
var DeliverabilityMetricPathFormat = "/api/v%d/metrics/deliverability"

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

type DeliverabilityMetrics struct {
	Results    []DeliverabilityMetricItem `json:"results,omitempty"`
	TotalCount int                        `json:"total_count,omitempty"`
	Links      []map[string]string        `json:"links,omitempty"`
	Errors     []interface{}              `json:"errors,omitempty"`

	ExtraPath string            `json:"-"`
	Params    map[string]string `json:"-"`
	Context   context.Context   `json:"-"`
}

// https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-domain
func (c *Client) QueryDeliverabilityMetrics(dm *DeliverabilityMetrics) (*Response, error) {
	var finalUrl string
	path := fmt.Sprintf(DeliverabilityMetricPathFormat, c.Config.ApiVersion)

	if dm.ExtraPath != "" {
		path = fmt.Sprintf("%s/%s", path, dm.ExtraPath)
	}

	if dm.Params == nil || len(dm.Params) == 0 {
		finalUrl = fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	} else {
		params := URL.Values{}
		for k, v := range dm.Params {
			params.Add(k, v)
		}

		finalUrl = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, params.Encode())
	}

	return dm.doMetricsRequest(c, finalUrl)
}

func (dm *DeliverabilityMetrics) doMetricsRequest(c *Client, finalUrl string) (*Response, error) {
	// Send off our request
	res, err := c.HttpGet(dm.Context, finalUrl)
	if err != nil {
		return res, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return res, err
	}

	err = res.ParseResponse()
	if err != nil {
		return res, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return res, err
	}

	// Parse expected response structure
	err = json.Unmarshal(bodyBytes, dm)
	if err != nil {
		return res, err
	}

	return res, nil
}

package gosparkpost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

var MetricsPathFormat = "/api/v%d/metrics/deliverability"

type MetricItem struct {
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

type Metrics struct {
	Results    []MetricItem        `json:"results,omitempty"`
	TotalCount int                 `json:"total_count,omitempty"`
	Links      []map[string]string `json:"links,omitempty"`
	Errors     []interface{}       `json:"errors,omitempty"`

	ExtraPath string            `json:"-"`
	Params    map[string]string `json:"-"`
}

// https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-domain
func (c *Client) QueryMetrics(m *Metrics) (*Response, error) {
	return c.QueryMetricsContext(context.Background(), m)
}

func (c *Client) QueryMetricsContext(ctx context.Context, m *Metrics) (*Response, error) {
	var finalUrl string
	path := fmt.Sprintf(MetricsPathFormat, c.Config.ApiVersion)

	if m.ExtraPath != "" {
		path = fmt.Sprintf("%s/%s", path, m.ExtraPath)
	}

	if m.Params == nil || len(m.Params) == 0 {
		finalUrl = fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	} else {
		params := url.Values{}
		for k, v := range m.Params {
			params.Add(k, v)
		}

		finalUrl = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, params.Encode())
	}

	return m.doMetricsRequest(ctx, c, finalUrl)
}

func (m *Metrics) doMetricsRequest(ctx context.Context, c *Client, finalUrl string) (*Response, error) {
	// Send off our request
	res, err := c.HttpGet(ctx, finalUrl)
	if err != nil {
		return res, err
	}

	var body []byte
	// Assert that we got a JSON Content-Type back
	if body, err = res.AssertJson(); err != nil {
		return res, err
	}

	err = res.ParseResponse()
	if err != nil {
		return res, err
	}

	// Parse expected response structure
	err = json.Unmarshal(body, m)
	if err != nil {
		return res, errors.Wrap(err, "unmarshaling response")
	}

	return res, nil
}

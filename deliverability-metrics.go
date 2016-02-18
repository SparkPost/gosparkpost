package gosparkpost

import (
	"encoding/json"
	"fmt"
	"log"

	URL "net/url"
)

// https://www.sparkpost.com/api#/reference/message-events
var deliverabilityMetricPathFormat = "/api/v%d/metrics/deliverability"

type DeliverabilityMetricItem struct {
	count_injected                 int    `json:"count_injected,omitempty"`
	count_bounce                   int    `json:"count_bounce,omitempty"`
	count_rejected                 int    `json:"count_rejected,omitempty"`
	count_delivered                int    `json:"count_delivered,omitempty"`
	count_delivered_first          int    `json:"count_delivered_first,omitempty"`
	count_delivered_subsequent     int    `json:"count_delivered_subsequent,omitempty"`
	total_delivery_time_first      int    `json:"total_delivery_time_first,omitempty"`
	total_delivery_time_subsequent int    `json:"total_delivery_time_subsequent,omitempty"`
	total_msg_volume               int    `json:"total_msg_volume,omitempty"`
	count_policy_rejection         int    `json:"count_policy_rejection,omitempty"`
	count_generation_rejection     int    `json:"count_generation_rejection,omitempty"`
	count_generation_failed        int    `json:"count_generation_failed,omitempty"`
	count_inband_bounce            int    `json:"count_inband_bounce,omitempty"`
	count_outofband_bounce         int    `json:"count_outofband_bounce,omitempty"`
	count_soft_bounce              int    `json:"count_soft_bounce,omitempty"`
	count_hard_bounce              int    `json:"count_hard_bounce,omitempty"`
	count_block_bounce             int    `json:"count_block_bounce,omitempty"`
	count_admin_bounce             int    `json:"count_admin_bounce,omitempty"`
	count_undetermined_bounce      int    `json:"count_undetermined_bounce,omitempty"`
	count_delayed                  int    `json:"count_delayed,omitempty"`
	count_delayed_first            int    `json:"count_delayed_first,omitempty"`
	count_rendered                 int    `json:"count_rendered,omitempty"`
	count_unique_rendered          int    `json:"count_unique_rendered,omitempty"`
	count_unique_confirmed_opened  int    `json:"count_unique_confirmed_opened,omitempty"`
	count_clicked                  int    `json:"count_clicked,omitempty"`
	count_unique_clicked           int    `json:"count_unique_clicked,omitempty"`
	count_targeted                 int    `json:"count_targeted,omitempty"`
	count_sent                     int    `json:"count_sent,omitempty"`
	count_accepted                 int    `json:"count_accepted,omitempty"`
	count_spam_complaint           int    `json:"count_spam_complaint,omitempty"`
	domain                         string `json:"domain,omitempty"`
}

type DeliverabilityMetricEventsWrapper struct {
	Results    []*DeliverabilityMetricItem `json:"results,omitempty"`
	TotalCount int                         `json:"total_count,omitempty"`
	Links      []map[string]string         `json:"links,omitempty"`
	Errors     []interface{}               `json:"errors,omitempty"`
	//{"errors":[{"param":"from","message":"From must be before to","value":"2014-07-20T09:00"},{"param":"to","message":"To must be in the format YYYY-MM-DDTHH:mm","value":"now"}]}
}

// https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-domain
func (c *Client) QueryDeliverabilityMetrics(extraPath string, parameters map[string]string) (*DeliverabilityMetricEventsWrapper, error) {

	var finalUrl string
	path := fmt.Sprintf(deliverabilityMetricPathFormat, c.Config.ApiVersion)

	if extraPath != "" {
		path = fmt.Sprintf("%s/%s", path, extraPath)
	}

	log.Printf("Path: %s", path)

	if parameters == nil || len(parameters) == 0 {
		finalUrl = fmt.Sprintf("%s/%s", c.Config.BaseUrl, path)
	} else {
		params := URL.Values{}
		for k, v := range parameters {
			params.Add(k, v)
		}

		finalUrl = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, params.Encode())
	}

	return DoMetricsRequest(c, finalUrl)
}

func (c *Client) MetricEventAsString(e *DeliverabilityMetricItem) string {

	

	return fmt.Sprintf("domain: %s, count_injected: %d", e.domain, e.count_injected)
}

func DoMetricsRequest(c *Client, finalUrl string) (*DeliverabilityMetricEventsWrapper, error) {
	// Send off our request
	res, err := c.HttpGet(finalUrl)
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
	var resMap DeliverabilityMetricEventsWrapper //map[string]interface{}
	err = json.Unmarshal(bodyBytes, &resMap)
	log.Printf("Decoded : %s", resMap)
	if err != nil {
		return nil, err
	}

	return &resMap, err
}


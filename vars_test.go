package gosparkpost_test

import (
	"encoding/json"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/SparkPost/gosparkpost/events"
)

func init() {
	json.Unmarshal([]byte(`"2017-04-26T21:37:32.000+00:00"`), &msgEventsPage1.out.Events[0].(*events.Click).Timestamp)
	json.Unmarshal([]byte(`"2017-04-26T21:37:17.000+00:00"`), &msgEventsPage2.out.Events[0].(*events.Open).Timestamp)
	json.Unmarshal([]byte(`"2017-04-26T21:37:17.000+00:00"`), &msgEventsPage3.out.Events[0].(*events.Open).Timestamp)
}

var click1 = &events.Click{
	EventCommon:     events.EventCommon{Type: "click"},
	CampaignID:      "",
	CustomerID:      "42",
	DeliveryMethod:  "esmtp",
	GeoIP:           &events.GeoIP{Country: "US", Region: "NY", City: "Bronx", Latitude: 40.8499, Longitude: -73.8769},
	IPAddress:       "66.102.8.2",
	MessageID:       "0001441ead58afdeb21d",
	Metadata:        map[string]interface{}{},
	Tags:            []string{},
	Recipient:       "developers@sparkpost.com",
	RecipientType:   "",
	TargetLinkName:  "",
	TargetLinkURL:   "https://sparkpost.com",
	TemplateID:      "best-template-ever",
	TemplateVersion: "10",
	TransmissionID:  "30494147183576458",
	UserAgent:       "lynx",
}

var open1 = &events.Open{
	EventCommon:     events.EventCommon{Type: "open"},
	CampaignID:      "",
	CustomerID:      "42",
	DeliveryMethod:  "esmtp",
	GeoIP:           &events.GeoIP{Country: "US", Region: "NY", City: "Bronx", Latitude: 40.8499, Longitude: -73.8769},
	IPAddress:       "66.102.8.28",
	MessageID:       "0001441ead58afdeb21d",
	Metadata:        map[string]interface{}{},
	Tags:            []string{},
	Recipient:       "developers@sparkpost.com",
	RecipientType:   "",
	TemplateID:      "best-template-ever",
	TemplateVersion: "10",
	TransmissionID:  "30494147183576458",
	UserAgent:       "Mozilla/5.0 (Windows NT 5.1; rv:11.0) Gecko Firefox/11.0 (via ggpht.com GoogleImageProxy)",
}

var open2 = &events.Open{
	EventCommon:     events.EventCommon{Type: "open"},
	CampaignID:      "",
	CustomerID:      "42",
	DeliveryMethod:  "esmtp",
	GeoIP:           &events.GeoIP{Country: "US", Region: "NY", City: "Bronx", Latitude: 40.8499, Longitude: -73.8769},
	IPAddress:       "66.102.8.2",
	MessageID:       "00042a25ad58fc0145e1",
	Metadata:        map[string]interface{}{},
	Tags:            []string{},
	Recipient:       "sales@sparkpost.com",
	RecipientType:   "",
	TemplateID:      "best-template-ever",
	TemplateVersion: "10",
	TransmissionID:  "84537445295705713",
	UserAgent:       "Mozilla/5.0 (Windows NT 5.1; rv:11.0) Gecko Firefox/11.0 (via ggpht.com GoogleImageProxy)",
}

var msgEventsPageAll = EventsPageResult{nil, 200, ``, &sp.EventsPage{
	TotalCount: 3,
	Events:     events.Events{click1, open1, open2},
}}

var msgEventsPage1 = EventsPageResult{nil, 200, ``, &sp.EventsPage{
	TotalCount: 3,
	Events:     events.Events{click1},
	NextPage:   "/api/v1/message-events?page=2&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:09&timezone=UTC",
	LastPage:   "/api/v1/message-events?page=3&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:09&timezone=UTC",
}}

var msgEventsPage2 = EventsPageResult{nil, 200, ``, &sp.EventsPage{
	TotalCount: 3,
	Events:     events.Events{open1},
	NextPage:   "/api/v1/message-events?page=3&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:12&timezone=UTC",
	PrevPage:   "/api/v1/message-events?page=1&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:12&timezone=UTC",
	FirstPage:  "/api/v1/message-events?page=1&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:12&timezone=UTC",
	LastPage:   "/api/v1/message-events?page=3&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:12&timezone=UTC",
}}

var msgEventsPage3 = EventsPageResult{nil, 200, ``, &sp.EventsPage{
	TotalCount: 3,
	Events:     events.Events{open2},
	PrevPage:   "/api/v1/message-events?page=2&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:13&timezone=UTC",
	FirstPage:  "/api/v1/message-events?page=1&per_page=1&from=2017-04-26T00:00&to=2017-04-27T01:13&timezone=UTC",
}}

package gosparkpost_test

import (
	"fmt"
	"net/http"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
)

// The "links" section is snipped for brevity
var delMetricsBaseNoArgs string = `{
    "errors": [
        {
            "message": "from is required",
            "param": "from"
        },
        {
            "message": "from must be in the format YYYY-MM-DDTHH:MM",
            "param": "from"
        },
        {
            "message": "from must be before to",
            "param": "from"
        }
    ],
    "links": [
        {
            "href": "/api/v1/metrics/deliverability",
            "method": "GET",
            "rel": "deliverability"
        },
        {
            "href": "/api/v1/metrics/deliverability/watched-domain",
            "method": "GET",
            "rel": "watched-domain"
        }
    ]
}`

func TestMetrics_Get_noArgsError(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	path := fmt.Sprintf(sp.MetricsPathFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(delMetricsBaseNoArgs))
	})

	m := &sp.Metrics{}
	res, err := testClient.QueryMetrics(m)
	if err != nil {
		testFailVerbose(t, res, "Metrics GET returned error: %+v", err)
	}

	if len(m.Errors) != 3 {
		testFailVerbose(t, res, "Expected 3 errors, got %d", len(m.Errors))
	}
}

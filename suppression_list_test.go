package gosparkpost

import (
	"fmt"
	"net/http"
	"testing"
)

var suppressionNotFound string = `{
	"errors": [
		{
			"message": "Recipient could not be found"
		}
	]
}`

// Test parsing of "not found" case
func TestSuppression_Get_notFound(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	path := fmt.Sprintf(suppressionListsPathFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusNotFound) // 404
		w.Write([]byte(suppressionNotFound))
	})

	// hit our local handler
	s, res, err := testClient.SuppressionList()
	if err != nil {
		testFailVerbose(t, res, "SuppressionList GET returned error: %v", err)
	}

	// basic content test
	if s.Results != nil {
		testFailVerbose(t, res, "SuppressionList GET returned non-nil Results (error expected)")
	} else if len(s.Results) != 0 {
		testFailVerbose(t, res, "SuppressionList GET returned %d results, expected %d", len(s.Results), 0)
	} else if len(res.Errors) != 1 {
		testFailVerbose(t, res, "SuppressionList GET returned %d errors, expected %d", len(res.Errors), 1)
	} else if res.Errors[0].Message != "Recipient could not be found" {
		testFailVerbose(t, res, "SuppressionList GET Unmarshal error; saw [%v] expected [%v]",
			res.Errors[0].Message, "Recipient could not be found")
	}
}

var combinedSuppressionList string = `{
  "results": [
    {
      "recipient": "rcpt_1@example.com",
      "transactional": true,
      "non_transactional": true,
      "source": "Manually Added",
      "description": "User requested to not receive any non-transactional emails.",
      "created": "2016-01-01T12:00:00+00:00",
      "updated": "2016-01-01T12:00:00+00:00"
    }
  ]
}`

// Test parsing of combined suppression list results
func TestSuppression_Get_combinedList(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	path := fmt.Sprintf(suppressionListsPathFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write([]byte(combinedSuppressionList))
	})

	// hit our local handler
	s, res, err := testClient.SuppressionList()
	if err != nil {
		t.Errorf("SuppressionList GET returned error: %v", err)
		for _, e := range res.Verbose {
			t.Error(e)
		}
		return
	}

	// basic content test
	if s.Results == nil {
		t.Error("SuppressionList GET returned nil Results")
	} else if len(s.Results) != 1 {
		t.Errorf("SuppressionList GET returned %d results, expected %d", len(s.Results), 1)
	} else if s.Results[0].Recipient != "rcpt_1@example.com" {
		t.Errorf("SuppressionList GET Unmarshal error; saw [%v] expected [rcpt_1@example.com]", s.Results[0].Recipient)
	}
}

var separateSuppressionList string = `{
  "results": [
    {
      "recipient": "rcpt_1@example.com",
      "non_transactional": true,
      "source": "Manually Added",
      "type": "non_transactional",
      "description": "User requested to not receive any non-transactional emails.",
      "created": "2016-01-01T12:00:00+00:00",
      "updated": "2016-01-01T12:00:00+00:00"
    },
    {
      "recipient": "rcpt_1@example.com",
      "transactional": true,
      "source": "Bounce Rule",
      "type": "transactional",
      "description": "550: 550-5.1.1 Invalid Recipient",
      "created": "2015-10-15T12:00:00+00:00",
      "updated": "2015-10-15T12:00:00+00:00"
    }
  ],
  "links": [],
  "total_count": 2
}`

// Test parsing of separate suppression list results
func TestSuppression_Get_separateList(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	path := fmt.Sprintf(suppressionListsPathFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write([]byte(separateSuppressionList))
	})

	// hit our local handler
	s, res, err := testClient.SuppressionList()
	if err != nil {
		t.Errorf("SuppressionList GET returned error: %v", err)
		for _, e := range res.Verbose {
			t.Error(e)
		}
		return
	}

	// basic content test
	if s.Results == nil {
		t.Error("SuppressionList GET returned nil Results")
	} else if len(s.Results) != 2 {
		t.Errorf("SuppressionList GET returned %d results, expected %d", len(s.Results), 2)
	} else if s.Results[0].Recipient != "rcpt_1@example.com" {
		t.Errorf("SuppressionList GET Unmarshal error; saw [%v] expected [rcpt_1@example.com]", s.Results[0].Recipient)
	}
}

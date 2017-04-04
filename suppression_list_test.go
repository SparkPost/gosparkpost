package gosparkpost_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
)

var suppressionNotFound = loadTestFile("test_data/suppression_not_found_error.json")

// Test parsing of "not found" case
func TestSuppression_Get_notFound(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	path := fmt.Sprintf(sp.SuppressionListsPathFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusNotFound) // 404
		w.Write([]byte(suppressionNotFound))
	})

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}

	res, err := testClient.SuppressionList(suppressionPage)
	if err != nil {
		testFailVerbose(t, res, "SuppressionList GET returned error: %v", err)
	}

	// basic content test
	if suppressionPage.Results != nil {
		testFailVerbose(t, res, "SuppressionList GET returned non-nil Results (error expected)")
	} else if len(suppressionPage.Results) != 0 {
		testFailVerbose(t, res, "SuppressionList GET returned %d results, expected %d", len(suppressionPage.Results), 0)
	} else if len(suppressionPage.Errors) != 1 {
		testFailVerbose(t, res, "SuppressionList GET returned %d errors, expected %d", len(suppressionPage.Errors), 1)
	} else if suppressionPage.Errors[0].Message != "Recipient could not be found" {
		testFailVerbose(t, res, "SuppressionList GET Unmarshal error; saw [%v] expected [%v]",
			res.Errors[0].Message, "Recipient could not be found")
	}
}

var combinedSuppressionList = loadTestFile("test_data/suppression_combined.json")

// Test parsing of combined suppression list results
func TestSuppression_Get_combinedList(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	path := fmt.Sprintf(sp.SuppressionListsPathFormat, testClient.Config.ApiVersion)
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

var separateSuppressionList = loadTestFile("test_data/suppression_seperate_lists.json")

// Test parsing of separate suppression list results
func TestSuppression_Get_separateList(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	path := fmt.Sprintf(sp.SuppressionListsPathFormat, testClient.Config.ApiVersion)
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

var suppressionListCursor = loadTestFile("test_data/suppression_cursor.json")

// Test parsing of separate suppression list results
func TestSuppression_links(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	path := fmt.Sprintf(sp.SuppressionListsPathFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write([]byte(suppressionListCursor))
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
	} else if s.TotalCount != 44 {
		t.Errorf("SuppressionList GET returned %d results, expected %d", s.TotalCount, 44)
	} else if len(s.Links) != 2 {
		t.Errorf("SuppressionList GET returned %d results, expected %d", len(s.Links), 2)
	} else if s.Links[0].Href != "The_HREF_Value1" {
		t.Error("SuppressionList GET returned invalid link[0].Href")
	} else if s.Links[1].Href != "The_HREF_Value2" {
		t.Error("SuppressionList GET returned invalid link[1].Href")
	} else if s.Links[0].Rel != "first" {
		t.Error("SuppressionList GET returned invalid s.Links[0].Rel")
	} else if s.Links[1].Rel != "next" {
		t.Error("SuppressionList GET returned invalid s.Links[1].Rel")
	}

}

func loadTestFile(fileToLoad string) string {
	b, err := ioutil.ReadFile(fileToLoad)
	if err != nil {
		fmt.Print(err)
	}

	return string(b)
}

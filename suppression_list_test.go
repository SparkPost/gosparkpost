package gosparkpost_test

import (
	"fmt"
	"net/http"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
)

// Test parsing of "not found" case
func TestSuppression_Get_notFound(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile("test_data/suppression_not_found_error.json")
	mockRestBuilderFormat(t, sp.SuppressionListsPathFormat, mockResponse)

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

func TestSuppression_Error_Bad_Path(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile("test_data/suppression_not_found_error.json")
	mockRestBuilderFormat(t, "/bad/path", mockResponse)

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}
	res, err := testClient.SuppressionList(suppressionPage)
	if err.Error() != "Expected json, got [text/plain] with code 404" {
		testFailVerbose(t, res, "SuppressionList GET returned error: %v", err)
	} else if res.HTTP.StatusCode != 404 {
		testFailVerbose(t, res, "Expected a 404 error: %v", res)
	}

}

func TestSuppression_Error_Bad_JSON(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	mockRestBuilderFormat(t, sp.SuppressionListsPathFormat, "ThisIsBadJSON")

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}

	// Bad JSON should generate an Error
	res, err := testClient.SuppressionList(suppressionPage)
	if err == nil {
		testFailVerbose(t, res, "Expected an error due to bad JSON: %v", err)
	}

}

func TestSuppression_Error_Wrong_JSON(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	mockRestBuilderFormat(t, sp.SuppressionListsPathFormat, "{\"errors\":\"\"")

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}

	// Bad JSON should generate an Error
	res, err := testClient.SuppressionList(suppressionPage)
	if err == nil {
		testFailVerbose(t, res, "Expected an error due to bad JSON: %v", err)
	}

}

// Test parsing of combined suppression list results
func TestSuppression_Get_combinedList(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile("test_data/suppression_combined.json")
	mockRestBuilderFormat(t, sp.SuppressionListsPathFormat, mockResponse)

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}
	res, err := testClient.SuppressionList(suppressionPage)
	if err != nil {
		t.Errorf("SuppressionList GET returned error: %v", err)
		for _, e := range res.Verbose {
			t.Error(e)
		}
		return
	}

	// basic content test
	if suppressionPage.Results == nil {
		t.Error("SuppressionList GET returned nil Results")
	} else if len(suppressionPage.Results) != 1 {
		t.Errorf("SuppressionList GET returned %d results, expected %d", len(suppressionPage.Results), 1)
	} else if suppressionPage.Results[0].Recipient != "rcpt_1@example.com" {
		t.Errorf("SuppressionList GET Unmarshal error; saw [%v] expected [rcpt_1@example.com]", suppressionPage.Results[0].Recipient)
	}
}

// Test parsing of separate suppression list results
func TestSuppression_Get_separateList(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile("test_data/suppression_seperate_lists.json")
	mockRestBuilderFormat(t, sp.SuppressionListsPathFormat, mockResponse)

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}
	res, err := testClient.SuppressionList(suppressionPage)
	if err != nil {
		t.Errorf("SuppressionList GET returned error: %v", err)
		for _, e := range res.Verbose {
			t.Error(e)
		}
		return
	}

	// basic content test
	if suppressionPage.Results == nil {
		t.Error("SuppressionList GET returned nil Results")
	} else if len(suppressionPage.Results) != 2 {
		t.Errorf("SuppressionList GET returned %d results, expected %d", len(suppressionPage.Results), 2)
	} else if suppressionPage.Results[0].Recipient != "rcpt_1@example.com" {
		t.Errorf("SuppressionList GET Unmarshal error; saw [%v] expected [rcpt_1@example.com]", suppressionPage.Results[0].Recipient)
	}
}

// Tests that links are generally parsed properly
func TestSuppression_links(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile("test_data/suppression_cursor.json")
	mockRestBuilderFormat(t, sp.SuppressionListsPathFormat, mockResponse)

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}
	res, err := testClient.SuppressionList(suppressionPage)
	if err != nil {
		t.Errorf("SuppressionList GET returned error: %v", err)
		for _, e := range res.Verbose {
			t.Error(e)
		}
		return
	}

	// basic content test
	if suppressionPage.Results == nil {
		t.Error("SuppressionList GET returned nil Results")
	} else if suppressionPage.TotalCount != 44 {
		t.Errorf("SuppressionList GET returned %d results, expected %d", suppressionPage.TotalCount, 44)
	} else if len(suppressionPage.Links) != 4 {
		t.Errorf("SuppressionList GET returned %d results, expected %d", len(suppressionPage.Links), 2)
	} else if suppressionPage.Links[0].Href != "The_HREF_first" {
		t.Error("SuppressionList GET returned invalid link[0].Href")
	} else if suppressionPage.Links[1].Href != "The_HREF_next" {
		t.Error("SuppressionList GET returned invalid link[1].Href")
	} else if suppressionPage.Links[0].Rel != "first" {
		t.Error("SuppressionList GET returned invalid s.Links[0].Rel")
	} else if suppressionPage.Links[1].Rel != "next" {
		t.Error("SuppressionList GET returned invalid s.Links[1].Rel")
	}

	// Check convenience links
	if suppressionPage.FirstPage != "The_HREF_first" {
		t.Errorf("Unexpected FirstPage value: %s", suppressionPage.FirstPage)
	} else if suppressionPage.LastPage != "The_HREF_last" {
		t.Errorf("Unexpected LastPage value: %s", suppressionPage.LastPage)
	} else if suppressionPage.PrevPage != "The_HREF_previous" {
		t.Errorf("Unexpected PrevPage value: %s", suppressionPage.PrevPage)
	} else if suppressionPage.NextPage != "The_HREF_next" {
		t.Errorf("Unexpected NextPage value: %s", suppressionPage.NextPage)
	}

}

func TestSuppression_Empty_NextPage(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile("test_data/suppression_single_page.json")
	mockRestBuilderFormat(t, sp.SuppressionListsPathFormat, mockResponse)

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}
	res, err := testClient.SuppressionList(suppressionPage)
	if err != nil {
		t.Errorf("SuppressionList GET returned error: %v", err)
		for _, e := range res.Verbose {
			t.Error(e)
		}
		return
	}

	nextResponse, res, err := suppressionPage.Next()

	if nextResponse != nil {
		t.Errorf("nextResponse should be nil but was: %v", nextResponse)
	} else if res != nil {
		t.Errorf("Response should be nil but was: %v", res)
	} else if err != nil {
		t.Errorf("Error should be nil but was: %v", err)
	}
}

//
func TestSuppression_NextPage(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile("test_data/suppression_page1.json")
	mockRestBuilderFormat(t, sp.SuppressionListsPathFormat, mockResponse)

	mockResponse = loadTestFile("test_data/suppression_page2.json")
	mockRestBuilder(t, "/test_data/suppression_page2.json", mockResponse)

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}
	res, err := testClient.SuppressionList(suppressionPage)
	if err != nil {
		t.Errorf("SuppressionList GET returned error: %v", err)
		for _, e := range res.Verbose {
			t.Error(e)
		}
		return
	}

	if suppressionPage.NextPage != "/test_data/suppression_page2.json" {
		t.Errorf("Unexpected NextPage value: %s", suppressionPage.NextPage)
	}

	nextResponse, res, err := suppressionPage.Next()

	if nextResponse.NextPage != "/test_data/suppression_pageLast.json" {
		t.Errorf("Unexpected NextPage value: %s", nextResponse.NextPage)
	}
}

func mockRestBuilderFormat(t *testing.T, pathFormat string, mockResponse string) {
	path := fmt.Sprintf(pathFormat, testClient.Config.ApiVersion)
	mockRestBuilder(t, path, mockResponse)
}

func mockRestBuilder(t *testing.T, path string, mockResponse string) {
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write([]byte(mockResponse))
	})
}

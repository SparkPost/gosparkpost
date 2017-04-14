package gosparkpost_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"encoding/json"

	sp "github.com/SparkPost/gosparkpost"
)

func TestUnmarshal_SupressionEvent(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	var suppressionEventString = loadTestFile(t, "test/json/suppression_entry_simple.json")

	suppressionEntry := &sp.SuppressionEntry{}
	err := json.Unmarshal([]byte(suppressionEventString), suppressionEntry)
	if err != nil {
		testFailVerbose(t, nil, "Unmarshal SuppressionEntry returned error: %v", err)
	}

	verifySuppressionEnty(t, suppressionEntry)
}

// Test parsing of "not found" case
func TestSuppression_Get_notFound(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile(t, "test/json/suppression_not_found_error.json")
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, mockResponse)

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

func TestSuppression_Retrieve(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile(t, "test/json/suppression_retrieve.json")
	status := http.StatusOK
	email := "john.doe@domain.com"
	mockRestResponseBuilderFormat(t, "GET", status, sp.SuppressionListsPathFormat+"/"+email, mockResponse)

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}
	res, err := testClient.SuppressionRetrieve(email, suppressionPage)
	if err != nil {
		testFailVerbose(t, res, "SuppressionList retrieve returned error: %v", err)
	} else if res == nil {
		testFailVerbose(t, res, "SuppressionList retrieve expected an HTTP response")
	}

	if len(suppressionPage.Results) != 1 {
		testFailVerbose(t, res, "SuppressionList retrieve expected 1 result: %v", suppressionPage)
	} else if suppressionPage.TotalCount != 1 {
		testFailVerbose(t, res, "SuppressionList retrieve expected 1 result: %v", suppressionPage)
	}

	verifySuppressionEnty(t, suppressionPage.Results[0])

}

func TestSuppression_Error_Bad_Path(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile(t, "test/json/suppression_not_found_error.json")
	mockRestBuilderFormat(t, "GET", "/bad/path", mockResponse)

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
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, "ThisIsBadJSON")

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
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, "{\"errors\":\"\"")

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
	var mockResponse = loadTestFile(t, "test/json/suppression_combined.json")
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, mockResponse)

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
	var mockResponse = loadTestFile(t, "test/json/suppression_seperate_lists.json")
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, mockResponse)

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
	var mockResponse = loadTestFile(t, "test/json/suppression_cursor.json")
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, mockResponse)

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
	var mockResponse = loadTestFile(t, "test/json/suppression_single_page.json")
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, mockResponse)

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
	var mockResponse = loadTestFile(t, "test/json/suppression_page1.json")
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, mockResponse)

	mockResponse = loadTestFile(t, "test/json/suppression_page2.json")
	mockRestBuilder(t, "GET", "/test/json/suppression_page2.json", mockResponse)

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

	if suppressionPage.NextPage != "/test/json/suppression_page2.json" {
		t.Errorf("Unexpected NextPage value: %s", suppressionPage.NextPage)
	}

	nextResponse, res, err := suppressionPage.Next()

	if nextResponse.NextPage != "/test/json/suppression_pageLast.json" {
		t.Errorf("Unexpected NextPage value: %s", nextResponse.NextPage)
	}
}

// Test parsing of combined suppression list results
func TestSuppression_Search_combinedList(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile(t, "test/json/suppression_combined.json")
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, mockResponse)

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}
	res, err := testClient.SuppressionSearch(suppressionPage)
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

// Test parsing of combined suppression list results
func TestSuppression_Search_params(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	var mockResponse = loadTestFile(t, "test/json/suppression_combined.json")
	mockRestBuilderFormat(t, "GET", sp.SuppressionListsPathFormat, mockResponse)

	// hit our local handler
	suppressionPage := &sp.SuppressionPage{}
	parameters := map[string]string{
		"from": "1970-01-01T00:00",
	}
	suppressionPage.Params = parameters

	res, err := testClient.SuppressionSearch(suppressionPage)
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

func TestClient_SuppressionUpsert_nil_entry(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	response, err := testClient.SuppressionUpsert(nil)
	if response != nil {
		t.Errorf("Expected nil response object but got: %v", response)
	} else if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestClient_SuppressionUpsert_bad_json(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	var mockResponse = "{bad json}"
	mockRestBuilderFormat(t, "PUT", sp.SuppressionListsPathFormat, mockResponse)

	entry := sp.WritableSuppressionEntry{
		Recipient:   "john.doe@domain.com",
		Description: "entry description",
		Type:        "non_transactional",
	}

	entries := []sp.WritableSuppressionEntry{
		entry,
	}

	response, err := testClient.SuppressionUpsert(entries)
	if response == nil {
		t.Errorf("Expected a response")
	} else if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestClient_SuppressionUpsert_1_entry(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	var expectedRequest = loadTestFile(t, "test/json/suppression_entry_simple_request.json")
	var mockResponse = "{}"
	mockRestRequestResponseBuilderFormat(t, "PUT", http.StatusOK, sp.SuppressionListsPathFormat, expectedRequest, mockResponse)

	entry := sp.WritableSuppressionEntry{
		Recipient:   "john.doe@domain.com",
		Description: "entry description",
		Type:        "non_transactional",
	}

	entries := []sp.WritableSuppressionEntry{
		entry,
	}

	response, err := testClient.SuppressionUpsert(entries)
	if response == nil {
		t.Errorf("Expected a response")
	} else if err != nil {
		t.Errorf("Did not expect an error: %v", err)
	}
}

func TestClient_SuppressionUpsert_error_response(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	var mockResponse = loadTestFile(t, "test/json/suppression_not_found_error.json")
	status := http.StatusBadRequest
	mockRestResponseBuilderFormat(t, "PUT", status, sp.SuppressionListsPathFormat, mockResponse)

	entry := sp.WritableSuppressionEntry{
		Recipient:   "john.doe@domain.com",
		Description: "entry description",
		Type:        "non_transactional",
	}

	entries := []sp.WritableSuppressionEntry{
		entry,
	}

	response, err := testClient.SuppressionUpsert(entries)
	if response == nil {
		t.Errorf("Expected a response")
	} else if err == nil {
		t.Errorf("Expected an error with the HTTP status code")
	}

	if response.HTTP.StatusCode != status {
		testFailVerbose(t, response, "Expected HTTP status code %d but got %d", status, response.HTTP.StatusCode)
	} else if len(response.Errors) != 1 {
		testFailVerbose(t, response, "SuppressionUpsert PUT returned %d errors, expected %d", len(response.Errors), 1)
	} else if response.Errors[0].Message != "Recipient could not be found" {
		testFailVerbose(t, response, "SuppressionUpsert PUT Unmarshal error; saw [%v] expected [%v]",
			response.Errors[0].Message, "Recipient could not be found")
	}
}

func TestClient_Suppression_Delete_nil_email(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	status := http.StatusNotFound
	mockRestResponseBuilderFormat(t, "DELETE", status, sp.SuppressionListsPathFormat+"/", "")

	_, err := testClient.SuppressionDelete("")
	if err == nil {
		t.Errorf("Expected an error indicating an email address is required")
	}
}

func TestClient_Suppression_Delete(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	email := "test@test.com"
	status := http.StatusNoContent
	mockRestResponseBuilderFormat(t, "DELETE", status, sp.SuppressionListsPathFormat+"/"+email, "")

	response, err := testClient.SuppressionDelete(email)
	if err != nil {
		t.Errorf("Did not expect an error")
	}

	if response.HTTP.StatusCode != status {
		testFailVerbose(t, response, "Expected HTTP status code %d but got %d", status, response.HTTP.StatusCode)
	}
}

func TestClient_Suppression_Delete_Errors(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	email := "test@test.com"
	status := http.StatusBadRequest
	var mockResponse = loadTestFile(t, "test/json/suppression_not_found_error.json")
	mockRestResponseBuilderFormat(t, "DELETE", status, sp.SuppressionListsPathFormat+"/"+email, mockResponse)

	response, err := testClient.SuppressionDelete(email)
	if err == nil {
		t.Errorf("Expected an error")
	}

	if response.HTTP.StatusCode != status {
		testFailVerbose(t, response, "Expected HTTP status code %d but got %d", status, response.HTTP.StatusCode)
	} else if len(response.Errors) != 1 {
		testFailVerbose(t, response, "SuppressionDelete DELETE returned %d errors, expected %d", len(response.Errors), 1)
	} else if response.Errors[0].Message != "Recipient could not be found" {
		testFailVerbose(t, response, "SuppressionDelete DELETE Unmarshal error; saw [%v] expected [%v]",
			response.Errors[0].Message, "Recipient could not be found")
	}
}

/////////////////////
// Internal Helpers
/////////////////////

func verifySuppressionEnty(t *testing.T, suppressionEntry *sp.SuppressionEntry) {
	if suppressionEntry.Recipient != "john.doe@domain.com" {
		testFailVerbose(t, nil, "Unexpected Recipient: %s", suppressionEntry.Recipient)
	} else if suppressionEntry.Description != "entry description" {
		testFailVerbose(t, nil, "Unexpected Description: %s", suppressionEntry.Description)
	} else if suppressionEntry.Source != "manually created" {
		testFailVerbose(t, nil, "Unexpected Source: %s", suppressionEntry.Source)
	} else if suppressionEntry.Type != "non_transactional" {
		testFailVerbose(t, nil, "Unexpected Type: %s", suppressionEntry.Type)
	} else if suppressionEntry.Created != "2016-05-02T16:29:56+00:00" {
		testFailVerbose(t, nil, "Unexpected Created: %s", suppressionEntry.Created)
	} else if suppressionEntry.Updated != "2016-05-02T17:20:50+00:00" {
		testFailVerbose(t, nil, "Unexpected Updated: %s", suppressionEntry.Updated)
	} else if suppressionEntry.NonTransactional != true {
		testFailVerbose(t, nil, "Unexpected NonTransactional value")
	}
}

func mockRestBuilderFormat(t *testing.T, method string, pathFormat string, mockResponse string) {
	mockRestResponseBuilderFormat(t, method, http.StatusOK, pathFormat, mockResponse)
}

func mockRestBuilder(t *testing.T, method string, path string, mockResponse string) {
	mockRestResponseBuilder(t, method, http.StatusOK, path, mockResponse)
}

func mockRestResponseBuilderFormat(t *testing.T, method string, status int, pathFormat string, mockResponse string) {
	path := fmt.Sprintf(pathFormat, testClient.Config.ApiVersion)
	mockRestResponseBuilder(t, method, status, path, mockResponse)
}

func mockRestResponseBuilder(t *testing.T, method string, status int, path string, mockResponse string) {
	mockRestRequestResponseBuilder(t, method, status, path, "", mockResponse)
}

func mockRestRequestResponseBuilderFormat(t *testing.T, method string, status int, pathFormat string, expectedBody string, mockResponse string) {
	path := fmt.Sprintf(pathFormat, testClient.Config.ApiVersion)
	mockRestRequestResponseBuilder(t, method, status, path, expectedBody, mockResponse)
}

func mockRestRequestResponseBuilder(t *testing.T, method string, status int, path string, expectedBody string, mockResponse string) {
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if expectedBody != "" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				testFailVerbose(t, nil, "error: %v", err)
			}

			ok, err := AreEqualJSON(expectedBody, string(body[:]))
			if err != nil {
				testFailVerbose(t, nil, "error: %v", err)
			}

			if !ok {
				testFailVerbose(t, nil, "Request did not match expected. \nExpected: \n%s\n\nActual:\n%s\n\n", err)
			}
		}

		testMethod(t, r, method)
		if mockResponse != "" {
			w.Header().Set("Content-Type", "application/json; charset=utf8")
		}
		w.WriteHeader(status)
		if mockResponse != "" {
			w.Write([]byte(mockResponse))
		}
	})
}

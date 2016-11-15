package gosparkpost

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const eventDocumentationFile = "test/event-docs.json"

var eventDocumentationBytes []byte

func init() {
	eventDocumentationBytes, err := ioutil.ReadFile(eventDocumentationFile)
	if err != nil {
		panic(err)
	}
}

func TestEventDocs_Get_parse(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	// set up the response handler
	path := fmt.Sprintf(eventDocumentationFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write(eventDocumentationBytes)
	})

	// hit our local handler
	w, res, err := testClient.EventDocumentation()
	if err != nil {
		t.Errorf("EventDocumentation GET returned error: %v", err)
		for _, e := range res.Verbose {
			t.Error(e)
		}
		return
	}

	// basic content test
	if w.Results == nil {
		t.Error("EventDocumentation GET returned nil Results")
	}
}

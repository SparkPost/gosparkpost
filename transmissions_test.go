package gosparkpost_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
)

var transmissionSuccess string = `{
  "results": {
    "total_rejected_recipients": 0,
    "total_accepted_recipients": 1,
    "id": "11111111111111111"
  }
}`

func TestTransmissions_Post_Success(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	path := fmt.Sprintf(sp.TransmissionsPathFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(transmissionSuccess))
	})

	testClient.Config.ApiKey = "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
	// set some headers on the client
	testClient.Headers.Add("X-Foo", "foo")
	testClient.Headers.Add("X-Foo", "baz")
	testClient.Headers.Add("X-Bar", "bar")
	testClient.Headers.Add("X-Baz", "baz")
	testClient.Headers.Del("X-Baz")
	// override one of the headers using a context
	header := http.Header{}
	header.Add("X-Foo", "bar")
	ctx := context.WithValue(context.Background(), "http.Header", header)
	tx := &sp.Transmission{
		CampaignID: "Post_Success",
		ReturnPath: "returnpath@example.com",
		Recipients: []string{"recipient1@example.com"},
		Content: sp.Content{
			Subject: "this is a test message",
			HTML:    "<h1>TEST</h1>",
			From:    map[string]string{"name": "test", "email": "from@example.com"},
		},
		Metadata: map[string]interface{}{"shoe_size": 9},
	}
	// send using the client and the context
	id, res, err := testClient.SendContext(ctx, tx)
	if err != nil {
		testFailVerbose(t, res, "Transmission POST returned error: %v", err)
	}

	if id != "11111111111111111" {
		testFailVerbose(t, res, "Unexpected value for id! (expected: 11111111111111111, saw: %s)", id)
	}

	var reqDump string
	var ok bool
	if reqDump, ok = res.Verbose["http_requestdump"]; !ok {
		testFailVerbose(t, res, "HTTP Request unavailable")
	}

	if !strings.Contains(reqDump, "X-Foo: bar") {
		testFailVerbose(t, res, "Header set on Client not sent")
	}
	if !strings.Contains(reqDump, "X-Bar: bar") {
		testFailVerbose(t, res, "Header set on Client not sent")
	}
	if strings.Contains(reqDump, "X-Baz: baz") {
		testFailVerbose(t, res, "Header set on Client should not have been sent")
	}
}

func TestTransmissions_Delete_Headers(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	path := fmt.Sprintf(sp.TransmissionsPathFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path+"/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})

	testClient.Config.Username = "testuser"
	testClient.Config.Password = "testpass"

	header := http.Header{}
	header.Add("X-Foo", "bar")
	ctx := context.WithValue(context.Background(), "http.Header", header)
	tx := &sp.Transmission{ID: "42"}
	res, err := testClient.TransmissionDeleteContext(ctx, tx)
	if err != nil {
		testFailVerbose(t, res, "Transmission DELETE failed")
	}

	var reqDump string
	var ok bool
	if reqDump, ok = res.Verbose["http_requestdump"]; !ok {
		testFailVerbose(t, res, "HTTP Request unavailable")
	}

	if !strings.Contains(reqDump, "X-Foo: bar") {
		testFailVerbose(t, res, "Header set on Transmission not sent")
	}
}

func TestTransmissions_ByID_Success(t *testing.T) {
	testSetup(t)
	defer testTeardown()

	tx := &sp.Transmission{
		CampaignID: "Post_Success",
		ReturnPath: "returnpath@example.com",
		Recipients: []string{"recipient1@example.com"},
		Content: sp.Content{
			Subject: "this is a test message",
			HTML:    "<h1>TEST</h1>",
			From:    map[string]string{"name": "test", "email": "from@example.com"},
		},
		Metadata: map[string]interface{}{"shoe_size": 9},
	}
	txBody := map[string]map[string]*sp.Transmission{"results": {"transmission": tx}}
	txBytes, err := json.Marshal(txBody)
	if err != nil {
		t.Error(err)
	}

	path := fmt.Sprintf(sp.TransmissionsPathFormat, testClient.Config.ApiVersion)
	testMux.HandleFunc(path+"/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write(txBytes)
	})

	tx1 := &sp.Transmission{ID: "42"}
	res, err := testClient.Transmission(tx1)
	if err != nil {
		testFailVerbose(t, res, "Transmission GET failed")
	}

	res, err = testClient.TransmissionContext(nil, tx1)
	if err != nil {
		testFailVerbose(t, res, "Transmission GET failed")
	}

	if tx1.CampaignID != tx.CampaignID {
		testFailVerbose(t, res, "CampaignIDs do not match")
	}
}

// Assert that options are actually ... optional,
// and that unspecified options don't default to their zero values.
func TestTransmissionOptions(t *testing.T) {
	var jsonb []byte
	var err error
	var opt bool

	tx := &sp.Transmission{}
	to := &sp.TxOptions{InlineCSS: &opt}
	tx.Options = to

	jsonb, err = json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(jsonb, []byte(`"options":{"inline_css":false}`)) {
		t.Fatalf("expected inline_css option to be false:\n%s", string(jsonb))
	}

	opt = true
	jsonb, err = json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(jsonb, []byte(`"options":{"inline_css":true}`)) {
		t.Fatalf("expected inline_css option to be true:\n%s", string(jsonb))
	}
}

package main

import (
	"encoding/json"
	"fmt"
	sp "github.com/SparkPost/gosparkpost"
	"os"
)

func main() {
	content := sp.Content{
		From:    "test@example.com",
		Subject: "transmission options example message",
		Text:    "This is a transmissions options example",
	}

	initialOpenTracking := false
	openTracking := true
	clickTracking := true
	options := sp.TmplOptions{
		InitialOpenTracking: &initialOpenTracking,
		OpenTracking: &openTracking,
		ClickTracking: &clickTracking,
	}

	txOptions := &sp.TxOptions{
		TmplOptions: options,
		IPPool: "my-ip-pool",
	}

	tx := &sp.Transmission{
		Recipients: []sp.Recipient{{Address: "optionstest@test.com.sink.sparkpostmail.com"}},
	}
	tx.Content = content
	tx.Options = txOptions

	jsonBytes, err := json.Marshal(tx)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stdout, string(jsonBytes))

}


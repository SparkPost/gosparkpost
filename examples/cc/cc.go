package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	sp "github.com/SparkPost/gosparkpost"
)

func main() {
	to := []string{
		"to1@test.com.sink.sparkpostmail.com",
		"to2@test.com.sink.sparkpostmail.com",
	}
	headerTo := strings.Join(to, ",")

	cc := []string{
		"cc1@test.com.sink.sparkpostmail.com",
		"cc2@test.com.sink.sparkpostmail.com",
	}

	bcc := []string{
		"bcc1@test.com.sink.sparkpostmail.com",
		"bcc2@test.com.sink.sparkpostmail.com",
	}

	content := sp.Content{
		From:    "test@example.com",
		Subject: "cc/bcc example message",
		Text:    "This is a cc/bcc example",
	}

	tx := &sp.Transmission{
		Recipients: []sp.Recipient{},
	}

	if len(to) > 0 {
		for _, t := range to {
			tx.Recipients = append(tx.Recipients.([]sp.Recipient), sp.Recipient{
				Address: sp.Address{Email: t, HeaderTo: headerTo},
			})
		}
	}

	if len(cc) > 0 {
		for _, c := range cc {
			tx.Recipients = append(tx.Recipients.([]sp.Recipient), sp.Recipient{
				Address: sp.Address{Email: c, HeaderTo: headerTo},
			})
		}
		// add cc header to content
		if content.Headers == nil {
			content.Headers = map[string]string{}
		}
		content.Headers["cc"] = strings.Join(cc, ",")
	}

	if len(bcc) > 0 {
		for _, b := range bcc {
			tx.Recipients = append(tx.Recipients.([]sp.Recipient), sp.Recipient{
				Address: sp.Address{Email: b, HeaderTo: headerTo},
			})
		}
	}

	tx.Content = content

	jsonBytes, err := json.Marshal(tx)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stdout, string(jsonBytes))

}

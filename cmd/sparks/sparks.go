// Sparks is a command-line tool for quickly sending email using SparkPost
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	sparkpost "github.com/SparkPost/gosparkpost"
)

var from = flag.String("from", "default@sparkpostbox.com", "where the mail came from")
var to = flag.String("to", "", "where the mail goes to")
var subject = flag.String("subject", "", "email subject")
var htmlFile = flag.String("html", "", "file containing html content")
var textFile = flag.String("text", "", "file containing text content")
var subsFile = flag.String("subs", "", "file containing substitution data (json object)")
var inline = flag.Bool("inline-css", false, "automatically inline css")
var dryrun = flag.Bool("dry-run", false, "dump json that would be sent to server")
var url = flag.String("url", "", "base url for api requests (optional)")

func main() {
	apiKey := os.Getenv("SPARKPOST_API_KEY")
	if strings.TrimSpace(apiKey) == "" {
		log.Fatal("FATAL: API key not found in environment!\n")
	}

	flag.Parse()

	if to == nil || strings.TrimSpace(*to) == "" {
		log.Fatal("SUCCESS: send mail to nobody!\n")
	}

	hasHtml := !(htmlFile == nil || strings.TrimSpace(*htmlFile) == "")
	hasText := !(textFile == nil || strings.TrimSpace(*textFile) == "")
	hasSubs := !(subsFile == nil || strings.TrimSpace(*subsFile) == "")

	if !hasHtml && !hasText {
		log.Fatal("FATAL: must specify one of --html or --text!\n")
	}

	cfg := &sparkpost.Config{ApiKey: apiKey}
	if url != nil && *url != "" {
		if !strings.HasPrefix(*url, "https://") {
			log.Fatal("FATAL: base url must be https!\n")
		}
		cfg.BaseUrl = *url
	}

	var sparky sparkpost.Client
	err := sparky.Init(cfg)
	if err != nil {
		log.Fatalf("SparkPost client init failed: %s\n", err)
	}

	content := sparkpost.Content{
		From:    *from,
		Subject: *subject,
	}
	if hasHtml {
		htmlBytes, err := ioutil.ReadFile(*htmlFile)
		if err != nil {
			log.Fatal(err)
		}
		content.HTML = string(htmlBytes)
	}

	if hasText {
		textBytes, err := ioutil.ReadFile(*textFile)
		if err != nil {
			log.Fatal(err)
		}
		content.Text = string(textBytes)
	}

	tx := &sparkpost.Transmission{
		Recipients: []string{*to},
		Content:    content,
	}

	if hasSubs {
		subsBytes, err := ioutil.ReadFile(*subsFile)
		if err != nil {
			log.Fatal(err)
		}
		recip := sparkpost.Recipient{Address: *to, SubstitutionData: json.RawMessage{}}
		err = json.Unmarshal(subsBytes, &recip.SubstitutionData)
		if err != nil {
			log.Fatal(err)
		}
		tx.Recipients = []sparkpost.Recipient{recip}
	}

	if inline != nil && *inline {
		if tx.Options == nil {
			tx.Options = &sparkpost.TxOptions{}
		}
		tx.Options.InlineCSS = true
	}

	if dryrun != nil && *dryrun {
		jsonBytes, err := json.Marshal(tx)
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write(jsonBytes)
		os.Exit(0)
	}

	id, req, err := sparky.Send(tx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("HTTP [%s] TX %s\n", req.HTTP.Status, id)
}

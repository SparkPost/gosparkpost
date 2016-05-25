// Sparks is a command-line tool for quickly sending email using SparkPost
package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	sparkpost "github.com/SparkPost/gosparkpost"
)

var from = flag.String("from", "default@sparkpostbox.com", "where the mail came from")
var to = flag.String("to", "", "where the mail goes to")
var cc = flag.String("cc", "", "carbon copy this address")
var bcc = flag.String("bcc", "", "blind carbon copy this address")
var subject = flag.String("subject", "", "email subject")
var htmlFlag = flag.String("html", "", "string/filename containing html content")
var textFlag = flag.String("text", "", "string/filename containing text content")
var subsFlag = flag.String("subs", "", "string/filename containing substitution data (json object)")
var imgFile = flag.String("img", "", "mimetype:cid:path for image to include")
var sendDelay = flag.String("send-delay", "", "delay delivery the specified amount of time")
var inline = flag.Bool("inline-css", false, "automatically inline css")
var dryrun = flag.Bool("dry-run", false, "dump json that would be sent to server")
var url = flag.String("url", "", "base url for api requests (optional)")

func main() {
	apiKey := os.Getenv("SPARKPOST_API_KEY")
	if strings.TrimSpace(apiKey) == "" {
		log.Fatal("FATAL: API key not found in environment!\n")
	}

	flag.Parse()

	if strings.TrimSpace(*to) == "" {
		log.Fatal("SUCCESS: send mail to nobody!\n")
	}

	hasHtml := strings.TrimSpace(*htmlFlag) != ""
	hasText := strings.TrimSpace(*textFlag) != ""
	hasSubs := strings.TrimSpace(*subsFlag) != ""
	hasImg := strings.TrimSpace(*imgFile) != ""

	if !hasHtml && !hasText {
		log.Fatal("FATAL: must specify one of --html or --text!\n")
	}

	cfg := &sparkpost.Config{ApiKey: apiKey}
	if strings.TrimSpace(*url) != "" {
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
		if strings.Contains(*htmlFlag, "/") {
			// read file to get html
			htmlBytes, err := ioutil.ReadFile(*htmlFlag)
			if err != nil {
				log.Fatal(err)
			}
			content.HTML = string(htmlBytes)
		} else {
			// html string passed on command line
			content.HTML = *htmlFlag
		}
	}

	if hasText {
		if strings.Contains(*textFlag, "/") {
			// read file to get text
			textBytes, err := ioutil.ReadFile(*textFlag)
			if err != nil {
				log.Fatal(err)
			}
			content.Text = string(textBytes)
		} else {
			// text string passed on command line
			content.Text = *textFlag
		}
	}

	if hasImg {
		imgra := strings.SplitN(*imgFile, ":", 3)
		if len(imgra) != 3 {
			log.Fatalf("--img format is mimetype:cid:path")
		}
		imgBytes, err := ioutil.ReadFile(imgra[2])
		if err != nil {
			log.Fatal(err)
		}
		img := sparkpost.InlineImage{
			MIMEType: imgra[0],
			Filename: imgra[1],
			B64Data:  base64.StdEncoding.EncodeToString(imgBytes),
		}
		content.InlineImages = append(content.InlineImages, img)
	}

	tx := &sparkpost.Transmission{}

	hasCc := strings.TrimSpace(*cc) != ""
	hasBcc := strings.TrimSpace(*bcc) != ""

	if hasCc {
		// need to set `header_to`; can't do that with string recipients
		tx.Recipients = []sparkpost.Recipient{
			{Address: sparkpost.Address{Email: *to}},
			{Address: sparkpost.Address{Email: *cc, HeaderTo: *to}},
		}
		if content.Headers == nil {
			content.Headers = map[string]string{}
		}
		content.Headers["cc"] = *cc
		if hasBcc {
			tx.Recipients = append(tx.Recipients.([]sparkpost.Recipient),
				sparkpost.Recipient{
					Address: sparkpost.Address{Email: *bcc, HeaderTo: *to}})
		}
	} else if hasBcc {
		tx.Recipients = []sparkpost.Recipient{
			{Address: sparkpost.Address{Email: *to}},
			{Address: sparkpost.Address{Email: *bcc, HeaderTo: *to}},
		}
	} else {
		tx.Recipients = []string{*to}
	}
	tx.Content = content

	if hasSubs {
		var subsBytes []byte
		if strings.Contains(*subsFlag, "/") {
			// read file to get substitution data
			subsBytes, err = ioutil.ReadFile(*subsFlag)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			subsBytes = []byte(*subsFlag)
		}

		recip := sparkpost.Recipient{Address: *to, SubstitutionData: json.RawMessage{}}
		err = json.Unmarshal(subsBytes, &recip.SubstitutionData)
		if err != nil {
			log.Fatal(err)
		}
		// FIXME: this clobbers any cc/bcc recipients
		tx.Recipients = []sparkpost.Recipient{recip}
	}

	if strings.TrimSpace(*sendDelay) != "" {
		if tx.Options == nil {
			tx.Options = &sparkpost.TxOptions{}
		}
		dur, err := time.ParseDuration(*sendDelay)
		if err != nil {
			log.Fatal(err)
		}
		start := sparkpost.RFC3339(time.Now().Add(dur))
		tx.Options.StartTime = &start
	}

	if *inline != false {
		if tx.Options == nil {
			tx.Options = &sparkpost.TxOptions{}
		}
		tx.Options.InlineCSS = true
	}

	if *dryrun != false {
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

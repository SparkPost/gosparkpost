// Sparks is a command-line tool for quickly sending email using SparkPost.
// It's like swaks, and apiaks sounded awkward.
package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	sp "github.com/SparkPost/gosparkpost"
)

type Strings []string

func (s *Strings) String() string {
	return strings.Join([]string(*s), ",")
}

func (s *Strings) Set(value string) error {
	*s = append([]string(*s), value)
	return nil
}

func main() {
	var to Strings
	var cc Strings
	var bcc Strings
	var headers Strings
	var images Strings
	var attachments Strings

	flag.Var(&to, "to", "where the mail goes to")
	flag.Var(&cc, "cc", "carbon copy this address")
	flag.Var(&bcc, "bcc", "blind carbon copy this address")
	flag.Var(&headers, "header", "custom header for your content")
	flag.Var(&images, "img", "mimetype:cid:path for image to include")
	flag.Var(&attachments, "attach", "mimetype:name:path for file to attach")

	var from = flag.String("from", "default@sparkpostbox.com", "where the mail came from")
	var subject = flag.String("subject", "", "email subject")
	var htmlFlag = flag.String("html", "", "string/filename containing html content")
	var textFlag = flag.String("text", "", "string/filename containing text content")
	var rfc822Flag = flag.String("rfc822", "", "string/filename containing raw message")
	var subsFlag = flag.String("subs", "", "string/filename containing substitution data (json object)")
	var sendDelay = flag.String("send-delay", "", "delay delivery the specified amount of time")
	var inline = flag.Bool("inline-css", false, "automatically inline css")
	var dryrun = flag.Bool("dry-run", false, "dump json that would be sent to server")
	var url = flag.String("url", "", "base url for api requests (optional)")
	var help = flag.Bool("help", false, "display a help message")
	var httpDump = flag.Bool("httpdump", false, "dump out http request and response")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if len(to) <= 0 {
		log.Fatal("SUCCESS: send mail to nobody!\n")
	}

	apiKey := os.Getenv("SPARKPOST_API_KEY")
	if strings.TrimSpace(apiKey) == "" && *dryrun == false {
		log.Fatal("FATAL: API key not found in environment!\n")
	}

	hasHtml := strings.TrimSpace(*htmlFlag) != ""
	hasText := strings.TrimSpace(*textFlag) != ""
	hasRFC822 := strings.TrimSpace(*rfc822Flag) != ""
	hasSubs := strings.TrimSpace(*subsFlag) != ""

	// rfc822 must be specified by itself, i.e. no text or html
	if hasRFC822 && (hasHtml || hasText) {
		log.Fatal("FATAL: --rfc822 cannot be combined with --html or --text!\n")
	} else if !hasRFC822 && !hasHtml && !hasText {
		log.Fatal("FATAL: must specify one of --html or --text!\n")
	}

	cfg := &sp.Config{ApiKey: apiKey}
	if strings.TrimSpace(*url) != "" {
		if !strings.HasPrefix(*url, "https://") {
			log.Fatal("FATAL: base url must be https!\n")
		}
		cfg.BaseUrl = *url
	}
	if *httpDump {
		cfg.Verbose = true
	}

	var sparky sp.Client
	err := sparky.Init(cfg)
	if err != nil {
		log.Fatalf("SparkPost client init failed: %s\n", err)
	}

	content := sp.Content{
		From:    *from,
		Subject: *subject,
	}

	if hasRFC822 {
		// these are pulled from the raw message
		content.From = nil
		content.Subject = ""
		if strings.HasPrefix(*rfc822Flag, "/") || strings.HasPrefix(*rfc822Flag, "./") {
			// read file to get raw message
			rfc822Bytes, err := ioutil.ReadFile(*rfc822Flag)
			if err != nil {
				log.Fatal(err)
			}
			content.EmailRFC822 = string(rfc822Bytes)
		} else {
			// raw message string passed on command line
			content.EmailRFC822 = *rfc822Flag
		}
	}

	if hasHtml {
		if strings.HasPrefix(*htmlFlag, "/") || strings.HasPrefix(*htmlFlag, "./") {
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
		if strings.HasPrefix(*textFlag, "/") || strings.HasPrefix(*textFlag, "./") {
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

	if len(images) > 0 {
		for _, imgStr := range images {
			img := strings.SplitN(imgStr, ":", 3)
			if len(img) != 3 {
				log.Fatalf("--img format is mimetype:cid:path")
			}
			imgBytes, err := ioutil.ReadFile(img[2])
			if err != nil {
				log.Fatal(err)
			}
			iimg := sp.InlineImage{
				MIMEType: img[0],
				Filename: img[1],
				B64Data:  base64.StdEncoding.EncodeToString(imgBytes),
			}
			content.InlineImages = append(content.InlineImages, iimg)
		}
	}

	if len(attachments) > 0 {
		for _, attStr := range attachments {
			att := strings.SplitN(attStr, ":", 3)
			if len(att) != 3 {
				log.Fatalf("--attach format is mimetype:name:path")
			}
			attBytes, err := ioutil.ReadFile(att[2])
			if err != nil {
				log.Fatal(err)
			}
			attach := sp.Attachment{
				MIMEType: att[0],
				Filename: att[1],
				B64Data:  base64.StdEncoding.EncodeToString(attBytes),
			}
			content.Attachments = append(content.Attachments, attach)
		}
	}

	tx := &sp.Transmission{}

	var subJson *json.RawMessage
	if hasSubs {
		var subsBytes []byte
		if strings.HasPrefix(*subsFlag, "/") || strings.HasPrefix(*subsFlag, "./") {
			// read file to get substitution data
			subsBytes, err = ioutil.ReadFile(*subsFlag)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			subsBytes = []byte(*subsFlag)
		}

		subJson = &json.RawMessage{}
		err = json.Unmarshal(subsBytes, subJson)
		if err != nil {
			log.Fatal(err)
		}
	}

	headerTo := strings.Join(to, ",")

	tx.Recipients = []sp.Recipient{}
	for _, r := range to {
		var recip sp.Recipient = sp.Recipient{Address: sp.Address{Email: r, HeaderTo: headerTo}}
		if hasSubs {
			recip.SubstitutionData = subJson
		}
		tx.Recipients = append(tx.Recipients.([]sp.Recipient), recip)
	}

	if len(cc) > 0 {
		for _, r := range cc {
			var recip sp.Recipient = sp.Recipient{Address: sp.Address{Email: r, HeaderTo: headerTo}}
			if hasSubs {
				recip.SubstitutionData = subJson
			}
			tx.Recipients = append(tx.Recipients.([]sp.Recipient), recip)
		}
		if content.Headers == nil {
			content.Headers = map[string]string{}
		}
		content.Headers["cc"] = strings.Join(cc, ",")
	}

	if len(bcc) > 0 {
		for _, r := range bcc {
			var recip sp.Recipient = sp.Recipient{Address: sp.Address{Email: r, HeaderTo: headerTo}}
			if hasSubs {
				recip.SubstitutionData = subJson
			}
			tx.Recipients = append(tx.Recipients.([]sp.Recipient), recip)
		}
	}

	if len(headers) > 0 {
		if content.Headers == nil {
			content.Headers = map[string]string{}
		}
		hb := regexp.MustCompile(`:\s*`)
		for _, hstr := range headers {
			hra := hb.Split(hstr, 2)
			content.Headers[hra[0]] = hra[1]
		}
	}

	tx.Content = content

	if strings.TrimSpace(*sendDelay) != "" {
		if tx.Options == nil {
			tx.Options = &sp.TxOptions{}
		}
		dur, err := time.ParseDuration(*sendDelay)
		if err != nil {
			log.Fatal(err)
		}
		start := sp.RFC3339(time.Now().Add(dur))
		tx.Options.StartTime = &start
	}

	if *inline != false {
		if tx.Options == nil {
			tx.Options = &sp.TxOptions{}
		}
		tx.Options.InlineCSS = inline
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

	if *httpDump {
		if reqDump, ok := req.Verbose["http_requestdump"]; ok {
			os.Stdout.Write([]byte(reqDump))
		} else {
			os.Stdout.Write([]byte("*** No request dump available! ***\n\n"))
		}

		if resDump, ok := req.Verbose["http_responsedump"]; ok {
			os.Stdout.Write([]byte(resDump))
			os.Stdout.Write([]byte("\n"))
		} else {
			os.Stdout.Write([]byte("*** No response dump available! ***\n"))
		}
	} else {
		log.Printf("HTTP [%s] TX %s\n", req.HTTP.Status, id)
	}
}

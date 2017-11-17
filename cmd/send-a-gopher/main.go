package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	gosp "github.com/SparkPost/gosparkpost"
)

func main() {
	var imgUrl = flag.String("url", "", "url of image to send in our email")
	var from = flag.String("from", "", "send using this from address")
	var to = flag.String("to", "", "send to this address")
	flag.Parse()

	if *imgUrl == "" {
		log.Fatal("Successfully sent nothing, but that's probably not what you wanted? Check out the --url option.")
	} else if *from == "" {
		log.Fatal("Must specify from address with --from option.")
	} else if *to == "" {
		log.Fatal("Must specify to address with --to option.")
	}

	parsed, err := url.Parse(*imgUrl)
	if err != nil {
		log.Fatal("That's very much not a URL. Please try again.")
	} else if parsed.Scheme == "" || parsed.Host == "" {
		log.Fatal("That doesn't really look like a URL. Please try again.")
	}

	res, err := http.Get(*imgUrl)
	if err != nil {
		log.Fatalf("Couldn't get that url: %s\n", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Couldn't read that url: %s\n", err)
	}
	ctype := http.DetectContentType(body)
	filename := path.Base(*imgUrl)
	iimg := gosp.InlineImage{
		MIMEType: ctype,
		Filename: filename,
		B64Data:  base64.StdEncoding.EncodeToString(body),
	}

	cfg := &gosp.Config{ApiKey: os.Getenv("SPARKPOST_API_KEY")}
	var sp gosp.Client
	err = sp.Init(cfg)
	if err != nil {
		log.Fatalf("SparkPost client init failed: %s\n", err)
	}

	html := fmt.Sprintf(`Here's that gopher you maybe asked for!<br/><img src="cid:%s" />`, filename)
	content := gosp.Content{
		From:    *from,
		Subject: "That gopher",
		HTML:    html,
	}
	content.InlineImages = append(content.InlineImages, iimg)
	tx := &gosp.Transmission{
		Content:    content,
		Recipients: []string{*to},
	}

	id, _, err := sp.Send(tx)
	if err != nil {
		log.Fatalf("Couldn't send that gopher image this time: %s\n", err)
	}
	log.Printf("Sent that gopher image to [%s]! (%s)\n", *to, id)
}

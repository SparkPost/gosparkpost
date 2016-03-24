package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var filename = flag.String("file", "", "path to email with a text/html part")
var dumpArf = flag.Bool("arf", false, "dump out multipart/report message")
var send = flag.String("send", "", "send the fbl report to host[:port]")
var verbose = flag.Bool("verbose", false, "print out lots of messages")

var cidPattern *regexp.Regexp = regexp.MustCompile(`"customer_id"\s*:\s*"(\d+)"`)
var fromPattern *regexp.Regexp = regexp.MustCompile(`"friendly_from"\s*:\s*"([^"\s]+)"`)
var toPattern *regexp.Regexp = regexp.MustCompile(`"r"\s*:\s*"([^"\s]+)"`)

func main() {
	flag.Parse()

	if filename == nil || strings.TrimSpace(*filename) == "" {
		log.Fatal("--file is required")
	}

	fh, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}

	msg, err := mail.ReadMessage(fh)
	if err != nil {
		log.Fatal(err)
	}

	b64hdr := strings.Replace(msg.Header.Get("X-MSFBL"), " ", "", -1)
	if verbose != nil && *verbose == true {
		log.Printf("X-MSFBL: %s\n", b64hdr)
	}

	var dec []byte
	b64 := base64.StdEncoding
	if strings.Index(b64hdr, "|") >= 0 {
		// Everything before the pipe is an encoded hmac
		// TODO: verify contents using hmac
		encs := strings.Split(b64hdr, "|")
		dec, err = b64.DecodeString(encs[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		dec, err = b64.DecodeString(b64hdr)
		if err != nil {
			log.Fatal(err)
		}
	}

	cidMatches := cidPattern.FindSubmatch(dec)
	if cidMatches == nil || len(cidMatches) < 2 {
		log.Fatalf("No key \"customer_id\" in X-MSFBL header:\n%s\n", string(dec))
	}
	cid, err := strconv.Atoi(string(cidMatches[1]))
	if err != nil {
		log.Fatal(err)
	}

	fromMatches := fromPattern.FindSubmatch(dec)
	if fromMatches == nil || len(fromMatches) < 2 {
		log.Fatalf("No key \"friendly_from\" in X-MSFBL header:\n%s\n", string(dec))
	}

	toMatches := toPattern.FindSubmatch(dec)
	if toMatches == nil || len(toMatches) < 2 {
		log.Fatalf("No key \"r\" (recipient) in X-MSFBL header:\n%s\n", string(dec))
	}

	if verbose != nil && *verbose == true {
		log.Printf("Decoded (%d):\n%s\n", cid, string(dec))
	}

	// from/to are opposite here, since we're simulating a reply
	to := string(fromMatches[1])
	from := string(toMatches[1])
	arf := BuildArf(from, to, b64hdr, cid)

	if dumpArf != nil && *dumpArf == true {
		fmt.Fprintf(os.Stdout, "%s", arf)
	}

	// TODO: send
}

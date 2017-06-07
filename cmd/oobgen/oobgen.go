package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/SparkPost/gosparkpost/helpers/loadmsg"
)

func main() {
	var filename = flag.String("file", "", "path to raw email")
	var send = flag.Bool("send", false, "send fbl report")
	var verboseOpt = flag.Bool("verbose", false, "print out lots of messages")

	flag.Parse()
	var verbose bool
	if *verboseOpt == true {
		verbose = true
	}

	if *filename == "" {
		log.Fatal("--file is required")
	}

	msg := loadmsg.Message{Filename: *filename}
	err := msg.Load()
	if err != nil {
		log.Fatal(err)
	}

	atIdx := strings.Index(msg.ReturnPath.Address, "@")
	if atIdx < 0 {
		log.Fatalf("Unsupported Return-Path header [%s]\n", msg.ReturnPath.Address)
	}
	oobDomain := msg.ReturnPath.Address[atIdx+1:]
	if verbose == true {
		log.Printf("Got domain [%s] from Return-Path\n", oobDomain)
	}

	fileBytes, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatal(err)
	}

	// from/to are opposite here, since we're simulating a reply
	to := msg.ReturnPath.Address
	from, err := mail.ParseAddress(msg.Message.Header.Get("From"))
	if err != nil {
		log.Fatal(err)
	}
	oob := BuildOob(to, from.Address, string(fileBytes))

	mxs, err := net.LookupMX(oobDomain)
	if err != nil {
		log.Fatal(err)
	}
	if mxs == nil || len(mxs) <= 0 {
		log.Fatal("No MXs for [%s]\n", oobDomain)
	}
	if verbose == true {
		log.Printf("Got MX [%s] for [%s]\n", mxs[0].Host, oobDomain)
	}
	smtpHost := fmt.Sprintf("%s:smtp", mxs[0].Host)

	if *send == true {
		log.Printf("Sending OOB from [%s] to [%s] via [%s]...\n",
			from, to, smtpHost)
		err = smtp.SendMail(smtpHost, nil, from.Address, []string{to}, []byte(oob))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Sent.\n")
	} else {
		if verbose == true {
			log.Printf("Would send OOB from [%s] to [%s] via [%s]\n",
				from.Address, to, smtpHost)
		}
	}
}

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/SparkPost/gosparkpost/helpers/loadmsg"
	"github.com/SparkPost/gosparkpost/helpers/smtptls"
)

func main() {
	var filename = flag.String("file", "", "path to raw email")
	var send = flag.Bool("send", false, "send oob bounce")
	var port = flag.Int("port", 25, "port for outbound smtp")
	var serverName = flag.String("servername", "", "override tls servername")
	var verboseOpt = flag.Bool("verbose", false, "print out lots of messages")
	var _smtpHost = flag.String("smtphost", "", "override smtp host")

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

	smtpHost := *_smtpHost
	if smtpHost == "" {
		// not set, auto-detect using domain from return-path header
		mxs, err := net.LookupMX(oobDomain)
		if err != nil {
			log.Fatal(err)
		}
		if mxs == nil || len(mxs) <= 0 {
			log.Fatalf("No MXs for [%s]\n", oobDomain)
		}
		if verbose == true {
			log.Printf("Got MX [%s] for [%s]\n", mxs[0].Host, oobDomain)
		}
		smtpHost = fmt.Sprintf("%s:%d", mxs[0].Host, *port)
	}

	to := msg.ReturnPath.Address
	// from/to are opposite here, since we're simulating a reply
	from, err := mail.ParseAddress(msg.Message.Header.Get("From"))
	if err != nil {
		log.Fatal(err)
	}
	oob := BuildOob(to, from.Address, string(fileBytes))

	var smtpTLS *smtp.Client
	if *serverName != "" {
		tlsc := &tls.Config{ServerName: *serverName}
		smtpTLS, err = smtptls.Connect(smtpHost, tlsc)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *send == true {
		log.Printf("Sending OOB from [%s] to [%s] via [%s]...\n",
			from, to, smtpHost)
		if *serverName != "" {
			err = smtptls.SendMail(smtpTLS, from.Address, []string{to}, []byte(oob))
		} else {
			err = smtp.SendMail(smtpHost, nil, from.Address, []string{to}, []byte(oob))
		}
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

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"

	"github.com/SparkPost/gosparkpost/helpers/loadmsg"
)

func main() {
	var filename = flag.String("file", "", "path to raw email")
	var dumpArf = flag.Bool("arf", false, "dump out multipart/report message")
	var send = flag.Bool("send", false, "send fbl report")
	var fblAddress = flag.String("fblto", "", "where to deliver the fbl report")
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

	if *fblAddress != "" {
		msg.SetReturnPath(*fblAddress)
	}

	atIdx := strings.Index(msg.ReturnPath.Address, "@")
	if atIdx < 0 {
		log.Fatalf("Unsupported Return-Path header [%s]\n", msg.ReturnPath.Address)
	}
	fblDomain := msg.ReturnPath.Address[atIdx+1:]
	fblTo := fmt.Sprintf("fbl@%s", fblDomain)
	if verbose == true {
		if *fblAddress != "" {
			log.Printf("Got domain [%s] from --fblto\n", fblDomain)
		} else {
			log.Printf("Got domain [%s] from Return-Path\n", fblDomain)
		}
	}

	// from/to are opposite here, since we're simulating a reply
	fblFrom := string(msg.Recipient)
	arf := BuildArf(fblFrom, fblTo, msg.MSFBL, msg.CustID)

	if *dumpArf == true {
		fmt.Fprintf(os.Stdout, "%s", arf)
	}

	mxs, err := net.LookupMX(fblDomain)
	if err != nil {
		log.Fatal(err)
	}
	if mxs == nil || len(mxs) <= 0 {
		log.Fatal("No MXs for [%s]\n", fblDomain)
	}
	if verbose == true {
		log.Printf("Got MX [%s] for [%s]\n", mxs[0].Host, fblDomain)
	}
	smtpHost := fmt.Sprintf("%s:smtp", mxs[0].Host)

	if *send == true {
		log.Printf("Sending FBL from [%s] to [%s] via [%s]...\n",
			fblFrom, fblTo, smtpHost)
		err = smtp.SendMail(smtpHost, nil, fblFrom, []string{fblTo}, []byte(arf))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Sent.\n")
	} else {
		if verbose == true {
			log.Printf("Would send FBL from [%s] to [%s] via [%s]\n",
				fblFrom, fblTo, smtpHost)
		}
	}
}

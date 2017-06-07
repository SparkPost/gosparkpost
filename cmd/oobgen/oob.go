package main

import (
	"fmt"
	"strings"
	"time"
)

// FIXME: allow swapping out the error message
var OobFormat string = `From: %s
Date: Mon, 02 Jan 2006 15:04:05 MST
Subject: Returned mail: see transcript for details
Auto-Submitted: auto-generated (failure)
To: %s
Content-Type: multipart/report; report-type=delivery-status;
	boundary="%s"

This is a MIME-encapsulated message

--%s

The original message was received at Mon, 02 Jan 2006 15:04:05 -0700
from example.com.sink.sparkpostmail.com [52.41.116.105]

   ----- The following addresses had permanent fatal errors -----
<%s>
    (reason: 550 5.0.0 <%s>... User unknown)

   ----- Transcript of session follows -----
... while talking to %s:
>>> DATA
<<< 550 5.0.0 <%s>... User unknown
550 5.1.1 <%s>... User unknown
<<< 503 5.0.0 Need RCPT (recipient)

--%s
Content-Type: message/delivery-status

Reporting-MTA: dns; %s
Received-From-MTA: DNS; %s
Arrival-Date: Mon, 02 Jan 2006 15:04:05 MST

Final-Recipient: RFC822; %s
Action: failed
Status: 5.0.0
Remote-MTA: DNS; %s
Diagnostic-Code: SMTP; 550 5.0.0 <%s>... User unknown
Last-Attempt-Date: Mon, 02 Jan 2006 15:04:05 MST

--%s
Content-Type: message/rfc822

%s

--%s--
`

func BuildOob(from, to, rawMsg string) string {
	boundary := fmt.Sprintf("_----%d===_61/00-25439-267B0055", time.Now().Unix())
	fromDomain := from[strings.Index(from, "@")+1:]
	toDomain := to[strings.Index(to, "@")+1:]
	msg := fmt.Sprintf(OobFormat,
		from, to, boundary,
		boundary, to, to, toDomain, to, to,
		boundary, toDomain, fromDomain, to, toDomain, to,
		boundary, rawMsg,
		boundary)
	return msg
}

package loadmsg

import (
	"encoding/base64"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

type Message struct {
	Filename  string
	File      *os.File
	Message   *mail.Message
	Json      []byte
	CustID    int
	Recipient []byte
}

func (m *Message) Load() error {
	var err error

	m.File, err = os.Open(m.Filename)
	if err != nil {
		return err
	}

	m.Message, err = mail.ReadMessage(m.File)
	if err != nil {
		return err
	}

	b64hdr := strings.Replace(m.Message.Header.Get("X-MSFBL"), " ", "", -1)

	if strings.Index(b64hdr, "|") >= 0 {
		// Everything before the pipe is an encoded HMAC
		// TODO: verify contents using HMAC
		b64hdr = strings.Split(b64hdr, "|")[1]
	}

	m.Json, err = base64.StdEncoding.DecodeString(b64hdr)
	if err != nil {
		return err
	}

	var cid []byte
	cid, _, _, err = jsonparser.Get(m.Json, "customer_id")
	if err != nil {
		return err
	}
	m.CustID, err = strconv.Atoi(string(cid))
	if err != nil {
		return err
	}

	m.Recipient, _, _, err = jsonparser.Get(m.Json, "r")
	if err != nil {
		return err
	}

	return nil
}

package loadmsg

import (
	"encoding/base64"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

type Message struct {
	Filename   string
	File       *os.File
	Message    *mail.Message
	MSFBL      string
	Json       []byte
	CustID     int
	Recipient  []byte
	ReturnPath *mail.Address
}

func (m *Message) Load() error {
	var err error

	m.File, err = os.Open(m.Filename)
	if err != nil {
		return errors.Wrap(err, "opening file")
	}

	m.Message, err = mail.ReadMessage(m.File)
	if err != nil {
		return errors.Wrap(err, "parsing message")
	}

	if m.ReturnPath == nil {
		err = m.SetReturnPath(m.Message.Header.Get("Return-Path"))
		if err != nil {
			return errors.Wrap(err, "setting return path")
		}
	}

	m.MSFBL = strings.Replace(m.Message.Header.Get("X-MSFBL"), " ", "", -1)

	if m.MSFBL == "" {
		// early return if there isn't a MSFBL header
		return nil
	}

	if strings.Index(m.MSFBL, "|") >= 0 {
		// Everything before the pipe is an encoded HMAC
		// TODO: verify contents using HMAC
		m.MSFBL = strings.Split(m.MSFBL, "|")[1]
	}

	m.Json, err = base64.StdEncoding.DecodeString(m.MSFBL)
	if err != nil {
		return errors.Wrap(err, "decoding fbl")
	}

	var cid []byte
	cid, _, _, err = jsonparser.Get(m.Json, "customer_id")
	if err != nil {
		return errors.Wrap(err, "getting customer_id")
	}
	m.CustID, err = strconv.Atoi(string(cid))
	if err != nil {
		return errors.Wrap(err, "int-ifying customer_id")
	}

	m.Recipient, _, _, err = jsonparser.Get(m.Json, "r")
	if err != nil {
		return errors.Wrap(err, "getting recipient")
	}

	return nil
}

func (m *Message) SetReturnPath(addr string) (err error) {
	if !strings.Contains(addr, "@") {
		return errors.Errorf("Unsupported Return-Path header: no @")
	}
	m.ReturnPath, err = mail.ParseAddress(addr)
	if err != nil {
		return errors.Wrap(err, "parsing return path")
	}
	return nil
}

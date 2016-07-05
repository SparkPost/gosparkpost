package loadmsg

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
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
		return err
	}

	m.Message, err = mail.ReadMessage(m.File)
	if err != nil {
		return err
	}

	if m.ReturnPath == nil {
		err = m.SetReturnPath(m.Message.Header.Get("Return-Path"))
		if err != nil {
			return err
		}
	}

	m.MSFBL = strings.Replace(m.Message.Header.Get("X-MSFBL"), " ", "", -1)

	if strings.Index(m.MSFBL, "|") >= 0 {
		// Everything before the pipe is an encoded HMAC
		// TODO: verify contents using HMAC
		m.MSFBL = strings.Split(m.MSFBL, "|")[1]
	}

	m.Json, err = base64.StdEncoding.DecodeString(m.MSFBL)
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

func (m *Message) SetReturnPath(addr string) (err error) {
	if !strings.Contains(addr, "@") {
		return fmt.Errorf("Unsupported Return-Path header: no @")
	}
	m.ReturnPath, err = mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	return nil
}

package smtptls

import (
	"crypto/tls"
	"net/smtp"

	"github.com/pkg/errors"
)

func Connect(host string, tlsc tls.Config) (*smtp.Client, error) {
	client, err := smtp.Dial(host)
	if err != nil {
		return nil, errors.Wrap(err, "dialing smtp")
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(&tlsc); err != nil {
			return nil, errors.Wrap(err, "starting tls")
		}
	} else {
		return nil, errors.Errorf("%s doesn't advertise STARTTLS support", host)
	}

	return client, nil
}

func SendMail(c *smtp.Client, from string, to []string, msg []byte) error {
	if err := c.Mail(from); err != nil {
		return errors.Wrap(err, "sending MAIL FROM")
	}

	if err := c.Rcpt(to[0]); err != nil {
		return errors.Wrap(err, "sending RCPT TO")
	}

	w, err := c.Data()
	if err != nil {
		return errors.Wrap(err, "sending DATA")
	}

	_, err = w.Write(msg)
	if err != nil {
		return errors.Wrap(err, "sending message body")
	}

	err = w.Close()
	if err != nil {
		return errors.Wrap(err, "closing connection")
	}

	c.Quit()
	return nil
}

package smtp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/pkg/errors"
)

type Client interface {
	SendMail(from, to string, content []byte) error
}

type client struct {
	c    *config.SMTPConfig
	tc   *tls.Config
	auth smtp.Auth
}

func NewClient(c *config.SMTPConfig) Client {
	auth := smtp.PlainAuth("", c.User, c.Password, c.Host)

	tc := &tls.Config{
		ServerName:         c.Host,
		InsecureSkipVerify: c.InsecureSkipVerify,
	}

	return &client{c, tc, auth}
}

func (s *client) SendMail(from, to string, message []byte) error {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.c.Host, s.c.Port), s.tc)
	if err != nil {
		return errors.Wrap(err, "smtp.SendMail.Dial")
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, s.c.Host)
	if err != nil {
		log.Panic(errors.Wrap(err, "smtp.client.SendMail"))
	}

	if err = c.Auth(s.auth); err != nil {
		return errors.Wrap(err, "smtp.client.SendMail.Auth")
	}

	if err = c.Mail(from); err != nil {
		return errors.Wrap(err, "smtp.client.SendMail.Mail")
	}

	if err = c.Rcpt(to); err != nil {
		return errors.Wrap(err, "smtp.client.SendMail.Rcpt")
	}

	d, err := c.Data()
	if err != nil {
		return errors.Wrap(err, "smtp.client.SendMail.Data")
	}

	_, err = d.Write(message)
	if err != nil {
		return errors.Wrap(err, "smtp.client.SendMail.Write")
	}

	err = d.Close()
	if err != nil {
		return errors.Wrap(err, "smtp.client.SendMail.Close")
	}

	c.Quit()

	return nil
}

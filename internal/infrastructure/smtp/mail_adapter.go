package smtp

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/adapter"
)

type mailAdapter struct {
	c    *config.Config
	auth *smtp.Auth
	addr *string
}

func NewMailAdapter(c *config.Config) adapter.MailAdapter {
	auth := smtp.PlainAuth(c.MailAddress.Default, c.SMTP.User, c.SMTP.Password, c.Server.Host)
	addr := fmt.Sprintf("%s:%d", c.SMTP.Host, c.SMTP.Port)
	return &mailAdapter{c, &auth, &addr}
}

func (a *mailAdapter) SendConfirmAccount(receiver string, token string) error {
	link := fmt.Sprintf("%s/%s", a.c.Frontend.ConfirmAccountURL, token)

	content := fmt.Sprintf(`
		<h2>Confirm Account</h2>
		<p>Click <a target="_blank" href="%s">link</a> to activate your account</p>
	`, link)

	err := smtp.SendMail(*a.addr, *a.auth, a.c.MailAddress.Default, []string{receiver}, a.createMessage("confirm", receiver, content))
	if err != nil {
		return err
	}

	return nil
}

func (a *mailAdapter) SendPasswordReset(receiver string, token string) error {
	link := fmt.Sprintf("%s/%s", a.c.Frontend.ResetPasswordURL, token)

	msg := fmt.Sprintf(`
		<h2>Reset Password</h2>
		<p>Click <a target="_blank" href="%s">link</a> to reset your password</p>
	`, link)

	err := smtp.SendMail(*a.addr, *a.auth, a.c.MailAddress.Default, []string{receiver}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func (a *mailAdapter) createMessage(mailType string, receiver string, content string) []byte {
	msg := strings.Builder{}
	defer msg.Reset()

	msg.WriteString("To: " + receiver + "\r\n")
	msg.WriteString("From: " + a.c.MailAddress.Default + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	if mailType == "confirm" {
		msg.WriteString("Subject: [DashboardX] Confirm Account\r\n")
	} else if mailType == "reset" {
		msg.WriteString("Subject: [DashboardX] Reset Password\r\n")
	}
	msg.WriteString("\r\n")
	msg.WriteString("<div style=\"font-family: Verdana, sans-serif;\">")
	msg.WriteString(content)
	msg.WriteString("</div>")

	return []byte(msg.String())
}

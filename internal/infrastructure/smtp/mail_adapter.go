package smtp

import (
	"fmt"
	"strings"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/adapter"
	"github.com/Deve-Lite/DashboardX-API/pkg/smtp"
)

type mailAdapter struct {
	c *config.Config
	s smtp.Client
}

func NewMailAdapter(c *config.Config, s smtp.Client) adapter.MailAdapter {
	return &mailAdapter{c, s}
}

func (a *mailAdapter) SendConfirmAccount(receiver string, token string) error {
	link := fmt.Sprintf("%s/%s", a.c.Frontend.ConfirmAccountURL, token)

	content := fmt.Sprintf(`
		<h2>Confirm Account</h2>
		<p>Click <a target="_blank" rel="noreferrer nofollow" href="%s">link</a> to activate your account</p>
	`, link)

	message := a.createMessage("confirm", receiver, content)

	err := a.s.SendMail(a.c.MailAddress.Default, receiver, message)
	if err != nil {
		return err
	}

	return nil
}

func (a *mailAdapter) SendPasswordReset(receiver string, token string) error {
	link := fmt.Sprintf("%s/%s", a.c.Frontend.ResetPasswordURL, token)

	content := fmt.Sprintf(`
		<h2>Reset Password</h2>
		<p>Click <a target="_blank" rel="noreferrer nofollow" href="%s">link</a> to reset your password</p>
	`, link)

	message := a.createMessage("reset", receiver, content)

	err := a.s.SendMail(a.c.MailAddress.Default, receiver, message)
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

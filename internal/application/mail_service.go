package application

import (
	"log"

	"github.com/Deve-Lite/DashboardX-API/internal/domain/adapter"
	"github.com/pkg/errors"
)

type MailService interface {
	SendConfirmAccount(receiver string, token string)
	SendPasswordReset(receiver string, token string)
}

type mailService struct {
	ma adapter.MailAdapter
}

func NewMailService(ma adapter.MailAdapter) MailService {
	return &mailService{ma}
}

func (m *mailService) SendConfirmAccount(receiver string, token string) {
	if err := m.ma.SendConfirmAccount(receiver, token); err != nil {
		log.Print(errors.Wrap(err, "mailService.SendConfirmAccount"))
	}
}

func (m *mailService) SendPasswordReset(receiver string, token string) {
	if err := m.ma.SendPasswordReset(receiver, token); err != nil {
		log.Print(errors.Wrap(err, "mailService.SendPasswordReset"))
	}
}

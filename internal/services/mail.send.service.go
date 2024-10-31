package services

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/mikromolekula2002/Testovoe/internal/config"
)

type MailSendService struct {
	cfg *config.Config
}

func NewMailSendService(cfg *config.Config) *MailSendService {
	return &MailSendService{cfg}
}

// тут переименовать и написать логику
func (ms *MailSendService) SendMailWarning(subject string, content string, to []string) error {
	// отправка варнинга на почту
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", ms.cfg.SMTP.Name, ms.cfg.SMTP.Email_address)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to

	smtpAuth := smtp.PlainAuth("", ms.cfg.SMTP.Email_address, ms.cfg.SMTP.Email_password, ms.cfg.SMTP.Auth_address)
	return e.Send(ms.cfg.SMTP.Server_address, smtpAuth)
}

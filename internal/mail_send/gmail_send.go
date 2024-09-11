package mail_send

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
	) error
}

// Структура отправителя (name = псевдоним почты), (Address = адрес почты отправителя), (Password = для GMAIL это 16-значный пароль приложения)
type GmailSender struct {
	AuthAddress       string
	ServerAddress     string
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

// Объявляем нашу структуру
func NewGmailSender(name string,
	fromEmailAddress string,
	fromEmailPassword string,
	authAddress string,
	serverAddress string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
		AuthAddress:       authAddress,
		ServerAddress:     serverAddress,
	}
}

// Метод отправки сообщения
func (sender *GmailSender) SendEmail(subject string, content string, to []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, sender.AuthAddress)
	return e.Send(sender.ServerAddress, smtpAuth)
}

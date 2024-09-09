package mail_send

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type Smtp interface {
	SendEmailWarn(receiver string) error
}

// СМТП сервер - для рассылки Email Warning
type SMTPServer struct {
	SMTP_Domain string
	User        string
	ApiKey      string //Or password
	MailSender  string
}

// инициализация смтп сервера
func Init(smtpDomain, user, apiKey, mailSender string) *SMTPServer {
	return &SMTPServer{
		SMTP_Domain: smtpDomain,
		User:        user,
		ApiKey:      apiKey,
		MailSender:  mailSender,
	}
}

// Метод отправляющий Email Warning
func (s *SMTPServer) SendEmailWarn(receiver string) error {
	op := "mail_send.SendEmailWarn"

	// Создаем новое сообщение
	m := gomail.NewMessage()
	m.SetHeader("From", s.MailSender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "Auth_service - Предупреждение!")
	m.SetBody("text/plain", "Кто-то пытается получить ваш токен с другого `ip` адреса.")

	// Настраиваем SMTP сервер Mailtrap
	d := gomail.NewDialer(s.SMTP_Domain, 587, s.User, s.ApiKey)

	// Отправка письма
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("%s - Ошибка отправки сообщения на почту: \n%v", op, err)
	}
	return nil
}

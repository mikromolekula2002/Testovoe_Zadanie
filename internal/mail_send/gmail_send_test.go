package mail_send

import (
	"testing"

	"github.com/mikromolekula2002/Testovoe/internal/config"
	"github.com/stretchr/testify/require"
)

const (
	EmailRecipient = "kortymalik@gmail.com" //replace to your mail
)

func TestSendEmailWithGmail(t *testing.T) {
	//Загрузка конфига
	cfg := config.LoadConfig("./config/config.yaml")

	sender := NewGmailSender(cfg.SMTP.Name,
		cfg.SMTP.Email_address,
		cfg.SMTP.Email_password,
		cfg.SMTP.Auth_address,
		cfg.SMTP.Server_address)

	subject := "A test email"
	content := `
    <h1>Hello world</h1>
    <p>This is a test message from Auth_Service</p>
    `
	to := []string{EmailRecipient}

	err := sender.SendEmail(subject, content, to)
	require.NoError(t, err)
}

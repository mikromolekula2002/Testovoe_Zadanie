package mail_send

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	Name           = "Auth Service API"
	Email_address  = "INSERT"
	Email_password = "INSERT"
	Auth_address   = "smtp.gmail.com"
	Server_address = "smtp.gmail.com:587"
	EmailRecipient = "kortymalik@gmail.com" //replace to your mail
)

func TestSendEmailWithGmail(t *testing.T) {

	sender := NewGmailSender(Name,
		Email_address,
		Email_password,
		Auth_address,
		Server_address)

	subject := "A test email"
	content := `
    <h1>Hello world</h1>
    <p>This is a test message from Auth_Service</p>
    `
	to := []string{EmailRecipient}

	err := sender.SendEmail(subject, content, to)
	require.NoError(t, err)
}

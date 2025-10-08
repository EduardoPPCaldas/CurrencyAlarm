package email

import (
	"net/smtp"
	"os"
)

type EmailSender interface {
	SendEmail(to []string, message []byte) error
}

type smtpEmailSender struct {
}

func NewEmailSender() EmailSender {
	return &smtpEmailSender{}
}

func (es smtpEmailSender) SendEmail(to []string, message []byte) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	return err
}

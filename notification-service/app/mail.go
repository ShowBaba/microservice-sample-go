package app

import (
	"net/smtp"

	"github.com/showbaba/microservice-sample-go/shared"
)

func SendEmail(mail shared.Mail) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	message := mail.BuildMessage()
	// Create authentication
	auth := smtp.PlainAuth("", GetConfig().MailUsername, GetConfig().MailPassword, smtpHost)
	// Send actual message
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, mail.Sender, mail.To, []byte(message))
}

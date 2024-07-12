package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"vokki_cloud/internal/models"
)

func SendVerificationEmail(user models.User, token string) error {
	from := os.Getenv("FROM_EMAIL")
	appPassword := os.Getenv("FROM_EMAIL_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, appPassword, smtpHost)
	toList := []string{user.Email}
	subject := "Subject: Email Verification\n"
	body := fmt.Sprintf("Please verify your email by clicking the following link: http://localhost:8080/verify?token=%s", token)
	msg := []byte(subject + "\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, msg)
	if err != nil {
		log.Println("Error sending email: ", err)
		return err
	}
	return nil
}

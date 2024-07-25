package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	vokki_constants "vokki_cloud/internal/constants"
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
	body := fmt.Sprintf("Please verify your email by clicking the following link: https://3.145.165.81%s/?token=%s", vokki_constants.RouteVerifyEmail, token)
	msg := []byte(subject + "\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, msg)
	if err != nil {
		log.Println("Error sending email: ", err)
		return err
	}
	return nil
}

func SendPasswordResetEmail(user models.User, token string) error {
	from := os.Getenv("FROM_EMAIL")
	appPassword := os.Getenv("FROM_EMAIL_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, appPassword, smtpHost)
	toList := []string{user.Email}
	subject := "Subject: Password Reset\n"
	body := fmt.Sprintf("Please reset your password by clicking the following link: https://")

	msg := []byte(subject + "\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, msg)
	if err != nil {
		log.Println("Error sending email: ", err)
		return err
	}
	return nil

}

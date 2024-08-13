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

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"Priority: Urgent\n" +
		"X-Priority: 1\n" +
		"Importance: High\n"

	body := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Email Verification</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				margin: 0;
				padding: 0;
			}
			.container {
				width: 100%%;
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
				background-color: #ffffff;
				border-radius: 8px;
				box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			h1 {
				color: #333333;
			}
			p {
				color: #555555;
			}
			a {
				display: inline-block;
				padding: 10px 20px;
				background-color: #007bff;
				text-decoration: none;
				border-radius: 5px;
			}
			span {
				color: #ffffff;
			}
			a:hover {
				background-color: #0056b3;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Email Verification</h1>
			<p>Please verify your email by clicking the following link:</p>
			<p><a href="https://vokki.net%s/?token=%s"><span>Verify Email<span></a></p>
		</div>
	</body>
	</html>
`, vokki_constants.RouteVerifyEmail, token)

	msg := []byte(subject + headers + "\n" + body)

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

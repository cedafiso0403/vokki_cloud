package services

import (
	"bytes"
	"html/template"
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

	// Parse the HTML template file
	tmpl, err := template.ParseFiles("internal/views/email_verification.html")
	if err != nil {
		log.Println("Error parsing template: ", err)
		return err
	}

	// Create a buffer to hold the executed template
	var body bytes.Buffer

	// Data to pass to the template
	data := struct {
		Route string
		Token string
	}{
		Route: vokki_constants.RouteVerifyEmail,
		Token: token,
	}

	// Execute the template and store it in the buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Println("Error executing template: ", err)
		return err
	}

	// Combine subject, headers, and body into the final message
	msg := []byte(subject + headers + "\n" + body.String())

	// Send the email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, msg)
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

	subject := "Subject: Create New Password\n"

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"Priority: Urgent\n" +
		"X-Priority: 1\n" +
		"Importance: High\n"

	// Parse the HTML template file
	tmpl, err := template.ParseFiles("internal/views/new_password.html")
	if err != nil {
		log.Println("Error parsing template: ", err)
		return err
	}

	// Create a buffer to hold the executed template
	var body bytes.Buffer

	// Data to pass to the template
	data := struct {
		Route string
		Token string
	}{
		Route: vokki_constants.RouteCreateNewPassword,
		Token: token,
	}

	// Execute the template and store it in the buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Println("Error executing template: ", err)
		return err
	}

	// Combine subject, headers, and body into the final message
	msg := []byte(subject + headers + "\n" + body.String())

	// Send the email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, msg)
	if err != nil {
		log.Println("Error sending email: ", err)
		return err
	}

	return nil
}

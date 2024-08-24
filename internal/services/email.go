package services

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"vokki_cloud/internal/constants"
	"vokki_cloud/internal/models"
)

// Common data structure for email templates
type EmailData struct {
	Route      string
	Token      string
	BaseScheme string
	BaseHost   string
	BasePath   string
}

// Helper function to populate shared data
func getEmailData(route, token string) EmailData {
	return EmailData{
		Route:      route,
		Token:      token,
		BaseScheme: vokki_constants.BaseScheme,
		BaseHost:   vokki_constants.BaseHost,
		BasePath:   vokki_constants.BasePath,
	}
}

// Common function to send emails
func sendEmail(to string, subject string, tmplPath string, data EmailData) error {
	from := os.Getenv("FROM_EMAIL")
	appPassword := os.Getenv("FROM_EMAIL_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, appPassword, smtpHost)
	toList := []string{to}

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"Priority: Urgent\n" +
		"X-Priority: 1\n" +
		"Importance: High\n"

	// Parse the HTML template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println("Error parsing template: ", err)
		return err
	}

	// Create a buffer to hold the executed template
	var body bytes.Buffer

	// Execute the template and store it in the buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Println("Error executing template: ", err)
		return err
	}

	// Combine subject, headers, and body into the final message
	msg := []byte("Subject: " + subject + "\n" + headers + "\n" + body.String())

	// Send the email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, msg)
	if err != nil {
		log.Println("Error sending email: ", err)
		return err
	}

	return nil
}

// SendVerificationEmail sends a verification email to the user
func SendVerificationEmail(user models.User, token string) error {
	data := getEmailData(vokki_constants.RouteVerifyEmail, token)
	return sendEmail(user.Email, "Email Verification", "internal/views/email_verification.html", data)
}

// SendPasswordResetEmail sends a password reset email to the user
func SendPasswordResetEmail(user models.User, token string) error {
	data := getEmailData(vokki_constants.RouteCreateNewPassword, token)
	return sendEmail(user.Email, "Create New Password", "internal/views/new_password.html", data)
}
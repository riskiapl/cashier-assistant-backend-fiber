package utils

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"time"
)

// Mailer represents the email sending service
type Mailer struct {
	Host     string
	Port     string
	Username string
	Password string
	FromName string
}

// NewMailer creates a new mailer instance using environment variables
func NewMailer() (*Mailer, error) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Validate SMTP configuration
	if host == "" || port == "" || username == "" || password == "" {
		return nil, errors.New("SMTP configuration is incomplete")
	}

	return &Mailer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		FromName: "Cashierly", // Set default sender name
	}, nil
}

// SendEmail is a generic method to send any email
func (m *Mailer) SendEmail(to string, subject string, htmlBody string) error {
	// Set up authentication
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	// Add debug log
	fmt.Printf("SMTP Config: Host=%s, Port=%s, Username=%s\n", m.Host, m.Port, m.Username)

	// Format the from address with display name: "Sender Name <email@example.com>"
	fromAddress := fmt.Sprintf("%s <%s>", m.FromName, m.Username)

	// Format the email with HTML content type
	message := fmt.Appendf(nil, "From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n", fromAddress, to, subject, htmlBody)

	// Send email
	err := smtp.SendMail(m.Host+":"+m.Port, auth, m.Username, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}

// SendOTPEmail sends an email containing the OTP code
func (m *Mailer) SendOTPEmail(email string, otp string, expiresInMinutes int) error {
	subject := "Your OTP Code"
	expiresAt := fmt.Sprintf("%d minutes", expiresInMinutes)
	htmlBody := OTPMail(otp, expiresAt)

	return m.SendEmail(email, subject, htmlBody)
}

// SendResetPasswordEmail sends an email with password reset link
func (m *Mailer) SendResetPasswordEmail(email string, resetLink string, expiresInMinutes int) error {
	subject := "Reset Your Password"
	expiresAt := fmt.Sprintf("%d minutes", expiresInMinutes)
	htmlBody := ResetPasswordMail(resetLink, expiresAt)

	return m.SendEmail(email, subject, htmlBody)
}

// GetExpirationTime returns a future time based on minutes
func GetExpirationTime(minutes int) time.Time {
	return time.Now().Add(time.Duration(minutes) * time.Minute)
}

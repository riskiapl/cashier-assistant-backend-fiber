package services

import (
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/riskiapl/fiber-app/database"
	"github.com/riskiapl/fiber-app/models"
	"github.com/riskiapl/fiber-app/repository"
	"github.com/riskiapl/fiber-app/types"
	"github.com/riskiapl/fiber-app/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		authRepo: repository.NewAuthRepository(database.DB),
	}
}

func (s *AuthService) Login(input types.LoginInput) (*types.LoginResponse, error) {
	// Cari member berdasarkan email
	member, err := s.authRepo.GetMemberByUserOrMail(input.Userormail)
	if err != nil {
		return nil, err
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	// Buat response
	response := &types.LoginResponse{
		ID:       member.ID,
		Username: member.Username,
		Email:    member.Email,
		Status:   member.Status,
		Avatar:   member.Avatar,
	}

	return response, nil
}

func (s *AuthService) Register(input types.RegisterInput) (*types.RegisterResponse, error) {
	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Buat pending member baru
	pendingMember := &models.PendingMember{
		Username:      input.Username,
		Email:         input.Email,
		Password:      string(hashedPassword),
		PlainPassword: input.Password,
	}

	// Simpan pending member ke database
	if err := s.authRepo.Register(pendingMember); err != nil {
		return nil, err
	}

	// Generate OTP
	otpCode := s.GenerateOTP()

	// Get expiration time
	expirationTime := time.Now().Add(15 * time.Minute)

	// Store OTP in database
	otp := &models.OTP{
		Email:      input.Email,
		OtpCode:    otpCode,
		IsVerified: false,
		ExpiredAt:  expirationTime, // OTP expires after 15 minutes
		ActionType: "I",            // I for initial registration
	}

	if err := s.authRepo.StoreOTP(otp); err != nil {
		// Clean up by deleting the pending member if OTP storage fails
		if deleteErr := s.authRepo.DeletePendingMember(input.Email); deleteErr != nil {
			// Log the error but continue
			fmt.Printf("Failed to delete pending member after OTP storage failure: %v\n", deleteErr)
		}
		return nil, err
	}

	// Send OTP via email
	if err := s.SendOTPEmail(input.Email, otpCode); err != nil {
		// Clean up by deleting both pending member and OTP if email sending fails
		if deleteErr := s.authRepo.DeleteOTP(input.Email); deleteErr != nil {
			fmt.Printf("Failed to delete OTP after email sending failure: %v\n", deleteErr)
		}
		if deleteErr := s.authRepo.DeletePendingMember(input.Email); deleteErr != nil {
			fmt.Printf("Failed to delete pending member after email sending failure: %v\n", deleteErr)
		}
		return nil, err
	}

	// Buat response
	response := &types.RegisterResponse{
		Success:   "OTP sent to your email",
		ExpiredAt: expirationTime,
	}

	return response, nil
}

func (s *AuthService) VerifyRegistration(input types.VerifyOTPInput) error {
	// Verify OTP
	_, err := s.authRepo.VerifyOTP(input.Email, input.OtpCode)
	if err != nil {
		return err
	}

	// Complete registration process
	if err := s.authRepo.CompleteRegistration(input.Email); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) IsUsernameExists(username string) (bool, error) {
	// Check if username already exists in database
	exists, err := s.authRepo.IsUsernameExists(username)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *AuthService) GenerateOTP() string {
	// Create a local random generator with a time-based source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random 6-digit OTP
	otp := r.Intn(900000) + 100000

	return fmt.Sprintf("%06d", otp)
}

func (s *AuthService) SendOTPEmail(email string, otp string) error {
	// Get email configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Validate SMTP configuration
	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" {
		return errors.New("SMTP configuration is incomplete")
	}

	// Set up authentication
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Add debug log
	fmt.Printf("SMTP Config: Host=%s, Port=%s, Username=%s\n", smtpHost, smtpPort, smtpUsername)

	// Compose the email
	from := smtpUsername
	to := []string{email}

	subject := "Your OTP Code"
	expiresAt := "15 minutes"
	htmlBody := utils.OTPMail(otp, expiresAt)

	// Format the email with HTML content type
	message := fmt.Appendf(nil, "To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n", email, subject, htmlBody)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) DeletePendingMember(email string) error {
	return s.authRepo.DeletePendingMember(email)
}

func (s *AuthService) ResendOTP(email string) (*types.RegisterResponse, error) {
	// Check if pending member exists
	_, err := s.authRepo.GetPendingMemberByEmail(email)
	if err != nil {
		return nil, errors.New("no pending member found, try registering again")
	}

	// Generate OTP
	otpCode := s.GenerateOTP()

	// Get expiration time
	expirationTime := time.Now().Add(15 * time.Minute)

	// Store OTP in database
	otp := &models.OTP{
		Email:      email,
		OtpCode:    otpCode,
		IsVerified: false,
		ExpiredAt:  expirationTime, // OTP expires after 15 minutes
		ActionType: "U",            // U for update
	}

	if err := s.authRepo.StoreOTP(otp); err != nil {
		return nil, err
	}

	// Send OTP via email
	if err := s.SendOTPEmail(email, otpCode); err != nil {
		// Clean up by deleting the OTP if email sending fails
		if deleteErr := s.authRepo.DeleteOTP(email); deleteErr != nil {
			fmt.Printf("Failed to delete OTP after email sending failure: %v\n", deleteErr)
		}
		return nil, err
	}

	// Buat response
	response := &types.RegisterResponse{
		Success:   "OTP resent to your email",
		ExpiredAt: expirationTime,
	}

	return response, nil
}

func (s *AuthService) ResetPassword(input types.ForgotPasswordInput) (*types.RegisterResponse, error) {
	// Check if member exists
	member, err := s.authRepo.GetMemberByUserOrMail(input.Email)
	if err != nil {
		return nil, err
	}

	// Get expiration time
	expirationTime := time.Now().Add(15 * time.Minute)

	// Prepare data for token generation and database storage
	tokenData := struct {
		Email   string
		Token   string
		Expired time.Time
	}{
		Email:   member.Email,
		Token:   utils.GenerateUUID(),
		Expired: expirationTime,
	}

	// Generate JWT token
	tokenString, err := utils.GenerateToken(map[string]any{
		"email":   tokenData.Email,
		"token":   tokenData.Token,
		"expired": tokenData.Expired.Unix(),
	})
	if err != nil {
		return nil, err
	}

	// Store token in database using the same data
	resetToken := &models.ResetPasswordToken{
		Email:   tokenData.Email,
		Token:   tokenData.Token,
		Expired: tokenData.Expired,
		IsUsed:  false,
	}

	if err := s.authRepo.StoreResetPasswordToken(resetToken); err != nil {
		return nil, err
	}

	// Get frontend URL from env
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		return nil, errors.New("frontend URL is not configured")
	}

	// Create reset password link
	resetLink := fmt.Sprintf("%s/auth/reset-password?email=%s&token=%s", frontendURL, input.Email, tokenString)

	// Send reset password email
	if err := s.SendResetPasswordEmail(member.Email, resetLink); err != nil {
		// Clean up by deleting the token if email sending fails
		if deleteErr := s.authRepo.DeleteResetPasswordToken(member.Email); deleteErr != nil {
			fmt.Printf("Failed to delete reset token after email sending failure: %v\n", deleteErr)
		}
		return nil, err
	}

	// Create response
	response := &types.RegisterResponse{
		Success:   "Reset password link sent to your email",
		ExpiredAt: expirationTime,
	}

	return response, nil
}

func (s *AuthService) SendResetPasswordEmail(email string, resetLink string) error {
	// Get email configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Validate SMTP configuration
	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" {
		return errors.New("SMTP configuration is incomplete")
	}

	// Set up authentication
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Compose the email
	from := smtpUsername
	to := []string{email}

	subject := "Reset Your Password"
	expiresAt := "15 minutes"
	htmlBody := utils.ResetPasswordMail(resetLink, expiresAt)

	// Format the email with HTML content type
	message := fmt.Appendf(nil, "To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n", email, subject, htmlBody)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ChangePassword(resetPasswordData types.ResetPasswordData, newPassword string) error {
	// Get the reset token from database to verify
	_, err := s.authRepo.GetResetPasswordToken(resetPasswordData.Email, resetPasswordData.Token)
	if err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the member's password
	if err := s.authRepo.UpdatePassword(resetPasswordData.Email, string(hashedPassword), newPassword); err != nil {
		return err
	}

	// Mark the token as used
	if err := s.authRepo.MarkResetTokenAsUsed(resetPasswordData.Email); err != nil {
		return err
	}

	return nil
}

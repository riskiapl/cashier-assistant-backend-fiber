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

	// Buat member baru
	member := &models.Member{
		Username:      input.Username,
		Email:         input.Email,
		Password:      string(hashedPassword),
		PlainPassword: input.PlainPassword,
	}

	// Simpan member ke database
	if err := s.authRepo.Register(member); err != nil {
		return nil, err
	}

	// Buat response
	response := &types.RegisterResponse{
		Message: "Registration successful",
	}

	return response, nil
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

	// Set up authentication
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

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

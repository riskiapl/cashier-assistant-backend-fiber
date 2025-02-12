package services

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type AuthService struct {
	// Tambahkan dependencies yang diperlukan seperti DB
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(input LoginInput) (*fiber.Map, error) {
	// TODO: Implement actual DB lookup
	// Contoh user hardcoded untuk demo
	user := &User{
		ID:       1,
		Email:    "john@example.com",
		Password: "$2a$10$YourHashedPasswordHere", // Hashed password
		Name:     "John Doe",
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Sign token
	t, err := token.SignedString([]byte("your_jwt_secret")) // Gunakan env variable untuk secret
	if err != nil {
		return nil, err
	}

	return &fiber.Map{
		"token": t,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	}, nil
}

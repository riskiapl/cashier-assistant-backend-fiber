package services

import (
	"errors"

	"github.com/riskiapl/fiber-app/database"
	"github.com/riskiapl/fiber-app/models"
	"github.com/riskiapl/fiber-app/repository"
	"github.com/riskiapl/fiber-app/types"
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
	member, err := s.authRepo.GetMemberByEmail(input.Userormail)
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

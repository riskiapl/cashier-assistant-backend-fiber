package services

import (
	"errors"

	"github.com/riskiapl/fiber-app/database"
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
	}

	return response, nil
}

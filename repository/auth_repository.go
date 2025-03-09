package repository

import (
	"errors"

	"github.com/riskiapl/fiber-app/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(DB *gorm.DB) *AuthRepository {
	if DB == nil {
		panic("database connection is not initialized")
	}
	return &AuthRepository{DB: DB}
}

func (r *AuthRepository) GetMemberByUserOrMail(userormail string) (*models.Member, error) {
	var member models.Member
	result := r.DB.Where("email = ? OR username = ?", userormail, userormail).First(&member)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("member not found")
		}
		return nil, result.Error
	}

	return &member, nil
}

func (r *AuthRepository) Register(member *models.Member) error {
	// Cek apakah email sudah terdaftar
	var existingMember models.Member
	result := r.DB.Where("email = ?", member.Email).First(&existingMember)
	if result.Error == nil {
		return errors.New("email already registered")
	}

	// Cek apakah username sudah terdaftar
	result = r.DB.Where("username = ?", member.Username).First(&existingMember)
	if result.Error == nil {
		return errors.New("username already taken")
	}

	// Simpan member baru ke database
	member.Status = "member"
	member.ActionType = "I"
	if err := r.DB.Create(member).Error; err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) IsUsernameExists(username string) (bool, error) {
	var member models.Member
	result := r.DB.Where("username = ?", username).First(&member)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

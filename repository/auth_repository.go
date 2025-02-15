package repository

import (
	"errors"
	"log"

	"github.com/riskiapl/fiber-app/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(DB *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: DB}
}

func (r *AuthRepository) GetMemberByEmail(userormail string) (*models.Member, error) {
	log.Println(userormail, "masuk userormail")
	if r.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

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

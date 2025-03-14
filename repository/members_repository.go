package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/riskiapl/fiber-app/models"
)

type MemberRepository struct {
	DB *gorm.DB
}

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{DB: db}
}

func (r *MemberRepository) GetMembers(limit, offset int) ([]models.Member, int64, error) {
	var members []models.Member
	var count int64

	if err := r.DB.Model(&models.Member{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Limit(limit).Offset(offset).Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, count, nil
}

func (r *MemberRepository) GetMemberByID(id uint) (*models.Member, error) {
	var member models.Member
	if err := r.DB.First(&member, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("member not found")
		}
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepository) UpdateMember(member *models.Member) error {
	return r.DB.Save(member).Error
}

func (r *MemberRepository) DeleteMember(id uint) error {
	return r.DB.Delete(&models.Member{}, id).Error
}

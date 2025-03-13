package repository

import (
	"errors"
	"time"

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

func (r *AuthRepository) Register(pendingMember *models.PendingMember) error {
	// Cek apakah email sudah terdaftar di members
	// Cek apakah email atau username sudah terdaftar di members
	var existingMember models.Member
	result := r.DB.Where("email = ? OR username = ?", pendingMember.Email, pendingMember.Username).First(&existingMember)
	if result.Error == nil {
		if existingMember.Email == pendingMember.Email {
			return errors.New("email already registered")
		}
		return errors.New("username already taken")
	}

	// Cek apakah email atau username sudah terdaftar di pending_members
	var existingPending models.PendingMember
	result = r.DB.Where("email = ? OR username = ?", pendingMember.Email, pendingMember.Username).First(&existingPending)
	if result.Error == nil {
		if existingPending.Email == pendingMember.Email {
			return errors.New("email already in registration process")
		}
		return errors.New("username already in registration process")
	}

	// Simpan pending member baru ke database
	pendingMember.ActionType = "I"
	if err := r.DB.Create(pendingMember).Error; err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) StoreOTP(otp *models.OTP) error {
	// Delete existing OTP for this email if any
	if err := r.DB.Where("email = ?", otp.Email).Delete(&models.OTP{}).Error; err != nil {
		return err
	}

	// Save the new OTP
	if err := r.DB.Create(otp).Error; err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) VerifyOTP(email string, otpCode string) (*models.OTP, error) {
	var otp models.OTP

	result := r.DB.Where("email = ? AND otp_code = ? AND is_verified = ? AND expired_at > ?",
		email, otpCode, false, time.Now()).First(&otp)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired OTP")
		}
		return nil, result.Error
	}

	return &otp, nil
}

func (r *AuthRepository) CompleteRegistration(email string) error {
	var pendingMember models.PendingMember

	// Find the pending member
	if err := r.DB.Where("email = ?", email).First(&pendingMember).Error; err != nil {
		return err
	}

	// Begin transaction
	tx := r.DB.Begin()

	// Create a new member from pending data
	member := models.Member{
		Username:      pendingMember.Username,
		Email:         pendingMember.Email,
		Password:      pendingMember.Password,
		PlainPassword: pendingMember.PlainPassword,
		Status:        "member",
		Avatar:        "", // Default avatar
		ActionType:    "I",
	}

	if err := tx.Create(&member).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the pending member
	if err := tx.Delete(&pendingMember).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Mark OTP as verified
	if err := tx.Model(&models.OTP{}).Where("email = ?", email).
		Updates(map[string]interface{}{"is_verified": true}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) IsUsernameExists(username string) (bool, error) {
	// Check in members table
	var member models.Member
	result := r.DB.Where("username = ?", username).First(&member)
	if result.Error == nil {
		return true, nil
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, result.Error
	}

	// Check in pending_members table
	var pendingMember models.PendingMember
	result = r.DB.Where("username = ?", username).First(&pendingMember)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

func (r *AuthRepository) DeletePendingMember(email string) error {
	result := r.DB.Where("email = ?", email).Delete(&models.PendingMember{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("pending member not found")
	}

	return nil
}

func (r *AuthRepository) DeleteOTP(email string) error {
	result := r.DB.Where("email = ?", email).Delete(&models.OTP{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("OTP not found")
	}

	return nil
}

func (r *AuthRepository) DeleteExpiredPendingMembers(expirationDuration time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-expirationDuration)
	result := r.DB.Where("created_at < ?", cutoffTime).Delete(&models.PendingMember{})
	return result.RowsAffected, result.Error
}

func (r *AuthRepository) DeleteExpiredOTPs(expirationDuration time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-expirationDuration)
	result := r.DB.Where("created_at < ?", cutoffTime).Delete(&models.OTP{})
	return result.RowsAffected, result.Error
}

func (r *AuthRepository) GetPendingMemberByEmail(email string) (*models.PendingMember, error) {
	var pendingMember models.PendingMember
	result := r.DB.Where("email = ?", email).First(&pendingMember)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("pending member not found")
		}
		return nil, result.Error
	}

	return &pendingMember, nil
}

package models

import (
	"time"
)

type OTP struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	OtpCode   string         `gorm:"type:varchar(5);not null" json:"otp_code"` 
	IsVerified    bool       `gorm:"default:false;index" json:"is_verified"`
	ExpiredAt time.Time      `gorm:"not null" json:"expired_at"`
	ActionType string        `gorm:"type:varchar(100);not null;index" json:"action_type"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
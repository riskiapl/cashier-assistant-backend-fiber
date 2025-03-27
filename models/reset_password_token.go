package models

import (
	"time"
)

// ResetPasswordToken represents a token for resetting a user's password
type ResetPasswordToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Token     string    `gorm:"size:255;not null" json:"token"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	Expired   time.Time `gorm:"not null" json:"expired"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// No need for BeforeCreate hook since GORM handles auto-increment automatically

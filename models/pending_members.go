package models

import (
	"time"
)

type PendingMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"type:varchar(100);not null;unique" json:"username"`
	Email     string         `gorm:"type:varchar(100);not null;unique" json:"email"`
	Password  string         `gorm:"type:varchar(100);not null" json:"password"`
	PlainPassword string 		 `gorm:"type:varchar(100);not null" json:"plain_password"`
	ActionType string        `gorm:"type:varchar(100);not null;index" json:"action_type"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

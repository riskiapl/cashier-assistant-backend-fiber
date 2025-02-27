package models

import (
	"time"
)

type Member struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string     `gorm:"type:varchar(100);unique;not null" json:"username"`
	Email     string     `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string     `gorm:"type:varchar(100);not null" json:"password"`
	PlainPassword string `gorm:"type:varchar(100);not null" json:"plain_password"`
	Status string        `gorm:"type:varchar(100);not null;index" json:"status"`
	Avatar string        `gorm:"type:varchar(100);not null" json:"avatar"`
	ActionType string    `gorm:"type:varchar(100);not null;index" json:"action_type"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
package models

import (
	"gorm.io/gorm"
	"time"
)

type Token struct {
	gorm.Model
	UserId    int       `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	IssuedAt  time.Time `json:"issued_at" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
}

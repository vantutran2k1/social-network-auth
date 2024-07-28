package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username  string    `json:"username" gorm:"unique, not null"`
	Password  string    `json:"password" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique, not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

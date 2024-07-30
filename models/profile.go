package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"not null"`
	FirstName   string `json:"first_name" gorm:"not null"`
	LastName    string `json:"last_name" gorm:"not null"`
	DateOfBirth string `json:"date_of_birth" gorm:"not null"`
	Address     string `json:"address" gorm:"not null"`
	Phone       string `json:"phone" gorm:"not null"`
}

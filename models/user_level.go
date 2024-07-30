package models

import "gorm.io/gorm"

type UserLevel struct {
	gorm.Model
	UserID  uint `json:"user_id" gorm:"not null"`
	LevelID uint `json:"level_id" gorm:"not null"`
}

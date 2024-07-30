package models

import "gorm.io/gorm"

type Level struct {
	gorm.Model
	LevelName string `json:"level_name" gorm:"not null;unique"`
}

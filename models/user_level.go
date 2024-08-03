package models

import "gorm.io/gorm"

type UserLevel struct {
	gorm.Model
	UserID  uint `json:"user_id" gorm:"not null"`
	LevelID uint `json:"level_id" gorm:"not null"`
}

func (ul *UserLevel) Save(db *gorm.DB) error {
	return db.Create(&ul).Error
}

func (ul *UserLevel) Get(db *gorm.DB) (UserLevel, error) {
	dbUserLevel := UserLevel{}
	db.Where("user_id = ? AND level_id = ?", ul.UserID, ul.LevelID).First(&dbUserLevel)

	return dbUserLevel, db.Error
}

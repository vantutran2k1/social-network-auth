package models

import (
	"gorm.io/gorm"
	"time"
)

type Token struct {
	ID        uint      `gorm:"primarykey"`
	UserId    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	IssuedAt  time.Time `json:"issued_at" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
}

func (token *Token) Save(db *gorm.DB) error {
	return db.Create(token).Error
}

func (token *Token) Validate(db *gorm.DB, tokenString string) bool {
	dbToken := Token{}
	err := db.Where("token = ?", tokenString).First(&dbToken).Error
	if err != nil {
		return false
	}

	return dbToken.ExpiresAt.After(time.Now())
}

func (token *Token) Revoke(db *gorm.DB, tokenString string) error {
	dbToken := Token{}
	err := db.Where("token = ?", tokenString).First(&dbToken).Error
	if err != nil {
		return err
	}

	dbToken.ExpiresAt = time.Now()
	return db.Save(&dbToken).Error
}

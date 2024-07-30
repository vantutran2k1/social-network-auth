package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Token struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	IssuedAt  time.Time `json:"issued_at" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
}

func (token *Token) Save(db *gorm.DB) error {
	return db.Create(token).Error
}

func (token *Token) Validate(db *gorm.DB, tokenString string) bool {
	dbToken := Token{}
	db.Where("token = ?", tokenString).First(&dbToken)
	if dbToken.ID == 0 {
		return false
	}

	return dbToken.ExpiresAt.After(time.Now().UTC())
}

func (token *Token) Revoke(db *gorm.DB, tokenString string) error {
	dbToken := Token{}
	db.Where("token = ?", tokenString).First(&dbToken)
	if dbToken.ID == 0 {
		return errors.New("token not found")
	}

	if !dbToken.ExpiresAt.After(time.Now().UTC()) {
		return errors.New("token is expired")
	}

	dbToken.ExpiresAt = time.Now().UTC()
	return db.Save(&dbToken).Error
}

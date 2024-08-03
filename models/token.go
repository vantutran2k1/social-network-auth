package models

import (
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
	var dbToken Token
	err := db.Where(&Token{Token: tokenString}).First(&dbToken).Error
	if err != nil {
		return false
	}

	return !dbToken.isExpired()
}

func (token *Token) Revoke(db *gorm.DB, tokenString string) error {
	var dbToken Token
	err := db.Where(&Token{Token: tokenString}).First(&dbToken).Error
	if err != nil {
		return err
	}

	dbToken.ExpiresAt = time.Now().UTC()

	return db.Save(&dbToken).Error
}

func (token *Token) isExpired() bool {
	return time.Now().UTC().After(token.ExpiresAt)
}

package models

import (
	"time"

	"gorm.io/gorm"
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

func (token *Token) RevokeUserActiveTokens(db *gorm.DB, userID uint) error {
	activeTokens, err := token.getActiveTokensByUser(db, userID)
	if err != nil {
		return err
	}

	for _, t := range activeTokens {
		t.ExpiresAt = time.Now().UTC()
	}

	return db.Save(activeTokens).Error
}

func (token *Token) getActiveTokensByUser(db *gorm.DB, userID uint) ([]*Token, error) {
	var tokens []*Token
	if err := db.Where(&Token{UserID: userID}).Find(&tokens).Error; err != nil {
		return nil, err
	}

	activeTokens := make([]*Token, 0, len(tokens))
	for _, t := range tokens {
		if !t.isExpired() {
			activeTokens = append(activeTokens, t)
		}
	}

	return activeTokens, nil
}

func (token *Token) isExpired() bool {
	return time.Now().UTC().After(token.ExpiresAt)
}

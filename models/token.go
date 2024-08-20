package models

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/vantutran2k1/social-network-auth/utils"
	"gorm.io/gorm"
)

type Token struct {
	ID        uuid.UUID `gorm:"primarykey"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	IssuedAt  time.Time `json:"issued_at" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
}

func (token *Token) CreateLoginToken(db *gorm.DB, userID uuid.UUID) (*Token, error) {
	expirationAfter, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if err != nil {
		return nil, err
	}

	expirationTime := time.Now().UTC().Add(time.Duration(expirationAfter) * time.Minute)
	claims := &utils.Claims{
		UserID:         userID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenJwt.SignedString(utils.JwtKey)
	if err != nil {
		return nil, err
	}

	t := Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     tokenString,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: expirationTime,
	}
	if err := db.Create(&t).Error; err != nil {
		return nil, err
	}

	return &t, nil
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

func (token *Token) RevokeUserActiveTokens(db *gorm.DB, userID uuid.UUID) error {
	activeTokens, err := token.getActiveTokensByUser(db, userID)
	if err != nil {
		return err
	}

	for _, t := range activeTokens {
		t.ExpiresAt = time.Now().UTC()
	}

	return db.Save(activeTokens).Error
}

func (token *Token) getActiveTokensByUser(db *gorm.DB, userID uuid.UUID) ([]*Token, error) {
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

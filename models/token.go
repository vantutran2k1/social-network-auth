package models

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/vantutran2k1/social-network-auth/errors"
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

func (token *Token) CreateLoginToken(db *gorm.DB, userID uuid.UUID) (*Token, *errors.ApiError) {
	expirationAfter, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	expirationTime := time.Now().UTC().Add(time.Duration(expirationAfter) * time.Minute)
	claims := &utils.Claims{
		UserID:         userID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenJwt.SignedString(utils.JwtKey)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	t := Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     tokenString,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: expirationTime,
	}
	if err := db.Create(&t).Error; err != nil {
		return nil, errors.InternalServerError(err.Error())
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

func (token *Token) Revoke(db *gorm.DB, tokenString string) *errors.ApiError {
	var dbToken Token
	err := db.Where(&Token{Token: tokenString}).First(&dbToken).Error
	if err != nil {
		if utils.IsRecordNotFound(err) {
			return errors.BadRequestError("token not found")
		}

		return errors.InternalServerError(err.Error())
	}

	dbToken.ExpiresAt = time.Now().UTC()

	if err := db.Save(&dbToken).Error; err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (token *Token) RevokeUserActiveTokens(db *gorm.DB, userID uuid.UUID) *errors.ApiError {
	activeTokens, err := token.getActiveTokensByUser(db, userID)
	if err != nil {
		return err
	}

	if len(activeTokens) == 0 {
		return nil
	}

	for _, t := range activeTokens {
		t.ExpiresAt = time.Now().UTC()
	}

	if err := db.Save(activeTokens).Error; err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (token *Token) getActiveTokensByUser(db *gorm.DB, userID uuid.UUID) ([]*Token, *errors.ApiError) {
	var tokens []*Token
	if err := db.Where(&Token{UserID: userID}).Find(&tokens).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return []*Token{}, nil
		}

		return nil, errors.InternalServerError(err.Error())
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

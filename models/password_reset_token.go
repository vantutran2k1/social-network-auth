package models

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1/social-network-auth/errors"
	"github.com/vantutran2k1/social-network-auth/utils"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID          uuid.UUID `gorm:"primary_key"`
	UserID      uuid.UUID `gorm:"not null"`
	Token       string    `gorm:"not null"`
	TokenExpiry time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null,autoCreateTime:false"`
}

func (t *PasswordResetToken) CreateResetToken(db *gorm.DB, userID uuid.UUID) (*PasswordResetToken, *errors.ApiError) {
	if err := db.Where(&User{ID: userID}).First(&User{}).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return nil, errors.BadRequestError("user not found")
		}

		return nil, errors.InternalServerError(err.Error())
	}

	token, err := generateResetToken()
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	expirationAfter, err := strconv.Atoi(os.Getenv("RESET_PASSWORD_TOKEN_EXPIRATION_MINUTES"))
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	resetToken := &PasswordResetToken{
		ID:          uuid.New(),
		UserID:      userID,
		Token:       token,
		TokenExpiry: time.Now().UTC().Add(time.Duration(expirationAfter) * time.Minute),
		CreatedAt:   time.Now().UTC(),
	}

	if err := db.Create(&resetToken).Error; err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return resetToken, nil
}

func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

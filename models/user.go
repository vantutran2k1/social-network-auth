package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1/social-network-auth/transaction"
	"github.com/vantutran2k1/social-network-auth/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Level     Level     `json:"level" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime:false"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;autoUpdateTime:false"`
	DeletedAt gorm.DeletedAt
}

func (user *User) Register(
	db *gorm.DB,
	username string,
	password string,
	email string,
) error {
	err := db.Where(&User{Username: username}).Or(&User{Email: email}).First(&User{}).Error
	if err == nil {
		return errors.New("username or email already exists")
	}

	if !utils.IsRecordNotFound(err) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	user.ID = uuid.New()
	user.Level = BRONZE
	user.Username = username
	user.Email = email

	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	return db.Create(&user).Error
}

func (user *User) Authenticate(db *gorm.DB, username string, password string) (*User, error) {
	var dbUser User
	if err := db.Where(&User{Username: username}).First(&dbUser).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return nil, fmt.Errorf("user %s not found", username)
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &dbUser, nil
}

func (user *User) UpdateLevel(db *gorm.DB, userID uuid.UUID, levelName string) error {
	err := db.Where(&User{ID: userID}).First(&user).Error
	if err != nil {
		if utils.IsRecordNotFound(err) {
			return fmt.Errorf("user %v not found", userID)
		}
		return err
	}

	level := GetLevelFromName(levelName)
	if user.Level == level {
		return nil
	}
	user.Level = level

	user.UpdatedAt = time.Now().UTC()

	return db.Save(user).Error
}

func (user *User) UpdatePassword(db *gorm.DB, currentPassword string, newPassword string) error {
	if err := db.Where(&User{ID: user.ID}).First(&user).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return errors.New("user not found")
		}

		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("invalid password")
	}

	if currentPassword == newPassword {
		return errors.New("new password can not be the same as current one")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now().UTC()

	err = transaction.TxManager.WithTransaction(func(tx *gorm.DB) error {
		token := &Token{}
		if err := token.RevokeUserActiveTokens(db, user.ID); err != nil {
			return err
		}

		return db.Save(user).Error
	})

	return err
}

func (user *User) GetUserByUsernameOrEmail(db *gorm.DB, userIdentity string) (*User, error) {
	var dbUser User
	if err := db.Where(&User{Username: userIdentity}).Or(&User{Email: userIdentity}).First(&dbUser).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return &dbUser, nil
}

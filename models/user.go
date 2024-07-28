package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username  string    `json:"username" gorm:"unique, not null"`
	Password  string    `json:"-"`
	Email     string    `json:"email" gorm:"unique, not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

func (user *User) Register(db *gorm.DB) error {
	err := db.Where("username = ?", user.Username).First(&User{}).Error
	if err == nil {
		return errors.New("username already exists")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)
		return db.Create(&user).Error
	}

	return err
}

func (user *User) Authenticate(db *gorm.DB, password string) bool {
	err := db.Where("username = ?", user.Username).First(&user).Error
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

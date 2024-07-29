package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/vantutran2k1/social-network-auth/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique" binding:"required"`
	Password string `json:"-" binding:"required"`
	Email    string `json:"email" gorm:"unique" binding:"required"`
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

func (user *User) GenerateToken() (string, time.Time, error) {
	expirationAfter, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if err != nil {
		return "", time.Time{}, err
	}

	expirationTime := time.Now().UTC().Add(time.Duration(expirationAfter) * time.Minute)
	claims := &utils.Claims{
		Username:       user.Username,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(utils.JwtKey)

	return tokenString, expirationTime, err
}

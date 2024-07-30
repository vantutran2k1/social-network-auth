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
	var dbUser User
	db.Where("username = ? OR email = ?", user.Username, user.Email).First(&dbUser)
	if dbUser.ID != 0 {
		return errors.New("username or email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return db.Create(&user).Error
}

func (user *User) Authenticate(db *gorm.DB, password string) bool {
	db.Where("username = ?", user.Username).First(&user)
	if user.ID == 0 {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) GenerateToken() (string, time.Time, error) {
	expirationAfter, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if err != nil {
		return "", time.Time{}, err
	}

	expirationTime := time.Now().UTC().Add(time.Duration(expirationAfter) * time.Minute)
	claims := &utils.Claims{
		UserID:         user.ID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(utils.JwtKey)

	return tokenString, expirationTime, err
}

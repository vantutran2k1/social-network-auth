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
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique" binding:"required"`
	Password  string    `json:"-" binding:"required"`
	Email     string    `json:"email" gorm:"unique" binding:"required"`
	Level     Level     `json:"level" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`
	DeletedAt gorm.DeletedAt
}

func (user *User) Register(db *gorm.DB) error {
	var dbUser User

	err := db.Where(&User{Username: user.Username}).Or(&User{Email: user.Email}).First(&dbUser).Error
	if err == nil {
		return errors.New("username or email already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	dbUser.Password = string(hashedPassword)

	dbUser.Level = BRONZE
	dbUser.Username = user.Username
	dbUser.Email = user.Email

	return db.Create(&dbUser).Error

}

func (user *User) Authenticate(db *gorm.DB, password string) bool {
	err := db.Where(&User{Username: user.Username}).First(&user).Error
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) GenerateToken(db *gorm.DB) (string, error) {
	expirationAfter, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().UTC().Add(time.Duration(expirationAfter) * time.Minute)
	claims := &utils.Claims{
		UserID:         user.ID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenJwt.SignedString(utils.JwtKey)

	token := Token{
		UserID:    user.ID,
		Token:     tokenString,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: expirationTime,
	}
	err = token.Save(db)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (user *User) AssignLevel(db *gorm.DB, levelName string) error {
	var dbUser User
	err := db.Where(&User{ID: user.ID}).First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	level := GetLevelFromName(levelName)
	if dbUser.Level == level {
		return nil
	}
	dbUser.Level = level

	return db.Save(dbUser).Error
}

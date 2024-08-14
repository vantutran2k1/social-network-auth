package models

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vantutran2k1/social-network-auth/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Level     Level     `json:"level" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime:false"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;autoUpdateTime:false"`
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

	dbUser.CreatedAt = time.Now().UTC()
	dbUser.UpdatedAt = time.Now().UTC()

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

	dbUser.UpdatedAt = time.Now().UTC()

	return db.Save(dbUser).Error
}

func (user *User) UpdatePassword(db *gorm.DB, currentPassword string, newPassword string) error {
	err := db.Where(&User{ID: user.ID}).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}

		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
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

	return db.Save(user).Error
}

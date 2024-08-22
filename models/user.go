package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1/social-network-auth/errors"
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
) *errors.ApiError {
	err := db.Where(&User{Username: username}).Or(&User{Email: email}).First(&User{}).Error
	if err == nil {
		return errors.BadRequestError("username or email already exists")
	}

	if !utils.IsRecordNotFound(err) {
		return errors.InternalServerError(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	user.Password = string(hashedPassword)

	user.ID = uuid.New()
	user.Level = BRONZE
	user.Username = username
	user.Email = email

	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	if err := db.Create(&user).Error; err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (user *User) Authenticate(db *gorm.DB, username string, password string) (*User, *errors.ApiError) {
	var dbUser User
	if err := db.Where(&User{Username: username}).First(&dbUser).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return nil, errors.BadRequestError("user %s not found", username)
		}

		return nil, errors.InternalServerError(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); err != nil {
		return nil, errors.BadRequestError("invalid password")
	}

	return &dbUser, nil
}

func (user *User) UpdateLevel(db *gorm.DB, userID uuid.UUID, levelName string) *errors.ApiError {
	err := db.Where(&User{ID: userID}).First(&user).Error
	if err != nil {
		if utils.IsRecordNotFound(err) {
			return errors.BadRequestError("user %v not found", userID)
		}

		return errors.InternalServerError(err.Error())
	}

	level := GetLevelFromName(levelName)
	if user.Level == level {
		return nil
	}
	user.Level = level

	user.UpdatedAt = time.Now().UTC()

	if err := db.Save(user).Error; err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (user *User) UpdatePassword(db *gorm.DB, currentPassword string, newPassword string) *errors.ApiError {
	if err := db.Where(&User{ID: user.ID}).First(&user).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return errors.BadRequestError("user not found")
		}

		return errors.InternalServerError(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.BadRequestError("invalid password")
	}

	if currentPassword == newPassword {
		return errors.BadRequestError("new password can not be the same as current one")
	}

	if err := user.updatePassword(db, newPassword); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (user *User) ResetPassword(db *gorm.DB, email string, resetToken string, newPassword string, confirmPassword string) *errors.ApiError {
	if newPassword != confirmPassword {
		return errors.BadRequestError("comfirm password does not match with new password")
	}

	var u User
	if err := db.Where(&User{Email: email}).First(&u).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return errors.BadRequestError("can not find user with email %v", email)
		}

		return errors.InternalServerError(err.Error())
	}

	if err := db.Where("token = ? AND token_expiry > ? AND user_id = ?", resetToken, time.Now().UTC(), u.ID).First(&PasswordResetToken{}).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return errors.BadRequestError("invalid or expired token")
		}

		return errors.InternalServerError(err.Error())
	}

	if err := u.updatePassword(db, newPassword); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (user *User) GetUserByUsernameOrEmail(db *gorm.DB, userIdentity string) (*User, *errors.ApiError) {
	var dbUser User
	if err := db.Where(&User{Username: userIdentity}).Or(&User{Email: userIdentity}).First(&dbUser).Error; err != nil {
		if utils.IsRecordNotFound(err) {
			return nil, errors.BadRequestError("user not found")
		}

		return nil, errors.InternalServerError(err.Error())
	}

	return &dbUser, nil
}

func (user *User) updatePassword(db *gorm.DB, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

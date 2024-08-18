package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	UserID      uint      `json:"user_id" gorm:"unique,not null"`
	FirstName   string    `json:"first_name" gorm:"not null"`
	LastName    string    `json:"last_name" gorm:"not null"`
	DateOfBirth string    `json:"date_of_birth" gorm:"not null"`
	Address     string    `json:"address" gorm:"not null"`
	Phone       string    `json:"phone" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;autoCreateTime:false"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;autoUpdateTime:false"`
	DeletedAt   gorm.DeletedAt
}

func (p *Profile) CreateProfile(
	db *gorm.DB,
	userID uint,
	firstName string,
	lastName string,
	dateOfBirth string,
	address string,
	phone string,
) (*Profile, error) {
	err := db.Where(&Profile{UserID: userID}).First(&Profile{}).Error
	if err == nil {
		return nil, errors.New("profile for current user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	err = db.Where(&User{ID: userID}).First(&User{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user %v not found", userID)
		}
		return nil, err
	}

	profile := Profile{
		UserID:      userID,
		FirstName:   firstName,
		LastName:    lastName,
		DateOfBirth: dateOfBirth,
		Address:     address,
		Phone:       phone,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	err = db.Create(&profile).Error
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (p *Profile) GetProfileByUser(db *gorm.DB) (*Profile, error) {
	var dbProfile Profile
	err := db.Where(&Profile{UserID: p.UserID}).First(&dbProfile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("profile for user %v not found", p.UserID)
		}
		return nil, err
	}

	return &dbProfile, nil
}

func (p *Profile) UpdateProfileByUser(db *gorm.DB) (*Profile, error) {
	var profile Profile
	err := db.Where(&Profile{UserID: p.UserID}).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("profile for user %v not found", p.UserID)
		}
		return nil, err
	}

	p.ID = profile.ID
	p.UpdatedAt = time.Now().UTC()

	err = db.Updates(&p).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}

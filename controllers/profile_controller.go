package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/models"
	"github.com/vantutran2k1/social-network-auth/utils"
	"net/http"
)

type CreateProfileRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	FirstName   string `json:"first_name" binding:"required,max=32"`
	LastName    string `json:"last_name" binding:"required,max=32"`
	DateOfBirth string `json:"date_of_birth" binding:"required,date,beforeToday"`
	Address     string `json:"address" binding:"required,max=255"`
	Phone       string `json:"phone" binding:"required,phoneNumber"`
}

func CreateProfile(c *gin.Context) {
	var request CreateProfileRequest
	errs := utils.BindAndValidate(c, &request)
	if errs != nil && len(errs) > 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	profile := models.Profile{
		UserID:      request.UserID,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		DateOfBirth: request.DateOfBirth,
		Address:     request.Address,
		Phone:       request.Phone,
	}
	p, err := profile.CreateProfile(config.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]any{
		"user_id":       p.UserID,
		"first_name":    p.FirstName,
		"last_name":     p.LastName,
		"date_of_birth": p.DateOfBirth,
		"address":       p.Address,
		"phone":         p.Phone,
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

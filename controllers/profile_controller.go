package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/models"
	"github.com/vantutran2k1/social-network-auth/utils"
	"net/http"
	"strconv"
)

type CreateProfileRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	FirstName   string `json:"first_name" binding:"required,max=32"`
	LastName    string `json:"last_name" binding:"required,max=32"`
	DateOfBirth string `json:"date_of_birth" binding:"required,date,beforeToday"`
	Address     string `json:"address" binding:"required,max=255"`
	Phone       string `json:"phone" binding:"required,phoneNumber"`
}

type ProfileResponse struct {
	ID          uint   `json:"id,omitempty"`
	UserID      uint   `json:"user_id,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Address     string `json:"address,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

func CreateProfile(c *gin.Context) {
	var request CreateProfileRequest
	errs := utils.BindAndValidate(c, &request)
	if errs != nil && len(errs) > 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	p := models.Profile{
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		DateOfBirth: request.DateOfBirth,
		Address:     request.Address,
		Phone:       request.Phone,
	}
	profile, err := p.CreateProfile(config.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := ProfileResponse{
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		DateOfBirth: profile.DateOfBirth,
		Address:     profile.Address,
		Phone:       profile.Phone,
	}
	c.JSON(http.StatusCreated, gin.H{"data": getProfileResponseData(r)})
}

func GetProfile(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusOK, gin.H{"data": make(map[string]any)})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be integer"})
		return
	}

	p := models.Profile{UserID: uint(userID)}
	profile, err := p.GetProfileByUser(config.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := ProfileResponse{
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		DateOfBirth: profile.DateOfBirth,
		Address:     profile.Address,
		Phone:       profile.Phone,
	}
	c.JSON(http.StatusOK, gin.H{"data": getProfileResponseData(r)})
}

func GetCurrentProfile(c *gin.Context) {
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"data": make(map[string]any)})
	}

	p := models.Profile{UserID: userID.(uint)}
	profile, err := p.GetProfileByUser(config.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := ProfileResponse{
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		DateOfBirth: profile.DateOfBirth,
		Address:     profile.Address,
		Phone:       profile.Phone,
	}
	c.JSON(http.StatusOK, gin.H{"data": getProfileResponseData(r)})
}

func getProfileResponseData(r ProfileResponse) map[string]any {
	data, err := json.Marshal(r)
	if err != nil {
		return nil
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}

	return result
}

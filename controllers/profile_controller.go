package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/middlewares"
	"github.com/vantutran2k1/social-network-auth/models"
	"github.com/vantutran2k1/social-network-auth/utils"
)

type CreateProfileRequest struct {
	FirstName   string `json:"first_name" binding:"required,max=32"`
	LastName    string `json:"last_name" binding:"required,max=32"`
	DateOfBirth string `json:"date_of_birth" binding:"required,date,beforeToday"`
	Address     string `json:"address" binding:"required,max=255"`
	Phone       string `json:"phone" binding:"required,phoneNumber"`
}

type UpdateProfileRequest struct {
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
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	userID, err := middlewares.GetUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p := models.Profile{}
	profile, err := p.CreateProfile(
		config.DB,
		userID,
		request.FirstName,
		request.LastName,
		request.DateOfBirth,
		request.Address,
		request.Phone,
	)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid syntax for user_id"})
		return
	}

	p := models.Profile{}
	profile, err := p.GetProfileByUser(config.DB, userID)
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
	userID, err := middlewares.GetUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p := models.Profile{}
	profile, err := p.GetProfileByUser(config.DB, userID)
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

func UpdateCurrentProfile(c *gin.Context) {
	var request UpdateProfileRequest
	errs := utils.BindAndValidate(c, &request)
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	userID, err := middlewares.GetUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p := models.Profile{}
	if err := p.UpdateProfileByUser(config.DB, userID, request.FirstName, request.LastName, request.DateOfBirth, request.Address, request.Phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := ProfileResponse{
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		DateOfBirth: p.DateOfBirth,
		Address:     p.Address,
		Phone:       p.Phone,
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

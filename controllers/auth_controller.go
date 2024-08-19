package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/middlewares"
	"github.com/vantutran2k1/social-network-auth/models"
	"github.com/vantutran2k1/social-network-auth/utils"
)

type UserRegistrationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

type UserAuthenticationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8,max=32"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=32"`
}

type UpdateLevelRequest struct {
	UserId    uint   `json:"user_id" binding:"required"`
	LevelName string `json:"level_name" binding:"required,oneof=BRONZE SILVER GOLD"`
}

func Register(c *gin.Context) {
	var creds UserRegistrationRequest
	errs := utils.BindAndValidate(c, &creds)
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	user := models.User{}
	err := user.Register(config.DB, creds.Username, creds.Password, creds.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]any{
		"username": user.Username,
		"email":    user.Email,
	}
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

func Login(c *gin.Context) {
	var auth UserAuthenticationRequest
	errs := utils.BindAndValidate(c, &auth)
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	user := models.User{}
	loginUser, err := user.Authenticate(config.DB, auth.Username, auth.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	t := models.Token{}
	token, err := t.CreateLoginToken(config.DB, loginUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": token.Token})
}

func Logout(c *gin.Context) {
	var token models.Token
	err := token.Revoke(config.DB, middlewares.GetAuthTokenFromRequest(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token revoked successfully"})
}

func UpdatePassword(c *gin.Context) {
	userID, err := middlewares.GetUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var request UpdatePasswordRequest
	errs := utils.BindAndValidate(c, &request)
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	u := models.User{ID: userID}
	err = u.UpdatePassword(config.DB, request.CurrentPassword, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

func UpdateUserLevel(c *gin.Context) {
	var request UpdateLevelRequest
	errs := utils.BindAndValidate(c, &request)
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	user := models.User{}
	err := user.UpdateLevel(config.DB, request.UserId, request.LevelName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]any{
		"user_id": user.ID,
		"level":   request.LevelName,
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

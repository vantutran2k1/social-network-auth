package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/config"
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

func Register(c *gin.Context) {
	var creds UserRegistrationRequest
	errs := utils.BindAndValidate(c, &creds)
	if errs != nil && len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	user := models.User{
		Username: creds.Username,
		Password: creds.Password,
		Email:    creds.Email,
	}
	err := user.Register(config.DB)
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
	if errs != nil && len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	user := models.User{Username: auth.Username}
	if !user.Authenticate(config.DB, auth.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	tokenString, err := user.GenerateToken(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tokenString})
}

func Logout(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	tokenString = utils.GetTokenFromString(tokenString)

	var token models.Token
	err := token.Revoke(config.DB, tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token revoked successfully"})
}

func UpdatePassword(c *gin.Context) {
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"error": make(map[string]any)})
	}

	var request UpdatePasswordRequest
	errs := utils.BindAndValidate(c, &request)
	if errs != nil && len(errs) > 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	u := models.User{ID: userID.(uint)}
	err := u.UpdatePassword(config.DB, request.CurrentPassword, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/models"
	"net/http"
	"time"
)

type UserRegistrationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserAuthenticationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var creds UserRegistrationRequest
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: creds.Username,
		Password: creds.Password,
		Email:    creds.Email,
	}
	err = user.Register(config.DB)
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
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Username: auth.Username}
	if !user.Authenticate(config.DB, auth.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	tokenString, expirationTime, err := user.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token := models.Token{
		UserId:    user.ID,
		Token:     tokenString,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: expirationTime,
	}
	err = token.Save(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tokenString})
}

func Validate(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token := models.Token{}
	if !token.Validate(config.DB, tokenString) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token is valid"})
}

func Logout(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token := models.Token{}
	err := token.Revoke(config.DB, tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token revoked successfully"})
}

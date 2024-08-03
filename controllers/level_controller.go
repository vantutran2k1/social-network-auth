package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/models"
	"github.com/vantutran2k1/social-network-auth/utils"
	"net/http"
)

type AssignLevelRequest struct {
	UserId    uint   `json:"user_id" binding:"required"`
	LevelName string `json:"level_name" binding:"required,oneof=BRONZE SILVER GOLD"`
}

func AssignLevelToUser(c *gin.Context) {
	var request AssignLevelRequest
	errs := utils.BindAndValidate(c, &request)
	if errs != nil && len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	user := models.User{ID: request.UserId}
	err := user.AssignLevel(config.DB, request.LevelName)
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

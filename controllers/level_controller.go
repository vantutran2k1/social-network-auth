package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/models"
	"net/http"
)

func GetLevels(c *gin.Context) {
	level := models.Level{}
	levels, err := level.GetLevels(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, map[string]any{"data": levels})
}

package middlewares

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1/social-network-auth/config"
	"github.com/vantutran2k1/social-network-auth/models"
	"github.com/vantutran2k1/social-network-auth/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		tokenString = utils.GetTokenFromString(tokenString)

		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		var dbToken models.Token
		if !dbToken.Validate(config.DB, tokenString) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token not found or expired"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func GetUserIDFromRequest(c *gin.Context) (uint, error) {
	userID, exist := c.Get("user_id")
	if !exist {
		return 0, errors.New("can not get user id from request")
	}

	return userID.(uint), nil
}

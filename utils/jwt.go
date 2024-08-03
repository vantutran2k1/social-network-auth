package utils

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GetTokenFromString(tokenString string) string {
	token := tokenString
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) == 2 {
		token = tokenParts[1]
	}

	return token
}

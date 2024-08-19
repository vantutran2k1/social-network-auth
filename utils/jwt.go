package utils

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

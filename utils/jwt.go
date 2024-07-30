package utils

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

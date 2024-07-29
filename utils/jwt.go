package utils

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

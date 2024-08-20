package utils

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

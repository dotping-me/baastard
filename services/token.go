package services

import (
	"os"
	"time"

	"github.com/dotping-me/baastard/schema"
	"github.com/golang-jwt/jwt/v5"
)

// JWT Implementation for user request authentication.

type Claims struct {
	ID         string `json:"id"`
	Email      string `json:"email"` // TODO: Maybe store the full user
	TokenState uint   // TODO: Make this toggable, so as to reduce compute times on requests

	jwt.RegisteredClaims
}

// Generate a token string for a given user
func GenerateToken(user *schema.User) (string, error) {
	claims := Claims{
		ID:    user.ID,
		Email: user.Email,

		// Default is 0. After logout, state is incremented and thus
		// User SessionState != JWT SessionState
		TokenState: user.TokenState,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)), // TODO: Make this configurable
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Signs token using secret
	// TODO: Make secret also configurable

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("SECRET"))

	return token.SignedString(secret)
}

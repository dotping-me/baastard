package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dotping-me/baastard/schema"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: Move to services/token.go
type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"` // TODO: Maybe store the full user

	jwt.RegisteredClaims
}

// TODO: Move to services/token.go
// Generate a token string for a given auth user
func GenerateToken(user *schema.Auth) (string, error) {
	claims := Claims{
		ID:    user.ID,
		Email: user.Email,

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

// Authentication middleware that checks if user is authenticated through its JWT
// token.
func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Getting Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing 'Authorization' header!",
			})

			return
		}

		// Getting the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing 'Bearer <Token>'!",
			})

			return
		}

		// Validating the JWT token
		claims := &Claims{}
		secret := []byte(os.Getenv("SECRET"))
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token!",
			})

			return
		}

		// All good!
		// JWT is valid, proceed to next function

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

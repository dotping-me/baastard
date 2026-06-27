package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dotping-me/baastard/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Authentication middleware that checks if user is authenticated through its JWT
// token.
func AuthRequired(auth *services.AuthService) gin.HandlerFunc {
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
		claims := &services.Claims{}
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

		// Validating the session state: Has the user logged out but JWT has not
		// expired yet?

		// NOTE: Extra DB query for every protected route

		user, err := auth.FindByID(claims.ID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "No user for this access token!",
			})

			return
		}

		if user.TokenState != claims.TokenState {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Session revoked!",
			})

			return
		}

		// All good!
		// JWT is valid, proceed to next function

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

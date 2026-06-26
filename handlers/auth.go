package handlers

import (
	"net/http"

	"github.com/dotping-me/baastard/middleware"
	"github.com/dotping-me/baastard/payload"
	"github.com/dotping-me/baastard/services"
	"github.com/gin-gonic/gin"
)

// Here routes and handler functions are defined and exposed for given
// service, in this case: Authentication Service.

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Handler for registering a user @ POST /api/auth/register
func (h *AuthHandler) Register(ctx *gin.Context) {

	var req payload.RegisterAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Registers the user using the Auth Service
	user, err := h.auth.Register(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"verified":  user.Verified,
		"createdAt": user.CreatedAt,
	})
}

// Handler for logging in a user @ POST /api/auth/login
func (h *AuthHandler) Login(ctx *gin.Context) {

	var req payload.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Logs in the user using the Auth Service
	user, err := h.auth.Login(req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Successful login, now generates a JWT token for that user
	token, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token, // TODO: Maybe use HTTP Only Cookie + CORS
	})
}

// Helper handler to return information about self + self-auth
func (h *AuthHandler) Me(ctx *gin.Context) {
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing claims!",
		})

		return
	}

	// Casting for type safety (and also so that I get autocompletion)
	jwtClaims, ok := claims.(*middleware.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing claims!",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":    jwtClaims.ID,
		"email": jwtClaims.Email,
	})
}

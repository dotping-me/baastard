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
	auth     *services.AuthService
	sessions *services.SessionService
}

func NewAuthHandler(auth *services.AuthService, sessions *services.SessionService) *AuthHandler {
	return &AuthHandler{
		auth:     auth,
		sessions: sessions,
	}
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

	// Successful login, now generates an access token for that user
	accessToken, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Creates a new session for that user
	session, err := h.sessions.Create(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": session.RefreshToken,

		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// Handler that refreshes a session for a logged in user
func (h *AuthHandler) Refresh(ctx *gin.Context) {

	var req payload.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Finds the valid session for that user
	session, err := h.sessions.Get(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	user, err := h.auth.FindByID(session.AuthID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Session!",
		})

		return
	}

	// User has been validated, regenerate access token
	accessToken, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Updates session accordingly and revoke old refresh token
	session, err = h.sessions.Refresh(session)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": session.RefreshToken,
	})
}

// Handler to log out a user
func (h *AuthHandler) Logout(ctx *gin.Context) {

	var req payload.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Deletes the session
	err := h.sessions.Delete(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out",
	})
}

// Handler to logout the user from all sessions
func (h *AuthHandler) LogoutAll(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middleware.Claims)
	err := h.sessions.DeleteAll(claims.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out all devices",
	})
}

// Helper handler to return information about self + self-auth
func (h *AuthHandler) Me(ctx *gin.Context) {

	// TODO: Turn this into a modular helper function or set claims explicitly
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

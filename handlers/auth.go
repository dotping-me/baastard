package handlers

import (
	"net/http"

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
	accessToken, err := services.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Creates a new session for that user

	deviceID := services.GetDeviceIDOrCreate(ctx) // Sets the Cookie
	session, err := h.sessions.Create(user.ID, deviceID,
		ctx.Request.UserAgent(), ctx.ClientIP())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	services.SetRefreshCookie(ctx, session.RefreshToken) // Sets the HTTP Only Cookie

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,

		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// Handler that refreshes a session for a logged in user
func (h *AuthHandler) Refresh(ctx *gin.Context) {

	refreshToken, err := services.GetRefreshTokenFromCookie(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing refresh cookie!",
		})

		return
	}

	// Finds the valid session for that user
	session, err := h.sessions.Get(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid session!",
		})

		return
	}

	// Device Binding Validation
	deviceID, err := services.GetExistingDeviceID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = h.sessions.ValidateDevice(session, deviceID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Device does not match!",
		})

		return
	}

	user, err := h.auth.FindByID(session.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid session!",
		})

		return
	}

	// User has been validated, regenerate access token
	accessToken, err := services.GenerateToken(user)
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

	services.SetRefreshCookie(ctx, session.RefreshToken) // Updates Cookie

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	})
}

// Handler to log out a user
func (h *AuthHandler) Logout(ctx *gin.Context) {

	refreshToken, err := services.GetRefreshTokenFromCookie(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing refresh cookie!",
		})

		return
	}

	// Deletes the session
	err = h.sessions.Delete(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	services.ClearRefreshCookie(ctx)

	// NOTE: Access tokens are not revoked here because it may happen that
	//		 the user didn't mean to log out (i.e. token expired and user was
	//		 forced out).

	// 		 However, a user can intentionally log out and still have a valid JWT
	// 		 token.

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out",
	})
}

// Handler to logout the user from all sessions
func (h *AuthHandler) LogoutAll(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*services.Claims)

	err := h.auth.IncrementTokenState(claims.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = h.sessions.DeleteAll(claims.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	services.ClearRefreshCookie(ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged out all devices",
	})
}

// Helper handler to return information about self + self-auth
func (h *AuthHandler) Me(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*services.Claims)
	ctx.JSON(http.StatusOK, gin.H{
		"id":    claims.ID,
		"email": claims.Email,
	})
}

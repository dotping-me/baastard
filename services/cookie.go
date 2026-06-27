package services

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Using HTTP Only Cookies to set refresh tokens, preventing browser JS from
// manipulating them.

const (
	RefreshCookieName   = "refresh_token"
	RefreshCookieMaxAge = 30 * 24 * time.Hour // TODO: Make this configurable
)

func SetRefreshCookie(ctx *gin.Context, token string) {
	ctx.SetCookie(RefreshCookieName, token, int(RefreshCookieMaxAge.Seconds()),
		"/", "", false, true) // TODO: Make secure tag be configured from an ENV key
}

func ClearRefreshCookie(ctx *gin.Context) {
	ctx.SetCookie(RefreshCookieName, "", -1, "/", "", false, true)
}

func GetRefreshTokenFromCookie(ctx *gin.Context) (string, error) {
	return ctx.Cookie(RefreshCookieName)
}

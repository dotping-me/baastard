package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// This HTTP server exposes routes only to the host machine. Having a seperate
// HTTP server acts as a security boundary.

func NewPrivateRouter() *gin.Engine {
	r := gin.Default()

	internal := r.Group("/internal")
	{
		internal.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
	}

	return r
}

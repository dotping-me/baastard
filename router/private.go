package router

import (
	"io/fs"
	"net/http"

	"github.com/dotping-me/baastard/public"
	"github.com/dotping-me/baastard/services"
	"github.com/dotping-me/baastard/views/pages"
	"github.com/gin-gonic/gin"
)

// This HTTP server exposes routes only to the host machine. Having a seperate
// HTTP server acts as a security boundary.

func NewPrivateRouter(cfg *services.TableConfigService) *gin.Engine {
	r := gin.Default()

	// Embedding static files
	cssFS, _ := fs.Sub(public.Files, "css")
	jsFS, _ := fs.Sub(public.Files, "js")
	assetsFS, _ := fs.Sub(public.Files, "assets")

	r.StaticFS("/css", http.FS(cssFS))
	r.StaticFS("/js", http.FS(jsFS))
	r.StaticFS("/assets", http.FS(assetsFS))

	// API routes
	internal := r.Group("/internal")
	{
		internal.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})

		internal.GET("/tables", func(ctx *gin.Context) {
			tables, err := cfg.GetAll()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})

				return
			}

			ctx.JSON(http.StatusOK, tables)
		})
	}

	// UI screens
	r.GET("/", func(c *gin.Context) {
		pages.Dashboard().Render(c.Request.Context(), c.Writer)
	})

	return r
}

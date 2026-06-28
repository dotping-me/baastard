package router

import (
	"github.com/dotping-me/baastard/handlers"
	"github.com/dotping-me/baastard/middleware"
	"github.com/dotping-me/baastard/services"
	"github.com/gin-gonic/gin"
)

// This starts a public HTTP server, accessible to users. It provides seperation
// from another HTTP server meant for localhost only.

func NewPublicRouter(auth *handlers.AuthHandler, authService *services.AuthService,
	crud *handlers.CrudHandler, cfg *services.TableConfigService) *gin.Engine {

	r := gin.Default()

	// API Routes
	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/refresh", auth.Refresh)
			authGroup.POST("/logout", auth.Logout)
			authGroup.POST("/logout-all",
				auth.AuthRequired(), auth.LogoutAll)

			authGroup.GET("/me", auth.AuthRequired(), auth.Me)
		}

		crudGroup := api.Group("/crud")
		crudGroup.Use(
			middleware.VerifyExposable(cfg),
			middleware.VerifyAccess(authService),
		)
		{
			crudGroup.GET("/:table", crud.GetRows)
			crudGroup.GET("/:table/:id", crud.GetRow)
			crudGroup.POST("/:table", crud.Insert)
			crudGroup.PATCH("/:table/:id", crud.Update)
			crudGroup.DELETE("/:table/:id", crud.Delete)
		}
	}

	return r
}

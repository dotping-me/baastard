package middleware

import (
	"net/http"

	"github.com/dotping-me/baastard/schema"
	"github.com/dotping-me/baastard/services"
	"github.com/gin-gonic/gin"
)

// Validates if table is exposable and includes its configuration in Gin context
func VerifyExposable(s *services.TableConfigService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		table := ctx.Param("table")
		if table == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing table name!",
			})

			return
		}

		config, err := s.Get(table)
		if err != nil {
			// Deliberately return 404 instead of exposing existence
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Table not found!",
			})

			return
		}

		// Hidden tables are treated as if they do not exist
		if config.Hidden {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Table not found!",
			})

			return
		}

		// Make available to downstream middleware and handlers.
		ctx.Set("table", table)
		ctx.Set("tableConfig", config)

		ctx.Next()
	}
}

func VerifyAccess(s *services.AuthService) gin.HandlerFunc {
	authMiddleware := AuthRequired(s)

	return func(ctx *gin.Context) {
		cfg := ctx.MustGet("tableConfig").(*schema.TableConfig)

		if cfg.AuthRequired {
			authMiddleware(ctx)
			if ctx.IsAborted() {
				return
			}
		}

		if cfg.ReadOnly {
			switch ctx.Request.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Table is read-only!",
				})

				return
			}
		}

		ctx.Next()
	}
}

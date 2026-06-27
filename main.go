package main

import (
	"database/sql"
	"os"

	"github.com/dotping-me/baastard/handlers"
	"github.com/dotping-me/baastard/middleware"
	"github.com/dotping-me/baastard/schema"
	"github.com/dotping-me/baastard/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "turso.tech/database/tursogo"
)

func main() {
	godotenv.Load()

	// Establishing Turso connection
	tursoConn, err := sql.Open("turso", "__db.db")
	if err != nil {
		panic(err)
	}

	defer tursoConn.Close()

	// Configuring GORM for Turso connection
	db, err := gorm.Open(sqlite.Dialector{Conn: tursoConn}, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Running database migrations
	db.Exec("PRAGMA foreign_keys = ON")
	db.AutoMigrate(
		&schema.User{},
		&schema.Session{},
		&schema.EntityMetadata{},
		&schema.FieldMetadata{},
	)

	// Registering services and handlers
	authService := services.NewAuthService(db)
	sessionService := services.NewSessionService(db)
	metadataService := services.NewMetadataService(db)

	authHandler := handlers.NewAuthHandler(authService, sessionService)
	metadataHandler := handlers.NewMetadataHandler(metadataService)

	// Setting up the HTTP server
	r := gin.Default()

	// API Routes
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/logout-all",
				middleware.AuthRequired(authService), authHandler.LogoutAll)

			auth.GET("/me", middleware.AuthRequired(authService), authHandler.Me)
		}

		// TODO: Add Auth middleware
		metadata := api.Group("/metadata")
		{
			metadata.POST("", metadataHandler.Create)

			metadata.GET("", metadataHandler.GetAll)
			metadata.GET("/:id", metadataHandler.GetById)

			metadata.DELETE("/:id", metadataHandler.Delete)
		}
	}

	// Configuring server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Starting the server
	err = r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

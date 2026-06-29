package main

import (
	"database/sql"
	"os"

	"github.com/dotping-me/baastard/handlers"
	"github.com/dotping-me/baastard/router"
	"github.com/dotping-me/baastard/schema"
	"github.com/dotping-me/baastard/services"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "turso.tech/database/tursogo"
)

func main() {
	godotenv.Load()

	// Establishing Turso connection
	// TODO: Make it configurable so that the database is encrypted using Turso's
	// 		 native encryption.

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
		&schema.TableConfig{},
	)

	// Syncs up table_configs table
	tableConfigService := services.NewTableConfigService(db)
	err = tableConfigService.Sync()
	if err != nil {
		panic(err)
	}

	// Registering services and handlers
	authService := services.NewAuthService(db)
	sessionService := services.NewSessionService(db)
	authHandler := handlers.NewAuthHandler(authService, sessionService)

	cfgService := services.NewTableConfigService(db)
	crudService := services.NewCrudService(tursoConn, cfgService)
	crudHandler := handlers.NewCrudHandler(crudService)

	// Starting the the public HTTP server
	publicPort := os.Getenv("PUBLICPORT")
	if publicPort == "" {
		publicPort = "8080"
	}

	public := router.NewPublicRouter(authHandler, authService, crudHandler, cfgService)
	go func() {
		err := public.Run(":" + publicPort)
		if err != nil {
			panic(err)
		}
	}()

	// Starting the private HTTP server, meant only for localhost
	privatePort := os.Getenv("PRIVATEPORT")
	if privatePort == "" {
		privatePort = "9090"
	}

	// TODO: Guard behind Unix Socket Authentication + Admin Token or Auth Login
	private := router.NewPrivateRouter(cfgService)
	err = private.Run("127.0.0.1:" + privatePort)
	if err != nil {
		panic(err)
	}
}

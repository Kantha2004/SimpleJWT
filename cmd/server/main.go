package main

import (
	"fmt"
	"log"

	"github.com/Kantha2004/SimpleJWT/internal/api"
	"github.com/Kantha2004/SimpleJWT/internal/api/handlers"
	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/config"
	"github.com/Kantha2004/SimpleJWT/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	config := config.LoadConfig()

	database, err := db.NewDatabase(config.DBPath)
	if err != nil {
		log.Fatal("Error while connectiing to DB", err.Error())
	}
	defer database.Close()

	// Initialize repositories
	// userRepo := db.NewUserRepository(database)

	// Set Gin to release mode in production
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	jwtService := auth.NewJWTService(config.JWTSecret)

	deps := api.NewDependencies(jwtService)
	handlerDeps := handlers.NewDependencies(database)

	router := gin.Default()

	api.SetupGinRoutes(router, deps, handlerDeps)

	// Get port from config
	port := ":" + config.Port
	fmt.Printf("SimpleJWT server starting on port %s\n", port)

	// Start server using Gin's Run method
	log.Fatal(router.Run(port))
}

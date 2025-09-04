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

	// Updated Swagger imports
	swaggerFiles "github.com/swaggo/files" // This replaces swaggerFiles
	ginSwagger "github.com/swaggo/gin-swagger"

	// Import the generated docs (after running swag init)
	_ "github.com/Kantha2004/SimpleJWT/docs"
)

// @title SimpleJWT API
// @version 1.0
// @description A simple JWT authentication service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	_ = godotenv.Load()

	loadConfig := config.LoadConfig()

	database, err := db.NewDatabase(loadConfig.DBPath)
	if err != nil {
		log.Fatal("Error while connecting to DB", err.Error())
	}
	defer func(database *db.Database) {
		err := database.Close()
		if err != nil {
			log.Fatal("Failed to close DB connection")
		}
	}(database)

	// Set Gin to release mode in production
	if loadConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	jwtService := auth.NewJWTService(loadConfig.JWTSecret)

	// Add nil check
	if jwtService == nil {
		log.Fatal("Failed to create JWT service")
	}

	deps := api.NewDependencies(jwtService)
	handlerDeps := handlers.NewDependencies(database, jwtService)

	router := gin.Default()

	api.SetupGinRoutes(router, deps, handlerDeps)

	// Add Swagger route - notice the updated import usage
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from loadConfig
	port := ":" + loadConfig.Port
	fmt.Printf("SimpleJWT server starting on port %s\n", port)
	fmt.Printf("Swagger documentation available at: http://localhost%s/swagger/index.html\n", port)

	// Start server using Gin's Run method
	log.Fatal(router.Run(port))
}

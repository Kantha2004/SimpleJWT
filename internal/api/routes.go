package api

import (
	"fmt"
	"net/http"

	"github.com/Kantha2004/SimpleJWT/internal/api/handlers"
	"github.com/Kantha2004/SimpleJWT/internal/api/services"
	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/repositories"
	"github.com/gin-gonic/gin"
)

// testHandler godoc
// @Summary      Test endpoint
// @Description  Returns a success response if JWT is valid
// @Tags         Protected
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Test successful"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /protected/test [get]
// @Security     BearerAuth
func testHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		apiresponse.SendInternalError(c, "Unable to get userId")
		return
	}
	fmt.Println("User ID:", userID)
	apiresponse.SendSuccess(c, http.StatusOK, gin.H{"user_id": userID}, "Test successful")
}

func SetupGinRoutes(router *gin.Engine, deps *Dependencies) {
	// Add global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize services
	userRepo := repositories.NewUserRepository(deps.DB)
	clientRepo := repositories.NewClientRepository(deps.DB)
	// clientUserRepo := repositories.NewClientUserRepository(deps.DB, "")

	// Initialize services with dependencies
	userService := services.NewUserService(userRepo, deps.JWTService)
	clientService := services.NewClientService(clientRepo, deps.DB)
	// clientUserService := services.NewClientUserService(clientRepo, clientUserRepo, deps.JWTService)
	authService := services.NewAuthService(userService)

	// Initialize handlers with services
	userHandler := handlers.NewUserHandler(userService, authService)
	clientHandler := handlers.NewClientHandler(clientService, authService)
	// clientUserHandler := handlers.NewClientUserHandler(clientUserService, authService, deps.JWTService)

	v1 := router.Group("api/v1")
	{
		// GET Methods
		v1.GET("/ping", PingHandler)

		// POST Methods
		v1.POST("/createUser", userHandler.CreateUser)
		v1.POST("/login", userHandler.Login)
	}

	// Client routes (public)
	// client := router.Group("api/v1/client")
	{
		// client.POST("/userlogin", clientUserHandler.ClientUserLogin)
	}

	// Protected routes
	protected := router.Group("api/v1/protected")
	protected.Use(JWTMiddleware(deps.JWTService))
	{
		// GET Methods
		protected.GET("/test", testHandler)
		protected.GET("/getAllClients", clientHandler.GetAllClients)

		// POST Methods
		protected.POST("/createClient", clientHandler.CreateClient)
		// protected.POST("/createClientUser", clientUserHandler.CreateClientUser)
	}
}

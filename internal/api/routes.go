package api

import (
	"fmt"
	"net/http"

	"github.com/Kantha2004/SimpleJWT/internal/api/handlers"
	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/gin-gonic/gin"
)

// testHandler godoc
// @Summary      Test endpoint
// @Description  Returns a success response if JWT is valid
// @Tags         Protected
// @Accept       json
// @Produce      json
// @Success      201  {object}  map[string]interface{}  "User created successfully"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /protected/test [get]
// @Security     BearerAuth
func testHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		apiresponse.SendInternalError(c, "Unable to get userId")
	}
	fmt.Println(userID)
	apiresponse.SendSuccess(c, http.StatusCreated, struct{}{}, "User created successfully")
}

func SetupGinRoutes(router *gin.Engine, deps *Dependencies, handlerDeps *handlers.Dependencies) {
	// Add global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("api/v1")
	{
		// GET Methods
		v1.GET("/ping", PingHandler)

		// POST Methods
		v1.POST("/createUser", handlerDeps.CreateUser)
		v1.POST("/login", handlerDeps.Login)
	}

	protected := router.Group("api/v1/protected")
	protected.Use(JWTMiddleware(deps.JWTService))
	{
		// GET Methods
		protected.GET("/test", testHandler)
		protected.GET("/getAllClients", handlerDeps.GetAllClients)

		// POST Methods
		protected.POST("/createClient", handlerDeps.CreateClient)
		protected.POST("/createClientUser", handlerDeps.CreateClientUser)
	}
}

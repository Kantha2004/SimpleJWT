package api

import (
	"github.com/Kantha2004/SimpleJWT/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupGinRoutes(router *gin.Engine, deps *Dependencies, handlerDeps *handlers.Dependencies) {
	// Add global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("api/v1")
	{
		v1.GET("/ping", PingHandler)
		v1.POST("/createUser", handlerDeps.CreateUser)
		v1.POST("/login", handlerDeps.Login)
	}
}

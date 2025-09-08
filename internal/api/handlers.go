package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler godoc
// @Summary Health check endpoint
// @Description Ping endpoint to check if the service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Service status response"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"service": "SimpleJWT",
		"message": "Application is running...",
	})
}

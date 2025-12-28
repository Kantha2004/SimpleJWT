package utils

import (
	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/gin-gonic/gin"
)

// WithBody is a higher-order function that wraps a Gin handler.
// It binds the request body to a struct of type T, validates it,
// and then calls the provided handler function with the context and the validated struct.
func WithBody[T any](handler func(*gin.Context, T)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindJSON(&req); err != nil {
			apiresponse.SendValidationError(c, err)
			return
		}
		handler(c, req)
	}
}

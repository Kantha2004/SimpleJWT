package utils

import (
	"errors"

	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/gin-gonic/gin"
)

func VerifyRequestModel(c *gin.Context, req any) bool {
	if err := c.ShouldBindBodyWithJSON(req); err != nil {
		apiresponse.SendValidationError(c, err)
		return false
	}
	return true
}

// GetUserIDFromContext extracts and validates user ID from context
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("unable to find userID from the request")
	}

	userID, ok := userIDVal.(int)
	if !ok {
		return 0, errors.New("invalid userID type")
	}

	return uint(userID), nil
}

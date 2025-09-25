package handlers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/repositories"
	"github.com/Kantha2004/SimpleJWT/internal/utils"
	"github.com/gin-gonic/gin"
)

const (
	USER_NOT_FOUND = "user not found"
)

// ValidateUser fetches and validates user existence
// Returns the user if found, otherwise returns an error
func (d *Dependencies) ValidateUser(userID uint) (*models.AdminUser, error) {
	userRepo := repositories.NewUserRepository(d.DB)
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	fmt.Println("user", user)

	if user == nil {
		return nil, errors.New(USER_NOT_FOUND)
	}

	return user, nil
}

// ValidateUserFromContext extracts user ID from context, validates the user,
// and handles HTTP error responses automatically
func (d *Dependencies) ValidateUserFromContext(c *gin.Context) (*models.AdminUser, bool) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		log.Printf("Failed to extract user ID from context: %v", err)
		apiresponse.SendValidationError(c, err)
		return nil, false
	}

	user, err := d.ValidateUser(userID)
	if err != nil {
		d.handleUserValidationError(c, err)
		return nil, false
	}

	return user, true
}

// handleUserValidationError centralizes error handling for user validation failures
func (d *Dependencies) handleUserValidationError(c *gin.Context, err error) {
	log.Printf("User validation failed: %v", err)

	if errors.Is(err, errors.New(USER_NOT_FOUND)) ||
		strings.Contains(strings.ToLower(err.Error()), USER_NOT_FOUND) {
		apiresponse.SendUnauthorized(c, "Invalid user")
		return
	}

	apiresponse.SendInternalError(c, "Authentication failed")
}

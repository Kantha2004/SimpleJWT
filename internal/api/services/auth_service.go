package services

import (
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	ValidateUserFromContext(c *gin.Context) (*models.AdminUser, bool)
}

type authService struct {
	userService UserService
}

func NewAuthService(userService UserService) AuthService {
	return &authService{
		userService: userService,
	}
}

func (s *authService) ValidateUserFromContext(c *gin.Context) (*models.AdminUser, bool) {
	// Extract user ID from context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return nil, false
	}

	// Validate user existence
	user, err := s.userService.ValidateUser(userID)
	if err != nil {
		return nil, false
	}

	return user, true
}

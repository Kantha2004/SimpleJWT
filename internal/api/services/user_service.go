// internal/services/user_service.go
package services

import (
	"fmt"
	"time"

	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/repositories"
)

type UserService interface {
	CreateUser(req models.CreateUser) (*models.CreateUserResponse, error)
	AuthenticateUser(req models.LoginRequest) (*models.LoginResponse, error)
	ValidateUser(userID uint) (*models.AdminUser, error)
}

type userService struct {
	userRepo   *repositories.UserRepository
	jwtService *auth.JWTService
}

func NewUserService(userRepo *repositories.UserRepository, jwtService *auth.JWTService) UserService {
	return &userService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *userService) CreateUser(req models.CreateUser) (*models.CreateUserResponse, error) {

	if exists, err := s.userRepo.UserNameExists(req.Username); err != nil {
		return nil, NewInternalError("Failed to validate username", err)
	} else if exists {
		return nil, NewConflictError("Username already exists")
	}

	if exists, err := s.userRepo.EmailExists(req.Email); err != nil {
		return nil, NewInternalError("Failed to validate email", err)
	} else if exists {
		return nil, NewConflictError("Email already exists")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, NewInternalError("Failed to process password", err)
	}

	user := &models.AdminUser{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	userID, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, NewInternalError("Failed to create user", err)
	}

	return &models.CreateUserResponse{
		UserID:   userID,
		Username: req.Username,
		Email:    req.Email,
	}, nil
}

func (s *userService) AuthenticateUser(req models.LoginRequest) (*models.LoginResponse, error) {

	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, NewInternalError("Authentication failed", err)
	}

	if user == nil || !auth.CheckPassword(user.PasswordHash, req.Password) {
		return nil, NewUnauthorizedError("Invalid username or password")
	}

	token, err := s.jwtService.CreateToken(user.ID)
	if err != nil {
		return nil, NewInternalError("Authentication failed", err)
	}

	expiresAt := time.Now().Add(time.Hour * 24)
	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *userService) ValidateUser(userID uint) (*models.AdminUser, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if user == nil {
		return nil, NewNotFoundError("User not found")
	}

	return user, nil
}

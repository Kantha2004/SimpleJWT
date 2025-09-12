package services

import (
	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/repositories"
	"github.com/Kantha2004/SimpleJWT/internal/utils"
)

// ClientUserService defines the business logic for client user operations
type ClientUserService interface {
	CreateUser(req models.CreateClientUser) (*models.ClientUser, error)
	AuthenticateUser(req models.ClientUserLoginRequest) (*models.LoginResponse, error)
}

type clientUserService struct {
	clientRepo     repositories.ClientRepository
	clientUserRepo repositories.ClientUserRepository
	jwtService     auth.JWTService
}

// NewClientUserService creates a new instance of ClientUserService
func NewClientUserService(
	clientRepo repositories.ClientRepository,
	clientUserRepo repositories.ClientUserRepository,
	jwtService auth.JWTService,
) ClientUserService {
	return &clientUserService{
		clientRepo:     clientRepo,
		clientUserRepo: clientUserRepo,
		jwtService:     jwtService,
	}
}

// CreateUser handles the business logic for creating a new client user
func (s *clientUserService) CreateUser(req models.CreateClientUser) (*models.ClientUser, error) {
	// Validate client exists
	client, err := s.clientRepo.GetClientById(req.ClientID)
	if err != nil {
		return nil, NewInternalError("Failed to validate client", err)
	}
	if client == nil {
		return nil, NewNotFoundError("Client not found")
	}

	// Check for existing username
	if exists, err := s.clientUserRepo.ClientUserNameExists(req.Username); err != nil {
		return nil, NewInternalError("Failed to validate username", err)
	} else if exists {
		return nil, NewConflictError("Username already exists")
	}

	// Check for existing email
	if exists, err := s.clientUserRepo.ClientUserEmailExists(req.Email); err != nil {
		return nil, NewInternalError("Failed to validate email", err)
	} else if exists {
		return nil, NewConflictError("Email already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, NewInternalError("Failed to process password", err)
	}

	// Create user model
	userModel := &models.ClientUser{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	// Save user
	user, err := s.clientUserRepo.CreateClientUser(userModel)
	if err != nil {
		return nil, NewInternalError("Failed to create user", err)
	}

	return user, nil
}

// AuthenticateUser handles user authentication and token generation
func (s *clientUserService) AuthenticateUser(req models.ClientUserLoginRequest) (*models.LoginResponse, error) {
	// Validate client
	client, err := s.clientRepo.GetClientBySecret(req.ClientSecret)
	if err != nil {
		return nil, NewInternalError("Failed to validate client", err)
	}
	if client == nil {
		return nil, NewNotFoundError("Client not found")
	}

	// Get user by username
	user, err := s.clientUserRepo.GetClientUserByUsername(req.Username)
	if err != nil {
		return nil, NewInternalError("Failed to validate user", err)
	}
	if user == nil {
		return nil, NewUnauthorizedError("Invalid username or password")
	}

	// Verify password
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		return nil, NewUnauthorizedError("Invalid username or password")
	}

	// Generate token
	token, err := s.jwtService.CreateToken(user.ID)
	if err != nil {
		return nil, NewInternalError("Unable to generate authentication token", err)
	}

	// Build response
	userInfo := models.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	response := &models.LoginResponse{
		User:      userInfo,
		Token:     token,
		ExpiresAt: utils.GetExpiryTime(),
	}

	return response, nil
}

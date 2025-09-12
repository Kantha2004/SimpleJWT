package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/repositories"
	"github.com/Kantha2004/SimpleJWT/internal/utils"
	"github.com/gin-gonic/gin"
)

// CreateClientUser godoc
// @Summary Create a new ClientUser
// @Description Create and add a new user account to the client
// @Tags Client
// @Accept json
// @Produce json
// @Param user body models.CreateClientUser true "User creation data"
// @Success 201 {object} apiresponse.SuccessResponse{data=models.ClientUser} "User created successfully"
// @Failure 400 {object} apiresponse.ErrorResponse "Bad request - validation error"
// @Failure 409 {object} apiresponse.ErrorResponse "Conflict - username or email already exists"
// @Failure 500 {object} apiresponse.ErrorResponse "Internal server error"
// @Router /protected/createClientUser [post]
// @Security BearerAuth
func (d *Dependencies) CreateClientUser(c *gin.Context) {
	var req models.CreateClientUser

	if verified := utils.VerifyRequestModel(c, &req); !verified {
		return
	}

	clientRepo := repositories.NewClientRepository(d.DB)

	_, ok := d.ValidateUserFromContext(c)

	if !ok {
		return
	}

	client, err := clientRepo.GetClientById(req.ClientID)

	if err != nil {
		apiresponse.SendInternalError(c, "Internal Server Error")
		return
	}

	if client == nil {
		apiresponse.SendValidationError(c, errors.New("Client Not Found"))
		return
	}

	clientUserRepo := repositories.NewClientUserRepository(d.DB, client.SchemaName)

	// Check if username exists
	if exists, err := clientUserRepo.ClientUserNameExists(req.Username); err != nil {
		log.Printf("Error checking username existence: %v", err)
		apiresponse.SendInternalError(c, "Failed to validate username")
		return
	} else if exists {
		apiresponse.SendConflict(c, "Username already exists")
		return
	}

	// Check if email exists
	if exists, err := clientUserRepo.ClientUserEmailExists(req.Email); err != nil {
		log.Printf("Error checking email existence: %v", err)
		apiresponse.SendInternalError(c, "Failed to validate email")
		return
	} else if exists {
		apiresponse.SendConflict(c, "Email already exists")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		apiresponse.SendInternalError(c, "Failed to process password")
		return
	}

	// Create ClientUser
	clientUserModel := &models.ClientUser{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	clientUser, err := clientUserRepo.CreateClientUser(clientUserModel)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		apiresponse.SendInternalError(c, "Failed to create user")
		return
	}

	apiresponse.SendSuccess(c, http.StatusCreated, clientUser, "User successfully added")

}

// ClientUserLogin godoc
// @Summary Authenticate a client user
// @Description Authenticate and login a client user with username/password
// @Tags Client
// @Accept json
// @Produce json
// @Param user body models.ClientUserLoginRequest true "User login data"
// @Success 200 {object} apiresponse.SuccessResponse{data=models.LoginResponse} "Login successful"
// @Failure 400 {object} apiresponse.ErrorResponse "Bad request - validation error"
// @Failure 401 {object} apiresponse.ErrorResponse "Unauthorized - invalid credentials"
// @Failure 404 {object} apiresponse.ErrorResponse "Client not found"
// @Failure 500 {object} apiresponse.ErrorResponse "Internal server error"
// @Router /client/userlogin [post]
// @Security BearerAuth
func (d *Dependencies) ClientUserLogin(c *gin.Context) {
	var req models.ClientUserLoginRequest

	// Validate request model
	if ok := utils.VerifyRequestModel(c, &req); !ok {
		return
	}

	// Validate client exists
	client, err := d.validateClient(c, req.ClientSecret)
	if err != nil || client == nil {
		return
	}

	// Authenticate user
	user, err := d.authenticateUser(c, client.SchemaName, req.Username, req.Password)
	if err != nil || user == nil {
		return
	}

	// Generate JWT token
	token, err := d.generateUserToken(c, user.ID)
	if err != nil {
		return
	}

	// Build and send response
	d.sendLoginResponse(c, user, token)
}

// validateClient validates that the client exists
func (d *Dependencies) validateClient(c *gin.Context, client_secret string) (*models.Client, error) {
	clientRepo := repositories.NewClientRepository(d.DB)
	client, err := clientRepo.GetClientBySecret(client_secret)

	if err != nil {
		log.Printf("Error retrieving client %s: %v", client_secret, err)
		apiresponse.SendInternalError(c, "Internal Server Error")
		return nil, err
	}

	if client == nil {
		apiresponse.SendNotFound(c, "Client not found")
		return nil, fmt.Errorf("client not found: %s", client_secret)
	}

	return client, nil
}

// authenticateUser validates user credentials
func (d *Dependencies) authenticateUser(c *gin.Context, schemaName, username, password string) (*models.ClientUser, error) {
	clientUserRepo := repositories.NewClientUserRepository(d.DB, schemaName)

	user, err := clientUserRepo.GetClientUserByUsername(username)
	if err != nil {
		log.Printf("Error retrieving user %s from schema %s: %v", username, schemaName, err)
		apiresponse.SendInternalError(c, "Failed to validate user")
		return nil, err
	}

	if user == nil {
		apiresponse.SendUnauthorized(c, "Invalid username or password")
		return nil, fmt.Errorf("user not found: %s", username)
	}

	if !auth.CheckPassword(user.PasswordHash, password) {
		apiresponse.SendUnauthorized(c, "Invalid username or password")
		return nil, fmt.Errorf("invalid password for user: %s", username)
	}

	return user, nil
}

// generateUserToken creates a JWT token for the authenticated user
func (d *Dependencies) generateUserToken(c *gin.Context, userID uint) (string, error) {
	token, err := d.jwtService.CreateToken(userID)
	if err != nil {
		log.Printf("Error generating token for user ID %s: %v", userID, err)
		apiresponse.SendInternalError(c, "Unable to generate authentication token")
		return "", err
	}

	return token, nil
}

// sendLoginResponse builds and sends the successful login response
func (d *Dependencies) sendLoginResponse(c *gin.Context, user *models.ClientUser, token string) {
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

	apiresponse.SendSuccess(c, http.StatusOK, response, "Login successful")
}

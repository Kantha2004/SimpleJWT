package handlers

import (
	"log"
	"net/http"
	"time"

	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/db"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User creation data"
// @Success 201 {object} apiresponse.SuccessResponse{data=models.CreateUserResponse} "User created successfully"
// @Failure 400 {object} apiresponse.ErrorResponse "Bad request - validation error"
// @Failure 409 {object} apiresponse.ErrorResponse "Conflict - username or email already exists"
// @Failure 500 {object} apiresponse.ErrorResponse "Internal server error"
// @Router /createUser [post]
func (d *Dependencies) CreateUser(c *gin.Context) {
	var req models.CreateUser

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		apiresponse.SendValidationError(c, err)
		return
	}

	userRepo := db.NewUserRepository(d.DB)

	// Check if username exists
	if exists, err := userRepo.UserNameExists(req.Username); err != nil {
		log.Printf("Error checking username existence: %v", err)
		apiresponse.SendInternalError(c, "Failed to validate username")
		return
	} else if exists {
		apiresponse.SendConflict(c, "Username already exists")
		return
	}

	// Check if email exists
	if exists, err := userRepo.EmailExists(req.Email); err != nil {
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

	// Create user
	user := &models.AdminUser{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	userId, err := userRepo.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		apiresponse.SendInternalError(c, "Failed to create user")
		return
	}

	// Prepare response data
	responseData := models.CreateUserResponse{
		UserID:   userId,
		Username: req.Username,
		Email:    req.Email,
	}

	apiresponse.SendSuccess(c, http.StatusCreated, responseData, "User created successfully")
}

// Login godoc
// @Summary User login
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} apiresponse.SuccessResponse{data=models.LoginResponse} "Login successful"
// @Failure 400 {object} apiresponse.ErrorResponse "Bad request"
// @Failure 401 {object} apiresponse.ErrorResponse "Invalid credentials"
// @Failure 500 {object} apiresponse.ErrorResponse "Internal server error"
// @Router /login [post]
func (d *Dependencies) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		apiresponse.SendValidationError(c, err)
		return
	}

	userRepo := db.NewUserRepository(d.DB)

	// Get user by username
	user, err := userRepo.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		apiresponse.SendInternalError(c, "Authentication failed")
		return
	}

	// Check user existence and password in one step for security
	if user == nil || !auth.CheckPassword(user.PasswordHash, req.Password) {
		apiresponse.SendUnauthorized(c, "Invalid username or password")
		return
	}

	// Generate JWT token
	token, err := d.jwtService.CreateToken(user.ID)
	if err != nil {
		log.Printf("Error creating token: %v", err)
		apiresponse.SendInternalError(c, "Authentication failed")
		return
	}

	// Prepare response
	expiresAt := time.Now().Add(time.Hour * 24) // Match your JWT expiration
	responseData := models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	apiresponse.SendSuccess(c, http.StatusOK, responseData, "Login successful")
}

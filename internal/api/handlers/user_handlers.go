package handlers

import (
	"net/http"

	"github.com/Kantha2004/SimpleJWT/internal/api/services"
	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
	authService services.AuthService
}

func NewUserHandler(userService services.UserService, authService services.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

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
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUser
	if !utils.VerifyRequestModel(c, &req) {
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		handleServiceError(c, err, "Failed to create user")
		return
	}

	apiresponse.SendSuccess(c, http.StatusCreated, user, "User created successfully")
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
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if !utils.VerifyRequestModel(c, &req) {
		return
	}

	response, err := h.userService.AuthenticateUser(req)
	if err != nil {
		handleServiceError(c, err, "Authentication failed")
		return
	}

	apiresponse.SendSuccess(c, http.StatusOK, response, "Login successful")
}

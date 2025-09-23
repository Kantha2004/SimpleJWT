package handlers

import (
	"net/http"

	"github.com/Kantha2004/SimpleJWT/internal/api/services"
	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/db"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/repositories"
	"github.com/Kantha2004/SimpleJWT/internal/utils"
	"github.com/gin-gonic/gin"
)

type ClientUserHandler struct {
	DB         *db.Database
	jwtService *auth.JWTService
}

func NewClientUserHandler(db *db.Database, authService *auth.JWTService) *ClientUserHandler {
	return &ClientUserHandler{
		DB:         db,
		jwtService: authService,
	}
}

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
func (d *ClientUserHandler) CreateClientUser(c *gin.Context) {
	var req models.CreateClientUser

	if ok := utils.VerifyRequestModel(c, &req); !ok {
		return
	}

	clientRepo := repositories.NewClientRepository(d.DB)

	client, err := clientRepo.GetClientById(req.ClientID)

	if err != nil {
		apiresponse.SendInternalError(c, err.Error())
	}

	if client == nil {
		apiresponse.SendNotFound(c, "Client not found")
	}

	service := services.NewClientUserService(
		clientRepo,
		repositories.NewClientUserRepository(d.DB, client.SchemaName),
		d.jwtService,
	)

	user, err := service.CreateUser(req)
	if err != nil {
		handleServiceError(c, err, "")
		return
	}

	apiresponse.SendSuccess(c, http.StatusCreated, user, "User successfully added")
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
func (d *ClientUserHandler) ClientUserLogin(c *gin.Context) {
	var req models.ClientUserLoginRequest

	if ok := utils.VerifyRequestModel(c, &req); !ok {
		return
	}

	clientRepo := repositories.NewClientRepository(d.DB)

	client, err := clientRepo.GetClientBySecret(req.ClientSecret)

	if err != nil {
		apiresponse.SendInternalError(c, err.Error())
	}

	if client == nil {
		apiresponse.SendNotFound(c, "Client not found")
	}

	service := services.NewClientUserService(
		clientRepo,
		repositories.NewClientUserRepository(d.DB, client.SchemaName),
		d.jwtService,
	)

	loginResp, err := service.AuthenticateUser(req)

	if err != nil {
		handleServiceError(c, err, "")
		return
	}

	apiresponse.SendSuccess(c, http.StatusOK, loginResp, "Login successful")
}

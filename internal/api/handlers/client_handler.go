// internal/handlers/client_handler.go
package handlers

import (
	"net/http"

	"github.com/Kantha2004/SimpleJWT/internal/api/services"
	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/utils"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientService services.ClientService
	authService   services.AuthService
}

func NewClientHandler(clientService services.ClientService, authService services.AuthService) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
		authService:   authService,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account in the system
// @Tags Client
// @Accept json
// @Produce json
// @Param user body models.CreateClient true "Client creation data"
// @Success 201 {object} apiresponse.SuccessResponse{data=models.CreateClientReponse} "Client created successfully"
// @Failure 400 {object} apiresponse.ErrorResponse "Bad request - validation error"
// @Failure 409 {object} apiresponse.ErrorResponse "Conflict - client name already exists"
// @Failure 500 {object} apiresponse.ErrorResponse "Internal server error"
// @Router /protected/createClient [post]
// @Security BearerAuth
func (h *ClientHandler) CreateClient(c *gin.Context) {
	var req models.CreateClient
	if !utils.VerifyRequestModel(c, &req) {
		return
	}

	user, ok := h.authService.ValidateUserFromContext(c)
	if !ok {
		return
	}

	response, err := h.clientService.CreateClient(req, user)
	if err != nil {
		handleServiceError(c, err, "Failed to create client")
		return
	}

	apiresponse.SendSuccess(c, http.StatusCreated, response, "Client created successfully")
}

// GetAllClients godoc
// @Summary Get all clients associated with the user
// @Description Retrieve all clients belonging to the authenticated user
// @Tags Client
// @Accept json
// @Produce json
// @Success 200 {object} apiresponse.SuccessResponse{data=[]models.Client} "Successfully retrieved all clients"
// @Failure 400 {object} apiresponse.ErrorResponse "Bad request - validation error"
// @Failure 401 {object} apiresponse.ErrorResponse "Unauthorized - invalid user"
// @Failure 404 {object} apiresponse.ErrorResponse "No clients found"
// @Failure 500 {object} apiresponse.ErrorResponse "Internal server error"
// @Router /protected/getAllClients [get]
// @Security BearerAuth
func (h *ClientHandler) GetAllClients(c *gin.Context) {
	user, ok := h.authService.ValidateUserFromContext(c)
	if !ok {
		return
	}

	clients, err := h.clientService.GetAllClientsByUser(user.ID)
	if err != nil {
		handleServiceError(c, err, "Error fetching clients")
		return
	}

	apiresponse.SendSuccess(c, http.StatusOK, clients, "Successfully retrieved all clients")
}

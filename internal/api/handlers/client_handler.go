package handlers

import (
	"fmt"
	"log"
	"net/http"

	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/db"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/utils"
	"github.com/gin-gonic/gin"
)

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
func (d *Dependencies) CreateClient(c *gin.Context) {
	var req models.CreateClient

	// Validate request body
	if verified := utils.VerifyRequestModel(c, &req); !verified {
		return
	}

	user, ok := d.ValidateUserFromContext(c)

	if !ok {
		return
	}

	userID := user.ID

	// Check if client name already exists for this user
	clientRepo := db.NewClientRepository(d.DB)
	exists, err := clientRepo.GetClientByNameForUser(req.ClientName, userID)
	if err != nil {
		log.Printf("Error checking client existence: %v", err)
		apiresponse.SendInternalError(c, "Error validating client name")
		return
	}

	if exists != nil {
		log.Printf("Client name '%s' already exists for user ID %d", req.ClientName, userID)
		apiresponse.SendAlreadyExistError(c, "Client name already exists")
		return
	}

	// Create client
	schemaName := fmt.Sprintf("%s_%s_client", user.Username, req.ClientName)
	client := &models.Client{
		ClientName: req.ClientName,
		UserID:     userID,
		SchemaName: schemaName,
	}

	clientSecret, err := clientRepo.CreateClient(client)
	if err != nil {
		log.Printf("Error creating client: %v", err)
		apiresponse.SendInternalError(c, "Failed to create client")
		return
	}

	// Create client schema and migrate tables
	if err := d.DB.CreateClientSchema(schemaName); err != nil {
		log.Printf("Error creating client schema: %v", err)
		apiresponse.SendInternalError(c, "Failed to initialize client schema")
		return
	}

	if err := d.DB.MigrateClientTables(schemaName); err != nil {
		log.Printf("Error migrating client tables: %v", err)
		apiresponse.SendInternalError(c, "Failed to migrate client tables")
		return
	}

	response := models.CreateClientReponse{
		ClientSecret: clientSecret,
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
func (d *Dependencies) GetAllClients(c *gin.Context) {

	user, ok := d.ValidateUserFromContext(c)

	if !ok {
		return
	}

	userID := user.ID

	// Fetch all clients for the user
	clientRepo := db.NewClientRepository(d.DB)
	clients, err := clientRepo.GetAllClientsByUserId(userID)

	if err != nil {
		log.Printf("Error fetching clients for user ID %d: %v", userID, err)
		apiresponse.SendInternalError(c, "Error fetching clients")
		return
	}

	// Handle empty results
	if len(clients) == 0 {
		apiresponse.SendError(c, http.StatusNotFound, "No clients found")
		return
	}

	apiresponse.SendSuccess(c, http.StatusOK, clients, "Successfully retrieved all clients")
}

package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/db"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account in the system
// @Tags Client
// @Accept json
// @Produce json
// @Param user body models.CreateClient true "Client creation data"
// @Success 201 {object} apiresponse.SuccessResponse{data=models.CreateClinetReponse} "Client created successfully"
// @Failure 400 {object} apiresponse.ErrorResponse "Bad request - validation error"
// @Failure 409 {object} apiresponse.ErrorResponse "Conflict - client name already exists"
// @Failure 500 {object} apiresponse.ErrorResponse "Internal server error"
// @Router /protected/createClient [post]
// @Security     BearerAuth
func (d *Dependencies) CreateClient(c *gin.Context) {
	var req models.CreateClient
	userIDVal, ok := c.Get("user_id")

	if !ok {
		apiresponse.SendValidationError(c, errors.New("Unable to find userID from the request"))
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		apiresponse.SendValidationError(c, err)
		return
	}

	fmt.Println("userIDVal", userIDVal)
	fmt.Println("userIDValint", (userIDVal).(int))
	userID := uint((userIDVal).(int))
	fmt.Println("userIDuint", userID)

	userRepo := db.NewUserRepository(d.DB)
	clientRepo := db.NewClientRepository(d.DB)

	user, err := userRepo.GetUserByID(userID)

	if err != nil {
		log.Printf("Error fetching user: %v", err)
		apiresponse.SendInternalError(c, "Authentication failed")
		return
	}

	if user == nil {
		apiresponse.SendUnauthorized(c, "Invalid user")
	}

	exists, err := clientRepo.GetClientByNameForUser(req.ClientName, userID)

	if err != nil {
		log.Printf("Error fetching user: %v", err)
		apiresponse.SendInternalError(c, "Authentication failed")
		return
	}

	if exists != nil {
		log.Print("Client name already exist for this user")
		apiresponse.SendAlreadyExistError(c, "Client name already exists")
		return
	}

	client := &models.Client{
		ClientName: req.ClientName,
		UserID:     userID,
	}

	clientSecret, err := clientRepo.CreateClient(client)

	if err != nil {
		apiresponse.SendInternalError(c, err.Error())
	}

	response := models.CreateClinetReponse{
		ClientSecret: clientSecret,
	}

	schemaName := fmt.Sprintf("%s_%s_client", user.Username, client.ClientName)

	d.DB.CreateClientSchema(schemaName)
	d.DB.MigrateClientTables(schemaName)

	apiresponse.SendSuccess(c, http.StatusCreated, response, "Client created successfully")
}

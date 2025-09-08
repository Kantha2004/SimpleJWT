package handlers

import (
	"errors"
	"log"
	"net/http"

	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/db"
	"github.com/Kantha2004/SimpleJWT/internal/models"
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

	clientRepo := db.NewClientRepository(d.DB)

	_, ok := d.ValidateUserFromContext(c)

	if !ok {
		return
	}

	client, err := clientRepo.GetClientId(req.ClientID)

	if err != nil {
		apiresponse.SendInternalError(c, "Internal Server Error")
		return
	}

	if client == nil {
		apiresponse.SendValidationError(c, errors.New("Client Not Found"))
		return
	}

	clientUserRepo := db.NewClientUserRepository(d.DB, client.SchemaName)

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

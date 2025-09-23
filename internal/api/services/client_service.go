// internal/services/client_service.go
package services

import (
	"fmt"

	"github.com/Kantha2004/SimpleJWT/internal/db"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"github.com/Kantha2004/SimpleJWT/internal/repositories"
)

type ClientService interface {
	CreateClient(req models.CreateClient, user *models.AdminUser) (*models.CreateClientReponse, error)
	GetAllClientsByUser(userID uint) ([]*models.Client, error)
}

type clientService struct {
	clientRepo *repositories.ClientRepository
	db         *db.Database
}

func NewClientService(
	clientRepo *repositories.ClientRepository,
	db *db.Database,
) ClientService {
	return &clientService{
		clientRepo: clientRepo,
		db:         db,
	}
}

func (s *clientService) CreateClient(req models.CreateClient, user *models.AdminUser) (*models.CreateClientReponse, error) {
	existing, err := s.clientRepo.GetClientByNameForUser(req.ClientName, user.ID)
	if err != nil {
		return nil, NewInternalError("Error validating client name", err)
	}
	if existing != nil {
		return nil, NewConflictError("Client name already exists")
	}

	schemaName := fmt.Sprintf("%s_%s_client", user.Username, req.ClientName)

	client := &models.Client{
		ClientName: req.ClientName,
		UserID:     user.ID,
		SchemaName: schemaName,
	}

	clientSecret, err := s.clientRepo.CreateClient(client)
	if err != nil {
		return nil, NewInternalError("Failed to create client", err)
	}

	if err := s.db.CreateClientSchema(schemaName); err != nil {
		return nil, NewInternalError("Failed to initialize client schema", err)
	}

	if err := s.db.MigrateClientTables(schemaName); err != nil {
		return nil, NewInternalError("Failed to migrate client tables", err)
	}

	return &models.CreateClientReponse{
		ClientSecret: clientSecret,
	}, nil
}

func (s *clientService) GetAllClientsByUser(userID uint) ([]*models.Client, error) {
	clients, err := s.clientRepo.GetAllClientsByUserId(userID)
	if err != nil {
		return nil, NewInternalError("Error fetching clients", err)
	}

	if len(clients) == 0 {
		return nil, NewNotFoundError("No clients found")
	}

	return clients, nil
}

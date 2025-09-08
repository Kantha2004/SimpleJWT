package db

import (
	"errors"

	"github.com/Kantha2004/SimpleJWT/internal/models"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *Database
}

func NewClientRepository(db *Database) *ClientRepository {
	return &ClientRepository{db: db}
}

func (cr *ClientRepository) CreateClient(client *models.Client) (string, error) {
	result := cr.db.DB.Create(client)
	return client.ClientSecret, result.Error
}

func (cr *ClientRepository) GetClientByNameForUser(clientName string, userID uint) (*models.Client, error) {
	var client models.Client

	result := cr.db.DB.Where("client_name = ? AND user_id = ?", clientName, userID).First(&client)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &client, nil
}

func (cr *ClientRepository) GetClientByUserId(userID uint) (*models.Client, error) {
	var client models.Client

	result := cr.db.DB.Where("user_id = ?", userID).First(&client)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &client, nil
}

func (cr *ClientRepository) GetAllClientsByUserId(userID uint) (*[]models.Client, error) {
	var clients []models.Client

	result := cr.db.DB.Where("user_id = ?", userID).Find(&clients)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &clients, nil
}

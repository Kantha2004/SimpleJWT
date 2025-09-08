package db

import (
	"errors"

	"github.com/Kantha2004/SimpleJWT/internal/models"
	"gorm.io/gorm"
)

type ClientUserRepository struct {
	db *gorm.DB
}

func NewClientUserRepository(db *Database, schemaName string) *ClientUserRepository {
	return &ClientUserRepository{db: db.TableWithSchema(schemaName, CLIENT_USER_TABLE)}
}

func (cur *ClientUserRepository) CreateClientUser(user *models.ClientUser) (*models.ClientUser, error) {
	result := cur.db.Create(user)
	return user, result.Error
}

func (cur *ClientUserRepository) GetClientUserByID(id uint) (*models.ClientUser, error) {
	var user models.ClientUser
	result := cur.db.Find(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (cur *ClientUserRepository) GetClientAllUser() (*[]models.ClientUser, error) {
	var users []models.ClientUser
	result := cur.db.Find(&users)
	return &users, result.Error
}

func (cur *ClientUserRepository) GetClientUserByEmail(email string) (*models.ClientUser, error) {
	var user models.ClientUser
	result := cur.db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (cur *ClientUserRepository) GetClientUserByUsername(username string) (*models.ClientUser, error) {
	var user models.ClientUser
	result := cur.db.Where("username = ?", username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (cur *ClientUserRepository) UpdateClientUser(user *models.ClientUser) error {
	result := cur.db.Save(user)
	return result.Error
}

func (cur *ClientUserRepository) ClientUserExists(username, email string) (bool, error) {
	var count int64
	result := cur.db.Model(&models.ClientUser{}).
		Where("username = ? AND email = ?", username, email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (cur *ClientUserRepository) ClientUserEmailExists(email string) (bool, error) {
	var count int64
	result := cur.db.Model(&models.ClientUser{}).Where("email = ?", email).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (cur *ClientUserRepository) ClientUserNameExists(username string) (bool, error) {
	var count int64
	result := cur.db.Model(&models.ClientUser{}).Where("username = ?", username).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

package db

import (
	"errors"

	"github.com/Kantha2004/SimpleJWT/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *Database
}

func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(user *models.AdminUser) (uint, error) {
	result := ur.db.DB.Create(user)
	return user.ID, result.Error
}

func (ur *UserRepository) GetUserByID(id uint) (*models.AdminUser, error) {
	var user models.AdminUser
	result := ur.db.DB.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.AdminUser, error) {
	var user models.AdminUser
	result := ur.db.DB.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*models.AdminUser, error) {
	var user models.AdminUser
	result := ur.db.DB.Where("username = ?", username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUser(user *models.AdminUser) error {
	result := ur.db.DB.Save(user)
	return result.Error
}

func (ur *UserRepository) UserExists(username, email string) (bool, error) {
	var count int64
	result := ur.db.DB.Model(&models.AdminUser{}).
		Where("username = ? AND email = ?", username, email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (ur *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	result := ur.db.DB.Model(&models.AdminUser{}).Where("email = ?", email).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (ur *UserRepository) UserNameExists(username string) (bool, error) {
	var count int64
	result := ur.db.DB.Model(&models.AdminUser{}).Where("username = ?", username).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

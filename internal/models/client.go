package models

import (
	"crypto/rand"
	"encoding/hex"

	"gorm.io/gorm"
)

type Client struct {
	ClientName   string    `json:"client_name" gorm:"unique;not null"`
	ClientSecret string    `json:"client_secret" gorm:"not null"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	User         AdminUser `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TableModel
}

func generateUniqueClientSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (c *Client) BeforeCreate(tx *gorm.DB) error {
	if c.ClientSecret == "" {
		secret, err := generateUniqueClientSecret()
		if err != nil {
			return err
		}
		c.ClientSecret = secret
	}
	return nil
}

type CreateClient struct {
	ClientName string `json:"client_name" binding:"required"`
}

type CreateClinetReponse struct {
	ClientSecret string `json:"client_secret" binding:"required"`
}

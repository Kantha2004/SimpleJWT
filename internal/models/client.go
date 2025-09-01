package models

import (
	"crypto/rand"
	"encoding/hex"

	"gorm.io/gorm"
)

type Client struct {
	ClientName   string `json:"client_name" gorm:"not null"`
	ClientSecret string `json:"client_secret" gorm:"unique;not null"`
	UserID       uint   `json:"user_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
	UserID     uint   `json:"user_id" binding:"required"`
}

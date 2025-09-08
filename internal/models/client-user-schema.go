package models

// ClientUser represents a user in the client
type ClientUser struct {
	AdminUser
}

// CreateClientUser represents the request payload for user creation
type CreateClientUser struct {
	CreateUser
	ClientID uint `json:"client_id" binding:"required" example:"1"`
}

package models

// ClientUser represents a user in the client
type ClientUser = AdminUser

// CreateClientUser represents the request payload for user creation
type CreateClientUser struct {
	CreateUser
	ClientID uint `json:"client_id" binding:"required" example:"0"`
}

type ClientRequest struct {
	ClientSecret string `json:"client_secret" binding:"required" example:"asdadsdasdsadasd"`
}

type ClientUserLoginRequest struct {
	LoginRequest
	ClientRequest
}

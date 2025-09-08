package models

import "time"

// AdminUser represents a user in the system
type AdminUser struct {
	ID           uint   `json:"id" gorm:"primaryKey" example:"1"`
	Username     string `json:"username" gorm:"unique;not null;size:50" example:"john_doe"`
	Email        string `json:"email" gorm:"unique;not null;size:100" example:"john@example.com"`
	PasswordHash string `json:"-" gorm:"not null;size:255"`
	TableModel
}

// CreateUser represents the request payload for user creation
type CreateUser struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"john_doe"`
	Password string `json:"password" binding:"required,min=8,max=100" example:"password123"`
	Email    string `json:"email" binding:"required,email,max=100" example:"john@example.com"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents the response for successful login
type LoginResponse struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time `json:"expiresAt" example:"2023-01-02T00:00:00Z"`
	User      UserInfo  `json:"user"`
}

// UserInfo represents public user information
type UserInfo struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" example:"john@example.com"`
}

// CreateUserResponse represents the response for user creation
type CreateUserResponse struct {
	UserID   uint   `json:"userId" example:"1"`
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" example:"john@example.com"`
}

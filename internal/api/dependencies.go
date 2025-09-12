package api

import (
	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/db"
)

// Dependencies holds all the required dependencies for the application
type Dependencies struct {
	DB         *db.Database
	JWTService *auth.JWTService
}

// NewDependencies creates a new Dependencies instance
func NewDependencies(db *db.Database, jwtService *auth.JWTService) *Dependencies {
	return &Dependencies{
		DB:         db,
		JWTService: jwtService,
	}
}

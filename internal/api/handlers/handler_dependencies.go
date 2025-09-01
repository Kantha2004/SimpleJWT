package handlers

import (
	"github.com/Kantha2004/SimpleJWT/internal/auth"
	"github.com/Kantha2004/SimpleJWT/internal/db"
)

type Dependencies struct {
	DB         *db.Database
	jwtService *auth.JWTService
}

func NewDependencies(db *db.Database, jwt *auth.JWTService) *Dependencies {
	return &Dependencies{
		DB:         db,
		jwtService: jwt,
	}
}

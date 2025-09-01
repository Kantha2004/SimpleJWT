package api

import (
	"github.com/Kantha2004/SimpleJWT/internal/auth"
)

type Dependencies struct {
	JWTService *auth.JWTService
}

func NewDependencies(jwtService *auth.JWTService) *Dependencies {
	return &Dependencies{
		JWTService: jwtService,
	}
}

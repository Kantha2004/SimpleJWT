package handlers

import (
	"github.com/Kantha2004/SimpleJWT/internal/db"
)

type Dependencies struct {
	DB *db.Database
}

func NewDependencies(db *db.Database) *Dependencies {
	return &Dependencies{
		DB: db,
	}
}

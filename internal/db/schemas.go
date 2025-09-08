package db

import (
	"fmt"

	"github.com/Kantha2004/SimpleJWT/internal/models"
)

// Create a schema for the client
func (db *Database) CreateClientSchema(schemaName string) error {
	return db.DB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)).Error
}

func (db *Database) MigrateClientTables(schemaName string) error {
	userTable := fmt.Sprintf("%s.users", schemaName)
	configTable := fmt.Sprintf("%s.configs", schemaName) // Avoid reserved word conflict

	if err := db.DB.Table(userTable).AutoMigrate(&models.ClientUser{}); err != nil {
		return err
	}

	if err := db.DB.Table(configTable).AutoMigrate(&models.ClientConfig{}); err != nil {
		return err
	}

	return nil
}

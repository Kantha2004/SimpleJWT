package db

import (
	"fmt"

	"github.com/Kantha2004/SimpleJWT/internal/models"
)

const (
	CLIENT_USER_TABLE   = "users"
	CLIENT_CONFIG_TABLE = "configs"
)

// Create a schema for the client
func (db *Database) CreateClientSchema(schemaName string) error {
	return db.DB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)).Error
}

func (db *Database) MigrateClientTables(schemaName string) error {
	userTable := fmt.Sprintf("%s.%s", schemaName, CLIENT_USER_TABLE)
	configTable := fmt.Sprintf("%s.%s", schemaName, CLIENT_CONFIG_TABLE)

	if err := db.DB.Table(userTable).AutoMigrate(&models.ClientUser{}); err != nil {
		return err
	}

	if err := db.DB.Table(configTable).AutoMigrate(&models.ClientConfig{}); err != nil {
		return err
	}

	return nil
}

// internal/services/database_service.go
package services

type DatabaseService interface {
	CreateClientSchema(schemaName string) error
	MigrateClientTables(schemaName string) error
}

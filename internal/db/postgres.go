package db

import (
	"fmt"
	"log"

	"github.com/Kantha2004/SimpleJWT/internal/config"
	"github.com/Kantha2004/SimpleJWT/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

// NewDatabase connects to PostgreSQL
func NewDatabase(cfg *config.DBConfig) (*Database, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// PostgreSQL DSN

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	database := &Database{DB: db}

	if err := database.migrate(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.DB
}

func (d *Database) migrate() error {
	err := d.DB.AutoMigrate(
		&models.AdminUser{},
		&models.Client{},
	)

	if err != nil {
		return err
	}

	log.Println("PostgreSQL database migration completed successfully")
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

package db

import (
	"log"

	"github.com/Kantha2004/SimpleJWT/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(sqlite.Open(dbPath), config)

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
		&models.User{},
		&models.Client{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

func (d *Database) Close() error {
	sqliteDB, err := d.DB.DB()

	if err != nil {
		return err
	}

	return sqliteDB.Close()
}

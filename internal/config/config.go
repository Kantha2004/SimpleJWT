package config

import (
	"log"
	"os"
	"strconv"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

type DBConfig struct {
	DBName     string
	DBPassword string
	DBHost     string
	DBPort     int
	DBUser     string
	SSLMode    string
}

type Config struct {
	JWTSecret   string
	DBConfig    DBConfig
	Port        string
	Environment Environment
}

func LoadConfig() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT secret is missing")
	}

	DBName := os.Getenv("DB_NAME")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBHost := os.Getenv("DB_HOST")
	DBPortStr := os.Getenv("DB_PORT")
	DBUser := os.Getenv("DB_USER")

	if DBPortStr == "" {
		log.Fatal("DB_PORT is missing")
	}
	DBPort, err := strconv.Atoi(DBPortStr)
	if err != nil {
		log.Fatal("DB_PORT must be a number")
	}

	return &Config{
		JWTSecret: jwtSecret,
		DBConfig: DBConfig{
			DBName:     DBName,
			DBPassword: DBPassword,
			DBHost:     DBHost,
			DBPort:     DBPort,
			DBUser:     DBUser,
			SSLMode:    getEnvWithDefault("SSL_MODE", "disable"),
		},
		Port:        getEnvWithDefault("PORT", "9000"),
		Environment: Environment(getEnvWithDefault("ENV", "development")),
	}
}

func getEnvWithDefault(key string, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

package config

import (
	"log"
	"os"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

type Config struct {
	JWTSecret   string
	DBPath      string
	Port        string
	Environment Environment
}

func LoadConfig() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT secret is missing")
	}

	return &Config{
		JWTSecret:   jwtSecret,
		DBPath:      getEnvWithDefault("DB_PATH", "simpleJWT.db"),
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

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DatabaseUrl string
	Port        string
	Migrations  string
	JWTSecret   string
}

// public
func LoadEnv() *EnvConfig {
	var err error = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}
	var envConfig = &EnvConfig{
		DatabaseUrl: os.Getenv("DB_URL"),
		Port:        os.Getenv("PORT"),
		Migrations:  os.Getenv("MIGRATIONS_DIR"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
	return envConfig
}

// private helper
func ValidateEnv(cfg *EnvConfig) error {
	// validation logic
	return nil
}

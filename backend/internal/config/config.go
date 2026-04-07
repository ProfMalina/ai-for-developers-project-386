package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	ServerPort  string
	DatabaseURL string
	Env         string
}

// LoadConfig loads configuration from environment and .env file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/booking_db?sslmode=disable"),
		Env:         getEnv("APP_ENV", "development"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

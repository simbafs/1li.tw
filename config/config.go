package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DBPath      string
	BotToken    string
	JWTSecret   string
	ServerPort  string
	Environment string
}

// LoadConfig loads configuration from environment variables or a .env file.
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	return &Config{
		DBPath:      getEnv("DB_PATH", "data/1li.db"),
		BotToken:    getEnv("BOT_TOKEN", ""),
		JWTSecret:   getEnv("JWT_SECRET", "a-very-secret-key"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}, nil
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

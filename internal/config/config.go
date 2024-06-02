package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("REDIS_GO_PORT")
	if port == "" {
		port = "3000" // default port if environment variable is not set
	}
	return &Config{
		Port: port,
	}, nil
}

package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/stevenxchung/redis.go/pkg/util"
)

type Config struct {
	Port string
}

func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		util.LogInfo("Environment file `.env` not found. Setting defaults...")
	}

	port := os.Getenv("REDIS_GO_PORT")
	if port == "" {
		// Default port if environment variable is not set
		port = "3000"
	}
	return &Config{
		Port: port,
	}, nil
}

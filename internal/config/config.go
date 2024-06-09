package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/stevenxchung/redis.go/pkg/util"
)

type Config struct {
	ServerPort string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		util.LogInfo("Environment file `.env` not found. Setting defaults...")
	}

	port := os.Getenv("REDIS_GO_SERVER_PORT")
	if port == "" {
		// Default server port if environment variable is not set
		port = "6379"
	}

	return &Config{
		ServerPort: port,
	}, nil
}

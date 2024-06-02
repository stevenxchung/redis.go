package main

import (
	"fmt"

	"github.com/stevenxchung/redis.go/internal/config"
	"github.com/stevenxchung/redis.go/internal/server"
	"github.com/stevenxchung/redis.go/pkg/util"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		util.LogError("Error loading configuration", err)
	}

	fmt.Println("Starting Redis Go...")
	s := server.NewServer(cfg)
	if err := s.Start(); err != nil {
		util.LogError("Failed to start server", err)
	}
}

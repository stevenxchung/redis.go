package main

import (
	"time"

	"github.com/stevenxchung/redis.go/external/client"
	"github.com/stevenxchung/redis.go/internal/config"
	"github.com/stevenxchung/redis.go/internal/server"
	"github.com/stevenxchung/redis.go/pkg/util"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		util.LogError("Error loading configuration", err)
	}

	go server.StartServer(cfg)
	time.Sleep(1 * time.Second)
	client.StartClient(cfg)
}

package server

import (
	"fmt"
	"net"

	"github.com/stevenxchung/redis.go/internal/config"
	"github.com/stevenxchung/redis.go/pkg/util"
)

type Server struct {
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) Start() error {
	server, err := net.Listen("tcp", ":"+s.config.ServerPort)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}
	defer server.Close()

	qh := NewQueryHandler()
	util.LogInfo(fmt.Sprintf("Server started on localhost:%s\n", s.config.ServerPort))

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go qh.queryHandler(conn)
	}
}

func StartServer(config *config.Config) {
	util.LogInfo("Starting redis.go...")
	s := NewServer(config)
	if err := s.Start(); err != nil {
		util.LogError("Failed to start server", err)
	}
}

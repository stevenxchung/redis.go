package server

import (
	"fmt"
	"net/http"

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
	http.HandleFunc("/graphql", NewQueryHandler().graphQLHandler)
	util.LogInfo(fmt.Sprintf("Server started on http://localhost:%s\n", s.config.Port))
	return http.ListenAndServe(":"+s.config.Port, nil)
}

package server

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/stevenxchung/redis.go/internal/config"
	"github.com/stevenxchung/redis.go/pkg/util"
)

type Server struct {
	config     *config.Config
	inMemoryDB map[string]string
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config:     cfg,
		inMemoryDB: make(map[string]string),
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/graphql", s.graphQLHandler)
	util.LogInfo(fmt.Sprintf("Server started on http://localhost:%s\n", s.config.Port))
	return http.ListenAndServe(":"+s.config.Port, nil)
}

func (s *Server) graphQLHandler(w http.ResponseWriter, r *http.Request) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    s.defineQuery(),
		Mutation: s.defineMutation(),
	})
	if err != nil {
		util.LogError("Failed to create schema", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	h.ServeHTTP(w, r)
}

func (s *Server) defineQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{"get": s.getQueryField()},
	})
}

func (s *Server) defineMutation() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"set": s.setMutationField(),
			"del": s.delMutationField(),
		},
	})
}

func (s *Server) getQueryField() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.String,
		Description: "Get the value for a given key",
		Args: graphql.FieldConfigArgument{
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			key := params.Args["key"].(string)
			return s.inMemoryDB[key], nil
		},
	}
}

func (s *Server) setMutationField() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.String,
		Description: "Set the value for a given key",
		Args: graphql.FieldConfigArgument{
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"value": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			key := params.Args["key"].(string)
			value := params.Args["value"].(string)
			s.inMemoryDB[key] = value
			return value, nil
		},
	}
}

func (s *Server) delMutationField() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.String,
		Description: "Delete a key-value pair for a given key",
		Args: graphql.FieldConfigArgument{
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			key := params.Args["key"].(string)
			delete(s.inMemoryDB, key)
			return key, nil
		},
	}
}

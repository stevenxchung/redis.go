package server

import (
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/stevenxchung/redis.go/pkg/util"
)

type ValueWithExpiration struct {
	Value   string
	Expires *time.Time
}

type QueryHandler struct {
	inMemoryDB map[string]ValueWithExpiration
}

func NewQueryHandler() *QueryHandler {
	return &QueryHandler{
		inMemoryDB: make(map[string]ValueWithExpiration),
	}
}

func (qh *QueryHandler) graphQLHandler(w http.ResponseWriter, r *http.Request) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    qh.defineQuery(),
		Mutation: qh.defineMutation(),
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

func (qh *QueryHandler) defineQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{"get": qh.getQueryField()},
	})
}

func (qh *QueryHandler) defineMutation() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"set": qh.setMutationField(),
			"del": qh.delMutationField(),
		},
	})
}

func (qh *QueryHandler) getQueryField() *graphql.Field {
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
			vwe, found := qh.inMemoryDB[key]
			if !found {
				return "data not found", nil
			}

			if vwe.Expires != nil && time.Now().After(*vwe.Expires) {
				delete(qh.inMemoryDB, key) // Remove expired value
				return "data has expired", nil
			}

			return vwe.Value, nil
		},
	}
}

func (qh *QueryHandler) setMutationField() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.String,
		Description: "Set the value for a given key, with an optional expiration time",
		Args: graphql.FieldConfigArgument{
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"value": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"expires": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			key := params.Args["key"].(string)
			value := params.Args["value"].(string)
			var ttl *time.Time

			if params.Args["expires"] != nil {
				expSeconds := params.Args["expires"].(int)
				expTime := time.Now().Add(time.Duration(expSeconds) * time.Second)
				ttl = &expTime
			}

			qh.inMemoryDB[key] = ValueWithExpiration{
				Value:   value,
				Expires: ttl,
			}

			return value, nil
		},
	}
}

func (qh *QueryHandler) delMutationField() *graphql.Field {
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
			delete(qh.inMemoryDB, key)
			return key, nil
		},
	}
}

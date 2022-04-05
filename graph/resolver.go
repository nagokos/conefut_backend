package graph

import (
	"database/sql"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nagokos/connefut_backend/graph/generated"
)

// This file will not be regenerated automatically.
//go:generate go run github.com/99designs/gqlgen

// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dbConnection *sql.DB
}

func NewSchema(db *sql.DB) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{db},
	})
}

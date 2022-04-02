package graph

import (
	"database/sql"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/graph/generated"
)

// This file will not be regenerated automatically.
//go:generate go run github.com/99designs/gqlgen

// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client       *ent.Client
	dbConnection *sql.DB
}

func NewSchema(client *ent.Client, db *sql.DB) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{client, db},
	})
}

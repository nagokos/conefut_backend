package graph

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/generated"
)

// This file will not be regenerated automatically.
//go:generate go run github.com/99designs/gqlgen

// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dbPool *pgxpool.Pool
}

func NewSchema(dbPool *pgxpool.Pool) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{dbPool: dbPool},
	})
}

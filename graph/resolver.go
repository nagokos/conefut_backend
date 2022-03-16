package graph

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/graph/generated"
)

// This file will not be regenerated automatically.
//go:generate go run github.com/99designs/gqlgen

// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client *ent.Client
}

func NewSchema(client *ent.Client) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{client},
	})
}

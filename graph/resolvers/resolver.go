package resolvers

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/vektah/gqlparser/v2/gqlerror"
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
		Directives: generated.DirectiveRoot{
			HasLoggedIn: func(ctx context.Context, _ interface{}, next graphql.Resolver) (res interface{}, err error) {
				viewer := user.GetViewer(ctx)
				if viewer == nil {
					logger.NewLogger().Error("user not loggedIn")
					return nil, gqlerror.Errorf("HAS_LOGGED_IN_ERROR")
				}
				return next(ctx)
			},
			HasEmailVerified: func(ctx context.Context, _ interface{}, next graphql.Resolver, status model.EmailVerificationStatus) (res interface{}, err error) {
				viewer := user.GetViewer(ctx)
				fmt.Println(viewer)
				if model.EmailVerificationStatus(strings.ToUpper(viewer.EmailVerificationStatus.String())) != status {
					logger.NewLogger().Error("access denined")
					return nil, gqlerror.Errorf("HAS_EMAIL_VERIFIED_ERROR")
				}
				return next(ctx)
			},
		},
	})
}

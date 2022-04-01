package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/config"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph"
	"github.com/nagokos/connefut_backend/graph/domain/user"
	"github.com/nagokos/connefut_backend/logger"
)

func init() {
	os.Setenv("TZ", "Asia/Tokyo")
}

func main() {
	port := config.Config.Port

	client := db.DatabaseConnection()
	defer client.Close()

	r := chi.NewRouter()
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
			AllowCredentials: true,
		}),
		auth.Middleware(client),
		auth.CookieMiddleWare(),
	)

	srv := handler.NewDefaultServer(graph.NewSchema(client))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)
	r.Route("/accounts", func(r chi.Router) {
		r.Route("/email_verification", func(r chi.Router) {
			r.Route("/{token}", func(r chi.Router) {
				r.Get("/", user.EmailVerification)
			})
		})
	})

	logger.NewLogger().Sugar().Infof("connect to http://localhost:%d/ for GraphQL playground", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
}

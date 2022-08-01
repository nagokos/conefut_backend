package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/config"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/loader"
	"github.com/nagokos/connefut_backend/graph/models/oauth"
	"github.com/nagokos/connefut_backend/graph/resolvers"
	"github.com/nagokos/connefut_backend/logger"
)

func init() {
	os.Setenv("TZ", "Asia/Tokyo")
	err := godotenv.Load(fmt.Sprintf("./env/%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
}

func main() {
	var err error

	port := config.Config.Port

	dbPool := db.DatabaseConnection()
	defer dbPool.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowCredentials: true,
		}),
	)
	r.Use(auth.Middleware(dbPool))
	r.Use(auth.CookieMiddleWare())

	loaders := loader.NewLoaders(dbPool)
	srv := handler.NewDefaultServer(resolvers.NewSchema(dbPool))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", loader.Middleware(loaders, srv))
	r.Route("/oauth", func(r chi.Router) {
		r.Route("/google", func(r chi.Router) {
			r.Get("/", oauth.AuthGoogleRedirect)
			r.Get("/callback", oauth.AuthGoogleCallback)
		})
		r.Route("/line", func(r chi.Router) {
			r.Get("/", oauth.AuthLineRedirect)
			r.Get("/callback", oauth.AuthLineCallback)
		})
	})

	logger.NewLogger().Sugar().Infof("connect to http://localhost:%d/ for GraphQL playground", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
}

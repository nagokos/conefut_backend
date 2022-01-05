package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/99designs/gqlgen/cmd"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/nagokos/connefut_backend/config"
	"github.com/nagokos/connefut_backend/ent/db"
	"github.com/nagokos/connefut_backend/graph"
)

const defaultPort = "8080"

func main() {
	port := config.Config.Port

	srv := handler.NewDefaultServer(graph.NewSchema(db.Client))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

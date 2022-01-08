package db

import (
	"context"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/migrate"
)

var Client *ent.Client

func init() {
	var err error
	Client, err = ent.Open("postgres", "host=db dbname=connefut_db port=5432 user=root password=password sslmode=disable")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	ctx := context.Background()

	f, err := os.Create("ent/db/migrate.sql")
	if err != nil {
		log.Fatalf("create migrate file: %v", err)
	}

	err = Client.Schema.WriteTo(
		ctx,
		f,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

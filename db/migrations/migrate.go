//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/logger"
)

func main() {
	client := db.DatabaseConnection()
	defer client.Close()

	ctx := context.Background()

	f, err := os.OpenFile("log/db_migrate.sql", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("create migrate file: %v", err))
	}

	err = client.Debug().Schema.WriteTo(
		ctx,
		f,
	)
	if err != nil {
		logger.Log.Fatal().Msg(fmt.Sprintf("failed creating schema resources: %v", err))
	}
}

package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/logger"
)

func DatabaseConnection() *pgxpool.Pool {
	uri := "postgres://root:password@db:5432/connefut_db?sslmode=disable"

	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		logger.NewLogger().Sugar().Fatalf("failed opening connection to postgres: %v", err)
	}

	err = client.Ping()
	if err != nil {
		logger.NewLogger().Sugar().Fatalf("pingError: %v", err)
	}

	db := sqldblogger.OpenDriver(
		"postgres://root:password@db:5432/connefut_db?sslmode=disable",
		client.Driver(),
		&logger.Logger{},
		sqldblogger.WithWrapResult(false),
		sqldblogger.WithSQLQueryAsMessage(true),
		sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
	)

	return db
}

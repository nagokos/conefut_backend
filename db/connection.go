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

	cfg.ConnConfig.Logger = zapadapter.NewLogger(logger.NewLogger())
	cfg.ConnConfig.LogLevel = pgx.LogLevelDebug
	dbPool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	return dbPool
}

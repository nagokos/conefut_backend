package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nagokos/connefut_backend/logger"
)

func DatabaseConnection() *pgxpool.Pool {
	err := godotenv.Load(fmt.Sprintf("./env/%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}

	uri := os.Getenv("DB_URL")
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

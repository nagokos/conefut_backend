package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/nagokos/connefut_backend/logger"
	sqldblogger "github.com/simukti/sqldb-logger"
)

func DatabaseConnection() *sql.DB {
	client, err := sql.Open("postgres", "host=db dbname=connefut_db port=5432 user=root password=password sslmode=disable")
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

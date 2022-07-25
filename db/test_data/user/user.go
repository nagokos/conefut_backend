package main

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mattn/go-gimei"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/logger"
)

func main() {
	gofakeit.Seed(0)
	cmd := `
	  INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, $3, $4)
	`
	timeNow := time.Now().Local()

	dbPool := db.DatabaseConnection()
	ctx := context.Background()

	for i := 0; i < 20; i++ {
		name := gimei.NewName()
		_, err := dbPool.Exec(
			ctx, cmd,
			name.Kanji(), gofakeit.Email(), timeNow, timeNow,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
	}
}

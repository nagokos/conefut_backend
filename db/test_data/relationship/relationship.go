package main

import (
	"context"
	"time"

	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/logger"
)

func main() {
	dbPool := db.DatabaseConnection()
	ctx := context.Background()
	timeNow := time.Now().Local()

	cmd := `
	  SELECT id FROM users LIMIT 1
	`
	row := dbPool.QueryRow(ctx, cmd)

	var viewer int
	err := row.Scan(&viewer)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return
	}

	cmd = `
	  SELECT id FROM users WHERE id > $1
	`
	rows, err := dbPool.Query(ctx, cmd, viewer)
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
	defer rows.Close()

	var followerdIDs []int
	for rows.Next() {
		var followedID int
		err := rows.Scan(&followedID)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		followerdIDs = append(followerdIDs, followedID)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}

	cmd = `
	  INSERT INTO relationships (followed_id, follower_id, created_at, updated_at) VALUES ($1, $2, $3, $4)
	`
	for _, id := range followerdIDs {
		_, err := dbPool.Exec(
			ctx, cmd,
			id, viewer, timeNow, timeNow,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}

		_, err = dbPool.Exec(
			ctx, cmd,
			viewer, id, timeNow, timeNow,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
	}
	logger.NewLogger().Info("create relatioship data!")

}

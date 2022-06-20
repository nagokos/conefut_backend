package room

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

func CreateRoom(ctx context.Context, dbPool pgxpool.Pool) (string, error) {
	cmd := `
	  INSERT INTO rooms 
		  (id, created_at, updated_at)
		VALUES
		  ($1, $2, $3)
		RETURNING id
	`

	timeNow := time.Now().Local()
	row := dbPool.QueryRow(ctx, cmd, xid.New().String(), timeNow, timeNow)

	var roomID string
	err := row.Scan(&roomID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return roomID, err
	}

	return roomID, nil
}

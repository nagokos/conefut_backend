package room

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

func CreateRoom(ctx context.Context, tx pgx.Tx) (string, error) {
	cmd := `
	  INSERT INTO rooms 
		  (id, created_at, updated_at)
		VALUES
		  ($1, $2, $3)
		RETURNING id
	`

	timeNow := time.Now().Local()
	row := tx.QueryRow(ctx, cmd, xid.New().String(), timeNow, timeNow)

	var roomID string
	err := row.Scan(&roomID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return roomID, err
	}

	return roomID, nil
}

func GetviewerRooms(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Room, error) {
	currenUser := auth.ForContext(ctx)
	if currenUser == nil {
		logger.NewLogger().Error("user not loggedIn")
		return nil, errors.New("ログインしてください")
	}

	cmd := `
	  SELECT e_1.room_id, u.name, u.avatar
		FROM entries AS e_1
		INNER JOIN (
			SELECT room_id 
			FROM entries 
			WHERE user_id = $1
		) AS e_2
		ON e_1.room_id = e_2.room_id
		INNER JOIN users AS u 
		ON e_1.user_id = u.id
		WHERE NOT e_1.user_id = $2
	`

	rows, err := dbPool.Query(
		ctx, cmd,
		currenUser.ID, currenUser.ID,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	defer rows.Close()

	var rooms []*model.Room
	for rows.Next() {
		var room model.Room
		var entrie model.Entrie
		var user model.User
		err = rows.Scan(&room.ID, &user.Name, &user.Avatar)
		logger.NewLogger().Error(err.Error())
		entrie.User = &user
		room.Entrie = &entrie
		rooms = append(rooms, &room)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return rooms, nil
}

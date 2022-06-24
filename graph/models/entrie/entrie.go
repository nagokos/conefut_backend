package entrie

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func GetEntrieUser(ctx context.Context, dbPool *pgxpool.Pool, roomID string) (*model.User, error) {
	currenUser := auth.ForContext(ctx)
	if currenUser == nil {
		logger.NewLogger().Error("user not loggedIn")
		return nil, errors.New("ログインしてください")
	}

	cmd := `
	  SELECT u.id, u.name, u.avatar
		FROM users AS u 
		INNER JOIN (
			SELECT user_id 
			FROM entries
			WHERE room_id = $1
			AND NOT user_id = $2
		) AS e
		ON u.id = e.user_id
	`

	row := dbPool.QueryRow(
		ctx, cmd,
		roomID, currenUser.ID,
	)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Avatar)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &user, nil
}

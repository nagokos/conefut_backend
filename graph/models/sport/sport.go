package sport

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func GetSports(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Sport, error) {
	cmd := "SELECT id, name FROM sports"
	rows, err := dbPool.Query(ctx, cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var sports []*model.Sport
	for rows.Next() {
		var sport model.Sport
		err := rows.Scan(&sport.DatabaseID, &sport.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		sports = append(sports, &sport)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return sports, nil
}

func GetSport(ctx context.Context, dbPool *pgxpool.Pool, id int) (*model.Sport, error) {
	cmd := "SELECT id, name FROM sports WHERE id = $1"

	var sport model.Sport
	row := dbPool.QueryRow(ctx, cmd, id)
	err := row.Scan(&sport.DatabaseID, &sport.Name)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return &sport, nil
}

func GetSportsByUserID(ctx context.Context, dbPool *pgxpool.Pool, userID int) ([]*model.Sport, error) {
	cmd := `
	  SELECT s.id, s.name
		FROM sports as s
		INNER JOIN user_play_sports as u_s
		ON s.id = u_s.sport_id
		WHERE u_s.user_id = $1
	`
	rows, err := dbPool.Query(ctx, cmd, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sports []*model.Sport
	for rows.Next() {
		var sport model.Sport
		if err := rows.Scan(&sport.DatabaseID, &sport.Name); err != nil {
			return nil, err
		}
		sports = append(sports, &sport)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sports, nil
}

func AddUserPlaySport(ctx context.Context, tx pgx.Tx, userID, sportID int) error {
	cmd := `
	  INSERT INTO user_play_sports
		  (user_id, sport_id, created_at, updated_at)
		VALUES
		  ($1, $2, $3, $4)
	`
	now := time.Now().Local()
	if _, err := tx.Exec(ctx, cmd, userID, sportID, now, now); err != nil {
		return err
	}
	return nil
}

func RemoveUserPlaySport(ctx context.Context, tx pgx.Tx, userID, sportID int) error {
	cmd := `
	  DELETE FROM user_play_sports as u_s
		WHERE u_s.user_id = $1
		AND u_s.sport_id = $2
	`
	if _, err := tx.Exec(ctx, cmd, userID, sportID); err != nil {
		return err
	}
	return nil
}

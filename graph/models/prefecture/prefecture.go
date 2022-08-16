package prefecture

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

type NullablePrefecture struct {
	ID         *string
	DatabaseID *int
	Name       *string
}

func GetPrefectures(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Prefecture, error) {
	var prefectures []*model.Prefecture

	cmd := "SELECT id, name FROM prefectures"
	rows, err := dbPool.Query(ctx, cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prefecture model.Prefecture
		err := rows.Scan(&prefecture.DatabaseID, &prefecture.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		prefectures = append(prefectures, &prefecture)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return prefectures, nil
}

func GetPrefecture(ctx context.Context, dbPool *pgxpool.Pool, id int) (*model.Prefecture, error) {
	cmd := "SELECT id, name FROM prefectures WHERE id = $1"

	var prefecture model.Prefecture
	row := dbPool.QueryRow(ctx, cmd, id)
	err := row.Scan(&prefecture.DatabaseID, &prefecture.Name)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &prefecture, nil
}

func GetPrefecturesByUserID(ctx context.Context, dbPool *pgxpool.Pool, userID int) ([]*model.Prefecture, error) {
	cmd := `
	  SELECT p.id, p.name
		FROM prefectures as p
		INNER JOIN user_activity_areas as u_a
		ON p.id = u_a.prefecture_id
		WHERE u_a.user_id = $1
	`
	rows, err := dbPool.Query(ctx, cmd, userID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var prefectures []*model.Prefecture
	for rows.Next() {
		var prefecture model.Prefecture
		if err := rows.Scan(&prefecture.DatabaseID, &prefecture.Name); err != nil {
			logger.NewLogger().Error(err.Error())
		}
		prefectures = append(prefectures, &prefecture)
	}
	if err := rows.Err(); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return prefectures, nil
}

func AddUserActivityArea(ctx context.Context, tx pgx.Tx, userID, prefectureID int) error {
	cmd := `
	  INSERT INTO user_activity_areas
		  (user_id, prefecture_id, created_at, updated_at)
		VALUES
		  ($1, $2, $3, $4)
	`
	now := time.Now().Local()
	if _, err := tx.Exec(ctx, cmd, userID, prefectureID, now, now); err != nil {
		return err
	}
	return nil
}

func RemoveUserActivieArea(ctx context.Context, tx pgx.Tx, userID, prefectureID int) error {
	cmd := `
	  DELETE FROM user_activity_areas as u_a
		WHERE u_a.user_id = $1
		AND u_a.prefecture_id = $2
	`
	if _, err := tx.Exec(ctx, cmd, userID, prefectureID); err != nil {
		return err
	}
	return nil
}

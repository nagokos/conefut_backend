package stock

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

func CreateStock(ctx context.Context, dbPool *pgxpool.Pool, recId string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	cmd := "SELECT COUNT(DISTINCT id) FROM stocks WHERE user_id = $1 AND recruitment_id = $2"
	row := dbPool.QueryRow(ctx, cmd, currentUser.ID, recId)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	if count == 1 {
		logger.NewLogger().Error("Already stocked.")
		return false, errors.New("既にストックしています")
	}

	timeNow := time.Now().Local()

	cmd = "INSERT INTO stocks (id, recruitment_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err = dbPool.Exec(ctx, cmd, xid.New().String(), recId, currentUser.ID, timeNow, timeNow)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	return true, nil
}

func DeleteStock(ctx context.Context, dbPool *pgxpool.Pool, recId string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	cmd := "DELETE FROM stocks WHERE user_id = $1 AND recruitment_id = $2"
	_, err := dbPool.Exec(ctx, cmd, currentUser.ID, recId)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, errors.New("delete stock error")
	}

	return true, nil
}

func GetStockedCount(ctx context.Context, dbPool *pgxpool.Pool, recId string) (int, error) {
	cmd := "SELECT COUNT(DISTINCT id) FROM stocks WHERE recruitment_id = $1"
	row := dbPool.QueryRow(ctx, cmd, recId)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return count, err
	}

	return count, nil
}

func CheckStocked(ctx context.Context, dbPool *pgxpool.Pool, recId string) (bool, error) {
	currentUser := auth.ForContext(ctx)

	if currentUser == nil {
		return false, nil
	}

	cmd := "SELECT COUNT(DISTINCT id) FROM stocks WHERE user_id = $1 AND recruitment_id = $2"
	row := dbPool.QueryRow(ctx, cmd, currentUser.ID, recId)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	var isStocked bool
	if count == 1 {
		isStocked = true
	}

	return isStocked, nil
}

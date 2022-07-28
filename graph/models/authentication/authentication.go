package authentication

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/logger"
)

type Claims struct {
	ID       string `json:"sub"`
	Name     string `json:"name"`
	Avatar   string `json:"picture"`
	Email    string `json:"email"`
	Provider string
}

func (c *Claims) CreateFrom(ctx context.Context, dbPool *pgxpool.Pool) (int, error) {
	timeNow := time.Now().Local()

	cmd := `
		INSERT INTO users
		(name, email, avatar, email_verification_status, last_sign_in_at, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`
	row := dbPool.QueryRow(
		ctx, cmd,
		c.Name, c.Email, c.Avatar, "verified", timeNow, timeNow, timeNow,
	)

	var userID int
	if err := row.Scan(&userID); err != nil {
		logger.NewLogger().Error(err.Error())
		return 0, err
	}

	cmd = `
	  INSERT INTO authentications (provider, uid, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := dbPool.Exec(
		ctx, cmd,
		c.Provider, c.ID, userID, timeNow, timeNow,
	); err != nil {
		logger.NewLogger().Error(err.Error())
		return 0, nil
	}

	return userID, nil
}

func (c *Claims) CheckAuthAlreadyExists(ctx context.Context, dbPool *pgxpool.Pool) (bool, error) {
	cmd := `
	  SELECT COUNT(DISTINCT id)
		FROM authentications
		WHERE provider = $1
		AND uid = $2
	`
	row := dbPool.QueryRow(
		ctx, cmd,
		c.Provider, c.ID,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

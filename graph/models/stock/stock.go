package stock

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func AddStock(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID int) (*model.FeedbackStock, error) {
	feedback := model.FeedbackStock{ID: utils.GenerateUniqueID("Stock", recruitmentID)}
	viewer := user.GetViewer(ctx)
	cmd := `
	  SELECT COUNT(DISTINCT id) 
	  FROM stocks 
		WHERE user_id = $1 
		AND recruitment_id = $2
	`
	row := dbPool.QueryRow(ctx, cmd, viewer.DatabaseID, recruitmentID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	if count > 0 {
		logger.NewLogger().Error("Already stocked.")
		feedback.ViewerDoesStock = true
		return &feedback, nil
	}

	timeNow := time.Now().Local()
	cmd = `
	  INSERT INTO stocks 
		  (recruitment_id, user_id, created_at, updated_at) 
		VALUES 
		  ($1, $2, $3, $4)
	`
	if _, err = dbPool.Exec(ctx, cmd, recruitmentID, viewer.DatabaseID, timeNow, timeNow); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	cmd = `
	  SELECT id, title, closing_at, user_id
		FROM recruitments
		WHERE id = $1
	`
	row = dbPool.QueryRow(ctx, cmd, recruitmentID)
	var recruitment model.Recruitment
	if err := row.Scan(&recruitment.DatabaseID, &recruitment.Title, &recruitment.ClosingAt, &recruitment.UserID); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	feedback.FeedbackRecruitmentEdge = &model.RecruitmentEdge{
		Cursor: utils.GenerateUniqueID("Recruitment", recruitment.DatabaseID),
		Node:   &recruitment,
	}
	feedback.ViewerDoesStock = true
	return &feedback, nil
}

func RemoveStock(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID int) (*model.FeedbackStock, error) {
	viewer := user.GetViewer(ctx)
	cmd := `
	  DELETE FROM stocks 
		WHERE user_id = $1 
		AND recruitment_id = $2
	`
	if _, err := dbPool.Exec(ctx, cmd, viewer.DatabaseID, recruitmentID); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	removeID := utils.GenerateUniqueID("Recruitment", recruitmentID)
	feedback := model.FeedbackStock{
		ID:                   utils.GenerateUniqueID("Stock", recruitmentID),
		ViewerDoesStock:      false,
		RemovedRecruitmentID: &removeID,
	}
	return &feedback, nil
}

func CheckStocked(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID int) (*model.FeedbackStock, error) {
	feedback := model.FeedbackStock{ID: utils.GenerateUniqueID("Stock", recruitmentID)}
	viewer := user.GetViewer(ctx)
	if viewer == nil {
		return &feedback, nil
	}
	cmd := `
	  SELECT COUNT(DISTINCT id) 
		FROM stocks 
		WHERE user_id = $1 
		AND recruitment_id = $2
	`
	row := dbPool.QueryRow(ctx, cmd, viewer.DatabaseID, recruitmentID)
	var count int
	if err := row.Scan(&count); err != nil {
		logger.NewLogger().Error(err.Error())
		return &feedback, err
	}
	if count > 0 {
		feedback.ViewerDoesStock = true
	}
	return &feedback, nil
}

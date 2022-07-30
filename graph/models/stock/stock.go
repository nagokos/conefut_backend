package stock

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func AddStock(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID string) (*model.FeedbackStock, error) {
	feedback := model.FeedbackStock{
		ID: utils.GenerateUniqueID("Stock", utils.DecodeUniqueID(recruitmentID)),
	}

	viewer := user.GetViewer(ctx)

	cmd := `
	  SELECT COUNT(DISTINCT id) 
	  FROM stocks 
		WHERE user_id = $1 
		AND recruitment_id = $2
	`
	row := dbPool.QueryRow(ctx, cmd, viewer.DatabaseID, utils.DecodeUniqueID(recruitmentID))

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return &feedback, err
	}

	if count == 1 {
		logger.NewLogger().Error("Already stocked.")
		feedback.ViewerDoesStock = true
		return &feedback, nil
	}

	timeNow := time.Now().Local()

	cmd = "INSERT INTO stocks (recruitment_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4)"
	_, err = dbPool.Exec(
		ctx, cmd,
		utils.DecodeUniqueID(recruitmentID), viewer.DatabaseID, timeNow, timeNow,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return &feedback, err
	}

	cmd = `
	  SELECT id, title, closing_at, user_id
		FROM recruitments
		WHERE id = $1
	`
	row = dbPool.QueryRow(ctx, cmd, utils.DecodeUniqueID(recruitmentID))

	var recruitment model.Recruitment
	err = row.Scan(&recruitment.DatabaseID, &recruitment.Title, &recruitment.ClosingAt, &recruitment.UserID)
	if err != nil {
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

func RemoveStock(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID string) (*model.FeedbackStock, error) {
	feedback := model.FeedbackStock{
		ID: utils.GenerateUniqueID("Stock", utils.DecodeUniqueID(recruitmentID)),
	}

	viewer := user.GetViewer(ctx)

	cmd := `
	  DELETE FROM stocks 
		WHERE user_id = $1 
		AND recruitment_id = $2
		RETURNING id
	`
	_, err := dbPool.Exec(ctx, cmd, viewer.DatabaseID, utils.DecodeUniqueID(recruitmentID))
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return &feedback, errors.New("delete stock error")
	}
	feedback.RemovedRecruitmentID = &recruitmentID

	return &feedback, nil
}

func CheckStocked(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID string) (*model.FeedbackStock, error) {
	feedback := model.FeedbackStock{
		ID: utils.GenerateUniqueID("Stock", utils.DecodeUniqueID(recruitmentID)),
	}

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
	row := dbPool.QueryRow(ctx, cmd, viewer.DatabaseID, utils.DecodeUniqueID(recruitmentID))

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return &feedback, err
	}

	if count == 1 {
		feedback.ViewerDoesStock = true
	}

	return &feedback, nil
}

// func GetStockedCount(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID string) (int, error) {
// 	cmd := "SELECT COUNT(DISTINCT id) FROM stocks WHERE recruitment_id = $1"
// 	row := dbPool.QueryRow(ctx, cmd, recruitmentID)

// 	var count int
// 	err := row.Scan(&count)
// 	if err != nil {
// 		logger.NewLogger().Error(err.Error())
// 		return count, err
// 	}

// 	return count, nil
// }

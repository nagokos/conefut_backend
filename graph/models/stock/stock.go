package stock

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/recruitment"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/logger"
)

func AddStock(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID int) (*model.AddStockResult, error) {
	isStocked, err := IsViewerStocked(ctx, dbPool, recruitmentID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	if isStocked {
		logger.NewLogger().Error("Already stocked")
		return &model.AddStockResult{}, nil
	}

	timeNow := time.Now().Local()
	cmd := `
	  INSERT INTO stocks 
		  (recruitment_id, user_id, created_at, updated_at) 
		VALUES 
		  ($1, $2, $3, $4)
	`
	viewer := user.GetViewer(ctx)
	if _, err := dbPool.Exec(ctx, cmd, recruitmentID, viewer.DatabaseID, timeNow, timeNow); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	recruitment, err := recruitment.GetRecruitment(ctx, dbPool, recruitmentID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	result := model.AddStockResult{
		FeedbackStock: &model.FeedbackStock{
			IsViewerStocked: true,
			RecruitmentID:   recruitmentID,
		},
		RecruitmentEdge: &model.RecruitmentEdge{
			Node: recruitment,
		},
	}
	return &result, nil
}

func RemoveStock(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID int) (*model.RemoveStockResult, error) {
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

	result := model.RemoveStockResult{
		FeedbackStock: &model.FeedbackStock{
			RecruitmentID: recruitmentID,
		},
	}
	return &result, nil
}

func GetFeedbackStock(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID int) (*model.FeedbackStock, error) {
	isStocked, err := IsViewerStocked(ctx, dbPool, recruitmentID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	feedback := model.FeedbackStock{
		RecruitmentID:   recruitmentID,
		IsViewerStocked: isStocked,
	}
	return &feedback, nil
}

func IsViewerStocked(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID int) (bool, error) {
	viewer := user.GetViewer(ctx)
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
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

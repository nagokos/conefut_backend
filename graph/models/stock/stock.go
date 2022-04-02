package stock

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/ent/stock"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

func CreateStock(ctx context.Context, dbConnection *sql.DB, recId string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	timeNow := time.Now().Local()

	cmd := fmt.Sprintf(`INSERT INTO %s (id, recruitment_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`, db.StockTable)
	_, err := dbConnection.Exec(cmd, xid.New().String(), recId, currentUser.ID, timeNow, timeNow)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	return true, nil
}

func DeleteStock(ctx context.Context, dbConnection sql.DB, recId string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	cmd := fmt.Sprintf(`DELETE FROM %s WHERE recruitment_id = $1 AND user_id = $2`, db.StockTable)
	_, err := dbConnection.Exec(cmd, recId, currentUser.ID)
	logger.NewLogger().Info(err.Error())
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, errors.New("delete stock error")
	}

	return true, nil
}

func GetStockedCount(ctx context.Context, client ent.Client, recId string) (int, error) {
	res, err := client.Recruitment.
		Query().
		Where(recruitment.ID(recId)).
		QueryStocks().
		Count(ctx)
	if err != nil {
		logger.NewLogger().Sugar().Errorf("get stocked count error %s", err.Error())
		return 0, err
	}

	return res, nil
}

func CheckStocked(ctx context.Context, client ent.Client, recId string) (bool, error) {
	currentUser := auth.ForContext(ctx)

	if currentUser == nil {
		return false, nil
	}

	res, err := client.Stock.
		Query().
		Where(
			stock.UserID(currentUser.ID),
			stock.RecruitmentID(recId),
		).
		Exist(ctx)
	if err != nil {
		logger.NewLogger().Sugar().Errorf("check stocked error %s", err.Error())
		return false, err
	}

	return res, nil
}

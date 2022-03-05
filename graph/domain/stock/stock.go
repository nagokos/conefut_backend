package stock

import (
	"context"
	"errors"
	"fmt"

	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/ent/stock"
	"github.com/nagokos/connefut_backend/ent/user"
	"github.com/nagokos/connefut_backend/logger"
)

func CreateStock(ctx context.Context, client ent.Client, recId string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	_, err := client.Stock.
		Create().
		SetRecruitmentID(recId).
		SetUserID(currentUser.ID).
		Save(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("create stock error %s", err.Error()))
		return false, err
	}

	return true, nil
}

func DeleteStock(ctx context.Context, client ent.Client, recId string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	res, err := client.Stock.
		Delete().
		Where(
			stock.HasUserWith(
				user.ID(currentUser.ID),
			),
			stock.RecruitmentID(recId),
		).
		Exec(ctx)
	if res == 0 {
		logger.Log.Error().Msg("delete stock error: res zero")
		return false, errors.New("ストックを削除できませんでした")
	}

	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("delete stock error %s", err.Error()))
		return false, err
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
		logger.Log.Error().Msg(fmt.Sprintf("get stocked count error %s", err.Error()))
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
			stock.HasUserWith(
				user.ID(currentUser.ID),
			),
			stock.HasRecruitmentWith(
				recruitment.ID(recId),
			),
		).
		Exist(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("check stocked error %s", err.Error()))
		return false, err
	}

	return res, nil
}

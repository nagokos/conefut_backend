package applicant

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func CheckApplied(ctx context.Context, client ent.Client, recId string) (bool, error) {
	currentUser := auth.ForContext(ctx)

	isApplied, err := client.Applicant.
		Query().
		Where(
			applicant.UserID(currentUser.ID),
			applicant.RecruitmentID(recId),
		).
		Exist(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	return isApplied, nil
}

func GetAppliedCounts(ctx context.Context, client ent.Client, recId string) (int, error) {
	res, err := client.Applicant.
		Query().
		Where(
			applicant.RecruitmentID(recId),
		).
		Count(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return res, err
	}

	return res, nil
}

func CreateApplicant(ctx context.Context, client ent.Client, recId string, status model.ManagementStatus) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	isApplied, err := client.Applicant.
		Query().
		Where(
			applicant.UserID(currentUser.ID),
			applicant.RecruitmentID(recId),
		).
		Exist(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}

	if isApplied {
		logger.NewLogger().Error("already applied for this position")
		return false, errors.New("この募集には既に応募しています")
	}

	_, err = client.Applicant.
		Create().
		SetRecruitmentID(recId).
		SetManagementStatus(applicant.ManagementStatus(strings.ToLower(string(status)))).
		SetUserID(currentUser.ID).
		Save(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	return true, nil
}

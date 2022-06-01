package search

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

type SearchParams struct {
	UseAfter  bool
	After     string
	UseBefore bool
	Before    string
	NumRows   int
	Options   SearchRecruitmentParams
}

type SearchRecruitmentParams struct {
	UseCompetition bool
	CompetitionID  string
	UsePrefecture  bool
	PrefectureID   string
	UseType        bool
	Type           model.Type
	UseStartAt     bool
	StartAt        time.Time
}

func NewSearchParams(after *string, before *string, first *int, last *int, options *model.SearchRecruitmentInput) (SearchParams, error) {
	var sp = SearchParams{}

	sp.UseAfter = (after != nil)
	sp.UseBefore = (before != nil)
	useFirst := (first != nil)
	useLast := (last != nil)

	if useFirst && !sp.UseAfter && !useLast && !sp.UseBefore {
		sp.NumRows = *first
	} else if useFirst && sp.UseAfter && !useLast && !sp.UseBefore {
		sp.NumRows = *first
		sp.After = *after
	} else if useLast && sp.UseBefore && !useFirst && !sp.UseAfter {
		sp.NumRows = *last
		sp.Before = *before
	} else {
		logger.NewLogger().Error("search params validation error")
		return SearchParams{}, errors.New("{first}, {after, first}, {before, last}のいずれかの組み合わせで指定してください")
	}

	var srp = SearchRecruitmentParams{}

	srp.UseCompetition = (options.CompetitionID != nil)
	srp.UsePrefecture = (options.PrefectureID != nil)
	srp.UseStartAt = (options.StartAt != nil)
	srp.UseType = (options.Type != nil)

	if srp.UseCompetition {
		srp.CompetitionID = *options.CompetitionID
	}

	if srp.UsePrefecture {
		srp.PrefectureID = *options.PrefectureID
	}

	if srp.UseStartAt {
		srp.StartAt = *options.StartAt
	}

	if srp.UseType {
		srp.Type = model.Type(strings.ToLower(string(*options.Type)))
	}

	sp.Options = srp

	return sp, nil
}

func NextPageExists(ctx context.Context, dbPool *pgxpool.Pool, nextID string, params SearchParams, sort string) (bool, error) {
	cmd := fmt.Sprintf(`
		SELECT COUNT(DISTINCT r.id)
		FROM 
			(
				SELECT id FROM recruitments
				WHERE status = $1
				AND ($2 OR competition_id = $3) 
				AND id < $4
				ORDER BY id %s
			) AS r
		LIMIT 1
	`, sort)

	row := dbPool.QueryRow(
		ctx, cmd,
		"published", !params.Options.UseCompetition, params.Options.CompetitionID, nextID,
	)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	var isNextPage bool
	if count != 0 {
		isNextPage = true
	}

	return isNextPage, nil
}

func PreviousPageExists(ctx context.Context, dbPool *pgxpool.Pool, previousID string, params SearchParams, sort string) (bool, error) {
	cmd := fmt.Sprintf(`
		SELECT COUNT(DISTINCT r.id)
		FROM 
			(
				SELECT id FROM recruitments
				WHERE status = $1
				AND ($2 OR competition_id = $3) 
				AND id > $4
				ORDER BY id %s
			) AS r
		LIMIT 1
	`, sort)

	row := dbPool.QueryRow(
		ctx, cmd,
		"published", !params.Options.UseCompetition, params.Options.CompetitionID, previousID,
	)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	var isPreviousPage bool
	if count != 0 {
		isPreviousPage = true
	}

	return isPreviousPage, nil
}

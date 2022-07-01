package search

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/logger"
)

type SearchParams struct {
	UseAfter  bool
	After     string
	UseBefore bool
	Before    string
	NumRows   int
}

func NewSearchParams(after *string, before *string, first *int, last *int) (SearchParams, error) {
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

	return sp, nil
}

func NextPageExists(ctx context.Context, dbPool *pgxpool.Pool, nextID string, params SearchParams, sort string) (bool, error) {
	cmd := fmt.Sprintf(`
		SELECT COUNT(DISTINCT r.id)
		FROM 
			(
				SELECT id FROM recruitments
				WHERE status = $1
				AND id < $2
				ORDER BY id %s
			) AS r
		LIMIT 1
	`, sort)

	row := dbPool.QueryRow(
		ctx, cmd,
		"published", nextID,
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
				AND id > $2
				ORDER BY id %s
			) AS r
		LIMIT 1
	`, sort)

	row := dbPool.QueryRow(
		ctx, cmd,
		"published", previousID,
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

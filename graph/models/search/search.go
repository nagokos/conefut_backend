package search

import (
	"errors"

	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

type SearchParams struct {
	UseAfter  bool
	After     int
	UseBefore bool
	Before    int
	NumRows   int
}

func NewSearchParams(first *int, after *string, last *int, before *string) (SearchParams, error) {
	var sp = SearchParams{}

	sp.UseAfter = (after != nil)
	sp.UseBefore = (before != nil)
	useFirst := (first != nil)
	useLast := (last != nil)

	if useFirst && !sp.UseAfter && !useLast && !sp.UseBefore {
		sp.NumRows = *first
		sp.After = 0
	} else if useFirst && sp.UseAfter && !useLast && !sp.UseBefore {
		sp.NumRows = *first
		sp.After = utils.DecodeUniqueIDIdentifierOnly(*after)
	} else if useLast && sp.UseBefore && !useFirst && !sp.UseAfter {
		sp.NumRows = *last
		sp.Before = utils.DecodeUniqueIDIdentifierOnly(*before)
	} else {
		logger.NewLogger().Error("search params validation error")
		return SearchParams{}, errors.New("{first}, {after, first}, {before, last}のいずれかの組み合わせで指定してください")
	}

	return sp, nil
}

package competition

import (
	"context"
	"fmt"

	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/competition"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func GetCompetitions(ctx context.Context, client *ent.CompetitionClient) ([]*model.Competition, error) {
	res, err := client.
		Query().
		Order(ent.Asc(competition.FieldID)).
		All(ctx)

	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln("get competitions error:", err.Error()))
		return nil, err
	}

	var competitions []*model.Competition
	for _, v := range res {
		competitions = append(competitions, &model.Competition{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return competitions, nil
}

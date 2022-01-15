package prefecture

import (
	"context"
	"fmt"

	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/prefecture"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func GetPrefectures(client ent.PrefectureClient, ctx context.Context) ([]*model.Prefecture, error) {
	var prefectures []*model.Prefecture

	res, err := client.
		Query().
		Order(ent.Asc(prefecture.FieldID)).
		All(ctx)

	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln("get prefectures: ", err))
		return prefectures, err
	}

	for _, prefecture := range res {
		prefectures = append(prefectures, &model.Prefecture{
			ID:   prefecture.ID,
			Name: prefecture.Name,
		})
	}

	return prefectures, nil
}

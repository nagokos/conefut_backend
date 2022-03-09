//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/logger"
)

func main() {
	client := db.DatabaseConnection()
	defer client.Close()

	ctx := context.Background()

	comp, err := client.Competition.Query().First(ctx)
	if err != nil {
		logger.Log.Error().Msg(err.Error())
		return
	}

	pref, err := client.Prefecture.Query().First(ctx)
	if err != nil {
		logger.Log.Error().Msg(err.Error())
		return
	}

	user, err := client.User.Query().First(ctx)
	if err != nil {
		logger.Log.Error().Msg(err.Error())
		return
	}

	var recruitments []*ent.Recruitment
	for i := 0; i < 20; i++ {
		recruitment := &ent.Recruitment{
			Type:      "opponent",
			Title:     fmt.Sprintf("%v 明日の午後からサッカーできる人を探しています。ご連絡お待ちしています。", i),
			Content:   fmt.Sprintf("%v 明日の午後からサッカーできる人を探しています。場所は後ほど連絡します。", i),
			Place:     "埼玉スタジアム2002",
			StartAt:   time.Now().Add(time.Hour * 20),
			ClosingAt: time.Now().Add(time.Hour * 19),
			Capacity:  1,
			Edges: ent.RecruitmentEdges{
				Competition: comp,
				Prefecture:  pref,
				User:        user,
			},
		}
		recruitments = append(recruitments, recruitment)
	}

	bulk := make([]*ent.RecruitmentCreate, len(recruitments))

	for i, rec := range recruitments {
		bulk[i] = client.Recruitment.
			Create().
			SetTitle(rec.Title).
			SetType(rec.Type).
			SetNillableStartAt(&rec.StartAt).
			SetNillableContent(&rec.Content).
			SetNillablePlace(&rec.Place).
			SetStatus(recruitment.StatusPublished).
			SetCapacity(rec.Capacity).
			SetNillableClosingAt(&rec.ClosingAt).
			SetNillableCompetitionID(&rec.Edges.Competition.ID).
			SetNillablePrefectureID(&rec.Edges.Prefecture.ID).
			SetUserID(rec.Edges.User.ID)
	}

	_, err = client.Recruitment.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		logger.Log.Error().Msg(err.Error())
	}
}

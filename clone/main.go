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

func updateRecruitmentStatus(ctx context.Context, client ent.Client) {
	_, err := client.Recruitment.
		Update().
		Where(
			recruitment.StatusEQ(recruitment.StatusPublished),
			recruitment.ClosingAtLT(time.Now().Local()),
		).
		SetStatus(recruitment.StatusClosed).
		Save(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("tick update recruitment status error %s", err.Error()))
	}
}

func main() {
	client := db.DatabaseConnection()
	ctx := context.Background()
	t := time.NewTicker(time.Second * 60)
	defer t.Stop()

	for {
		<-t.C
		updateRecruitmentStatus(ctx, *client)
	}
}

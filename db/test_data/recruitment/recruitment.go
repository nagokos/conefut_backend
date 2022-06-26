//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/logger"
)

type recruitment struct {
	id            int
	recType       string
	title         string
	detail        string
	place         string
	startAt       time.Time
	closingAt     time.Time
	competitionID int
	prefectureID  int
	userID        int
	status        string
	published_at  time.Time
	createdAt     time.Time
	updatedAt     time.Time
}

func main() {
	dbPool := db.DatabaseConnection()
	defer dbPool.Close()

	ctx := context.Background()

	cmd := "SELECT id FROM competitions LIMIT 1"
	row := dbPool.QueryRow(ctx, cmd)

	var competitionID int
	err := row.Scan(&competitionID)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	cmd = "SELECT id FROM prefectures LIMIT 1"
	row = dbPool.QueryRow(ctx, cmd)

	var prefectureID int
	err = row.Scan(&prefectureID)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	cmd = "SELECT id FROM users LIMIT 1"
	row = dbPool.QueryRow(ctx, cmd)

	var userID int
	err = row.Scan(&userID)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	var recruitments []*recruitment
	for i := 0; i < 20; i++ {
		detail := `初めまして。
東京都社会人3部リーグに所属しているFortuna TOKYOと申します。
下記の通りグラウンドが取得できましたので、対戦相手の募集をいたします。
※先着順ではございません。
※他でも打診をしております。
応募を多数いただく場合はチーム内協議の上決定いたします。

日時:4月16日（土）8:30〜10:30
場所:朝霞中央公園陸上競技場(綺麗な人工芝)
費用:6000円

〈募集条件〉
①暴力、暴言、ラフプレーなどが無いよう、リスペクトの精神を持ってプレーできる事
②対戦決定後キャンセルしない事
③当日審判、グラウンドの準備、整備にご協力頂ける事
④13人以上揃う事
⑤競技思考である事
⑥コロナ感染対策にご協力いただける事

◆当チームプロフィール◆
チーム名  Fortuna TOKYO
ユニフォーム色 青 or 赤
平均年齢  27

対戦をご希望される方は、
チーム名：
代表者名：
代表者電話番号：
ユニフォーム色：
所属リーグ等チーム情報：

上記ご記入の上ご連絡ください。

以上、よろしくお願いいたします。`

		recruitment := &recruitment{
			recType:       "opponent",
			title:         fmt.Sprintf("対戦相手募集 朝霞中央公園陸上競技場(人工芝) %v", i+1),
			place:         "朝霞中央公園陸上競技場",
			startAt:       time.Now().Add(time.Hour * 240).Local(),
			closingAt:     time.Now().Add(time.Hour * 230).Local(),
			competitionID: competitionID,
			prefectureID:  prefectureID,
			userID:        userID,
			status:        "published",
			createdAt:     time.Now().Local(),
			updatedAt:     time.Now().Local(),
			detail:        detail,
			published_at:  time.Now().Local(),
		}
		recruitments = append(recruitments, recruitment)
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer tx.Rollback(ctx)

	cmd = `
	  INSERT INTO recruitments 
		  (title, type, place, start_at, detail, closing_at, competition_id, prefecture_id, user_id, created_at, updated_at, status, published_at)
		VALUES 
		  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`

	for _, recruitment := range recruitments {
		if _, err := tx.Exec(
			ctx, cmd,
			recruitment.title, recruitment.recType, recruitment.place, recruitment.startAt, recruitment.detail,
			recruitment.closingAt, recruitment.competitionID, recruitment.prefectureID, recruitment.userID, recruitment.createdAt, recruitment.updatedAt, recruitment.status, recruitment.published_at,
		); err != nil {
			logger.NewLogger().Fatal(err.Error())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	logger.NewLogger().Info("create recruitments data!")
}

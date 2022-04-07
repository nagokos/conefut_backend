package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

func InsertPrefectures(ctx context.Context, dbPool *pgxpool.Pool) {

	prefectures := []struct {
		name string
	}{
		{"北海道"},
		{"青森県"},
		{"岩手県"},
		{"宮城県"},
		{"秋田県"},
		{"山形県"},
		{"福島県"},
		{"茨城県"},
		{"栃木県"},
		{"群馬県"},
		{"埼玉県"},
		{"千葉県"},
		{"東京都"},
		{"神奈川県"},
		{"新潟県"},
		{"富山県"},
		{"石川県"},
		{"福井県"},
		{"山梨県"},
		{"長野県"},
		{"岐阜県"},
		{"静岡県"},
		{"愛知県"},
		{"三重県"},
		{"滋賀県"},
		{"京都府"},
		{"大阪府"},
		{"兵庫県"},
		{"奈良県"},
		{"和歌山県"},
		{"鳥取県"},
		{"島根県"},
		{"岡山県"},
		{"広島県"},
		{"山口県"},
		{"徳島県"},
		{"香川県"},
		{"愛媛県"},
		{"高知県"},
		{"福岡県"},
		{"佐賀県"},
		{"長崎県"},
		{"熊本県"},
		{"大分県"},
		{"宮崎県"},
		{"鹿児島県"},
		{"沖縄県"},
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer tx.Rollback(ctx)

	cmd := "INSERT INTO prefectures (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)"

	for _, prefecture := range prefectures {
		ID := xid.New().String()
		timeNow := time.Now().Local()
		if _, err := tx.Exec(ctx, cmd, ID, prefecture.name, timeNow, timeNow); err != nil {
			logger.NewLogger().Fatal(err.Error())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	logger.NewLogger().Info("create prefectures data!")
}

func InsertCompetitions(ctx context.Context, dbPool *pgxpool.Pool) {
	competitions := []struct {
		name string
	}{
		{"サッカー"},
		{"フットサル"},
		{"ソサイチ"},
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer tx.Rollback(ctx)

	cmd := "INSERT INTO competitions (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)"
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	for _, competition := range competitions {
		ID := xid.New().String()
		timeNow := time.Now().Local()
		if _, err := tx.Exec(ctx, cmd, ID, competition.name, timeNow, timeNow); err != nil {
			logger.NewLogger().Fatal(err.Error())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	logger.NewLogger().Info("create competitions data!")
}

func InsertTags(ctx context.Context, dbPool *pgxpool.Pool) {
	tags := []struct {
		name string
	}{
		{"エンジョイ"},
		{"男女mix"},
		{"シニア"},
		{"ガチ"},
		{"誰でもok"},
		{"経験者"},
		{"初心者歓迎"},
		{"競技志向"},
		{"急募"},
		{"人工芝"},
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer tx.Rollback(ctx)

	cmd := "INSERT INTO tags (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)"

	for _, tag := range tags {
		ID := xid.New().String()
		timeNow := time.Now().Local()
		if _, err := tx.Exec(ctx, cmd, ID, tag.name, timeNow, timeNow); err != nil {
			logger.NewLogger().Fatal(err.Error())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	logger.NewLogger().Info("create tags data!")
}

func main() {
	dbPool := db.DatabaseConnection()
	defer dbPool.Close()

	ctx := context.Background()

	InsertPrefectures(ctx, dbPool)
	InsertCompetitions(ctx, dbPool)
	InsertTags(ctx, dbPool)
}

package main

import (
	"database/sql"
	"time"

	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

func InsertPrefectures(dbConnection *sql.DB) {

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

	tx, err := dbConnection.Begin()
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO prefectures (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer stmt.Close()

	for _, prefecture := range prefectures {
		ID := xid.New().String()
		timeNow := time.Now().Local()
		if _, err := stmt.Exec(ID, prefecture.name, timeNow, timeNow); err != nil {
			logger.NewLogger().Fatal(err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	logger.NewLogger().Info("create prefectures data!")
}

func InsertCompetitions(dbConnection *sql.DB) {
	competitions := []struct {
		name string
	}{
		{"サッカー"},
		{"フットサル"},
		{"ソサイチ"},
	}

	tx, err := dbConnection.Begin()
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO competitions (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer stmt.Close()

	for _, competition := range competitions {
		ID := xid.New().String()
		timeNow := time.Now().Local()
		if _, err := stmt.Exec(ID, competition.name, timeNow, timeNow); err != nil {
			logger.NewLogger().Fatal(err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	logger.NewLogger().Info("create competitions data!")
}

func InsertTags(dbConnection *sql.DB) {
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

	tx, err := dbConnection.Begin()
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO tags (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}
	defer stmt.Close()

	for _, tag := range tags {
		ID := xid.New().String()
		timeNow := time.Now().Local()
		if _, err := stmt.Exec(ID, tag.name, timeNow, timeNow); err != nil {
			logger.NewLogger().Fatal(err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	logger.NewLogger().Info("create tags data!")
}

func main() {
	dbConnection := db.DatabaseConnection()
	defer dbConnection.Close()

	InsertPrefectures(dbConnection)
	InsertCompetitions(dbConnection)
	InsertTags(dbConnection)
}

package main

import (
	"context"

	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/logger"
)

func InsertPrefectures(ctx context.Context, client *ent.Client) error {
	prefectures := []string{
		"北海道",
		"青森県",
		"岩手県",
		"宮城県",
		"秋田県",
		"山形県",
		"福島県",
		"茨城県",
		"栃木県",
		"群馬県",
		"埼玉県",
		"千葉県",
		"東京都",
		"神奈川県",
		"新潟県",
		"富山県",
		"石川県",
		"福井県",
		"山梨県",
		"長野県",
		"岐阜県",
		"静岡県",
		"愛知県",
		"三重県",
		"滋賀県",
		"京都府",
		"大阪府",
		"兵庫県",
		"奈良県",
		"和歌山県",
		"鳥取県",
		"島根県",
		"岡山県",
		"広島県",
		"山口県",
		"徳島県",
		"香川県",
		"愛媛県",
		"高知県",
		"福岡県",
		"佐賀県",
		"長崎県",
		"熊本県",
		"大分県",
		"宮崎県",
		"鹿児島県",
		"沖縄県",
	}
	bulk := make([]*ent.PrefectureCreate, len(prefectures))
	for i, name := range prefectures {
		bulk[i] = client.Prefecture.Create().SetName(name)
	}
	_, err := client.Prefecture.CreateBulk(bulk...).Save(ctx)
	return err
}

func InsertCompetitions(ctx context.Context, client *ent.Client) error {
	competitions := []string{
		"サッカー",
		"フットサル",
		"ソサイチ",
	}
	bulk := make([]*ent.CompetitionCreate, len(competitions))
	for i, name := range competitions {
		bulk[i] = client.Competition.Create().SetName(name)
	}
	_, err := client.Competition.CreateBulk(bulk...).Save(ctx)
	return err
}

func InsertTags(ctx context.Context, client *ent.Client) error {
	tags := []string{
		"エンジョイ",
		"男女mix",
		"シニア",
		"ガチ",
		"誰でもok",
		"経験者",
		"初心者歓迎",
		"競技志向",
		"急募",
		"人工芝",
	}
	bulk := make([]*ent.TagCreate, len(tags))
	for i, name := range tags {
		bulk[i] = client.Tag.Create().SetName(name)
	}
	_, err := client.Tag.CreateBulk(bulk...).Save(ctx)
	return err
}

func main() {
	client, _ := db.DatabaseConnection()
	defer client.Close()

	ctx := context.Background()

	err := InsertPrefectures(ctx, client)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	err = InsertCompetitions(ctx, client)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	err = InsertTags(ctx, client)
	if err != nil {
		logger.NewLogger().Fatal(err.Error())
	}

	logger.NewLogger().Info("create initial data!")
}

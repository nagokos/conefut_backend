package message

import (
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

type Message struct {
	Content string
}

func (m Message) MessageValidate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Content,
			validation.Required.Error("メッセージを入力してください"),
			validation.RuneLength(1, 1000).Error("メッセージは1000文字以内で入力してください"),
		),
	)
}

func (m *Message) CreateMessage(ctx context.Context, dbPool *pgxpool.Pool, roomID string) (*model.Message, error) {
	viewer := auth.ForContext(ctx)
	if viewer == nil {
		logger.NewLogger().Error("user not loggedIn")
		return nil, errors.New("ログインしてください")
	}

	cmd := `
	  INSERT INTO messages 
		  (id, content, room_id, user_id, created_at, updated_at)
		VALUES
		  ($1, $2, $3, $4, $5, $6)
		RETURNING content
	`

	timeNow := time.Now().Local()

	var message model.Message
	row := dbPool.QueryRow(
		ctx, cmd,
		xid.New().String(), m.Content, roomID, viewer.ID, timeNow, timeNow,
	)
	err := row.Scan(&message.Content)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &message, nil
}

func GetRoomMessages(ctx context.Context, dbPool *pgxpool.Pool, roomID string) ([]*model.Message, error) {
	viewer := auth.ForContext(ctx)
	if viewer == nil {
		logger.NewLogger().Error("user not loggedIn")
		return nil, errors.New("ログインしてください")
	}

	cmd := `
	  SELECT 
		  m.content, m.created_at,
			u.name,  u.avatar,
			a.message,
			r.title, r.type, r.start_at,
			p.name,
			c.name
		FROM messages AS m 
		INNER JOIN users AS u 
		  ON u.id = m.user_id
		INNER JOIN applicants AS a 
		  ON a.id = m.applicant_id
		INNER JOIN recruitments as r 
	  	ON r.id = a.recruitment_id 
		INNER JOIN prefectures as p 
	  	ON p.id = r.prefecture_id
		INNER JOIN competitions as c 
	  	ON c.id = r.competition_id
		WHERE m.room_id = $1
	`

	rows, err := dbPool.Query(ctx, cmd, roomID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		var message model.Message
		var user model.User
		var applicant model.Applicant
		var recruitment model.Recruitment
		var prefecture model.Prefecture
		var competition model.Competition

		err := rows.Scan(
			&message.Content, &message.CreatedAt, &user.Name, &user.Avatar, &applicant.Message, &recruitment.Title,
			&recruitment.Type, &recruitment.StartAt, &prefecture.Name, &competition.Name,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}

		message.User = &user

		applicant.Recruitment = &recruitment
		message.Applicant = &applicant

		messages = append(messages, &message)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return messages, nil
}

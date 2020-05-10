package storage

import (
	"context"

	"github.com/impu1se/icq_bot/configs"
	"github.com/jackc/pgx/v4"
)

type Database struct {
	connect *pgx.Conn
}

type User struct {
	Id        int64  `db:"id"`
	ChatId    string `db:"chat_id"`
	LastVideo string `db:"last_video"`
	StartTime *int   `db:"start_time"`
	EndTime   *int   `db:"end_time"`
	UserName  string `db:"user_name"`
}

type Message struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
	Text string `db:"text"`
}

func NewDb(config *configs.Config) (*Database, error) {
	conn, err := pgx.Connect(context.Background(), config.Dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &Database{conn}, nil
}

func (db *Database) CreateUser(ctx context.Context, user *User) error {
	_, err := db.connect.Exec(ctx,
		"insert into users (chat_id, user_name) values ($1, $2) on conflict (chat_id) do update set user_name = $2",
		user.ChatId, user.UserName)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetText(ctx context.Context, message string) (string, error) {
	var text string
	err := db.connect.QueryRow(ctx, "select text from messages where name = $1", message).Scan(&text)
	if err != nil {
		return "", err
	}
	return text, nil
}

func (db *Database) GetUser(ctx context.Context, chatId string) (*User, error) {
	var user User
	err := db.connect.QueryRow(ctx, `select id, chat_id, last_video, start_time, end_time, user_name from users where chat_id = $1`, chatId).
		Scan(&user.Id, &user.ChatId, &user.LastVideo, &user.StartTime, &user.EndTime, &user.UserName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Database) ClearTime(ctx context.Context, chatId string) error {
	_, err := db.connect.Exec(ctx, "update users set start_time = null, end_time = null where chat_id = $1", chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateStartTime(ctx context.Context, chatId string, startTime int) error {
	_, err := db.connect.Exec(ctx, "update users set start_time = $1 where chat_id = $2", startTime, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateEndTime(ctx context.Context, chatId string, endTime int) error {
	_, err := db.connect.Exec(ctx, "update users set end_time = $1 where chat_id = $2", endTime, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateLastVideo(ctx context.Context, chatId, lastVideo string) error {
	_, err := db.connect.Exec(ctx, "update users set last_video = $1 where chat_id = $2", lastVideo, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetScale(ctx context.Context) (float64, error) {
	var scale float64
	err := db.connect.QueryRow(ctx, `select scale from settings`).
		Scan(&scale)
	if err != nil {
		return 0, err
	}
	return scale, nil
}

package sms

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type SqliteStorage struct {
	uuid uuid.UUID
	db   *sqlx.DB
}

func NewSqliteStorage(db *sqlx.DB) *SqliteStorage {
	return &SqliteStorage{
		uuid: uuid.NewV4(),
		db:   db,
	}
}

func (s *SqliteStorage) Insert(m Message) error {
	q := "INSERT INTO `sms_messages` (`id`, `time`, `from`, `to`, `text`, `provider`) VALUES (?,?,?,?,?,?)"
	_, err := s.db.Exec(q, m.ID, m.Time, m.From, m.To, m.Text, m.Provider)
	if err != nil {
		return fmt.Errorf("can not insert sms message: %w", err)
	}

	return nil
}

func (s *SqliteStorage) Total() (int64, error) {
	var count int64
	q := "SELECT count(*) FROM `sms_messages`"

	err := s.db.Get(&count, q)
	if err != nil {
		return 0, fmt.Errorf("can not count sms messages: %w", err)
	}

	return count, nil
}

func (s *SqliteStorage) LastMessages(limit int64) ([]Message, error) {

	q := "SELECT `id`, `time`, `from`, `to`, `text`, `provider` FROM `sms_messages` ORDER BY `time` DESC LIMIT ?"
	rows, err := s.db.Query(q, limit)
	if err != nil {
		return nil, fmt.Errorf("can not select %d last messages: %v", limit, err)
	}

	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var id string
		var time time.Time
		var from string
		var to string
		var text string
		var provider string

		err = rows.Scan(&id, &time, &from, &to, &text, &provider)
		if err != nil {
			return nil, fmt.Errorf("can not red row: v", err)
		}

		messages = append(messages, Message{
			ID:       id,
			Time:     time,
			From:     from,
			To:       to,
			Text:     text,
			Provider: provider,
		})
	}

	return messages, nil
}

func (s *SqliteStorage) GenID() string {
	return s.uuid.String()
}

package infrastructure

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/vasyahuyasa/librebread/api"
)

type SQLMSRepo struct {
	db *sqlx.DB
}

func NewSQLSMSRepo(db *sqlx.DB) *SQLMSRepo {
	return &SQLMSRepo{
		db: db,
	}
}

func (s *SQLMSRepo) Create(from, to, text, provider string) (string, error) {
	id := genID()
	now := time.Now()

	q := "INSERT INTO `sms_messages` (`id`, `time`, `from`, `to`, `text`, `provider`) VALUES (?,?,?,?,?,?)"
	_, err := s.db.Exec(q, id, now, from, to, text, provider)
	if err != nil {
		return "", fmt.Errorf("can not insert sms message: %w", err)
	}

	return id, nil
}

func (s *SQLMSRepo) LastMessages(limit int64) (api.SMSList, error) {

	q := "SELECT `id`, `time`, `from`, `to`, `text`, `provider` FROM `sms_messages` ORDER BY `time` DESC LIMIT ?"
	rows, err := s.db.Query(q, limit)
	if err != nil {
		return nil, fmt.Errorf("can not select %d last messages: %v", limit, err)
	}

	defer rows.Close()

	var smses api.SMSList

	for rows.Next() {
		var id string
		var time time.Time
		var from string
		var to string
		var text string
		var provider string

		err = rows.Scan(&id, &time, &from, &to, &text, &provider)
		if err != nil {
			return nil, fmt.Errorf("can not red row: %w", err)
		}

		smses = append(smses, api.SMS{
			ID:       id,
			Time:     time,
			From:     from,
			To:       to,
			Text:     text,
			Provider: provider,
		})
	}

	return smses, nil
}

func genID() string {
	return uuid.NewV4().String()
}

package sms

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vasyahuyasa/librebread/api"
)

type SMSMapper struct {
	//uuid uuid.UUID
	db *sqlx.DB
}

func NewSMSMapper(db *sqlx.DB) *SMSMapper {
	return &SMSMapper{
		//uuid: uuid.NewV4(),
		db: db,
	}
}

/*
func (s *SMSMapper) Insert(m Message) error {
	q := "INSERT INTO `sms_messages` (`id`, `time`, `from`, `to`, `text`, `provider`) VALUES (?,?,?,?,?,?)"
	_, err := s.db.Exec(q, m.ID, m.Time, m.From, m.To, m.Text, m.Provider)
	if err != nil {
		return fmt.Errorf("can not insert sms message: %w", err)
	}

	return nil
}
*/

/*
func (s *SMSMapper) Total() (int64, error) {
	var count int64
	q := "SELECT count(*) FROM `sms_messages`"

	err := s.db.Get(&count, q)
	if err != nil {
		return 0, fmt.Errorf("can not count sms messages: %w", err)
	}

	return count, nil
}
*/

func (s *SMSMapper) LastMessages(limit int64) (api.SMSes, error) {

	q := "SELECT `id`, `time`, `from`, `to`, `text`, `provider` FROM `sms_messages` ORDER BY `time` DESC LIMIT ?"
	rows, err := s.db.Query(q, limit)
	if err != nil {
		return nil, fmt.Errorf("can not select %d last messages: %v", limit, err)
	}

	defer rows.Close()

	var smses api.SMSes

	for rows.Next() {
		var id string
		var time time.Time
		var from string
		var to string
		var text string
		var provider string

		err = rows.Scan(&id, &time, &from, &to, &text, &provider)
		if err != nil {
			return nil, fmt.Errorf("can not red row: %v", err)
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

/*
func (s *SMSMapper) GenID() string {
	return s.uuid.String()
}
*/

package infrastructure

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vasyahuyasa/librebread/api"
)

type SQLHelpdeskRepo struct {
	db *sqlx.DB
}

func NewSQLHelpdeskRepo(db *sqlx.DB) *SQLHelpdeskRepo {
	return &SQLHelpdeskRepo{
		db: db,
	}
}

func (hd *SQLHelpdeskRepo) Create(title, description string, typeID, priorityID, departmentID int) error {
	q := `INSERT INTO helpdesk (title, description, typeID, priorityID, departmentID) VALUES (?,?,?,?,?)`

	_, err := hd.db.Exec(q, title, description, typeID, priorityID, departmentID)
	return err
}

func (hd *SQLHelpdeskRepo) LastTickets(limit int64) (api.HelpdeskEddyTicketList, error) {
	q := "SELECT `id`,`created_at`,`title`,`description`,`type_id`,`priority_id`,`department_id` FROM `helpdesk_messages` ORDER BY `created_at` DESC LIMIT ?"

	rows, err := hd.db.Query(q, limit)
	if err != nil {
		return nil, fmt.Errorf("can not select %d last messages: %w", limit, err)
	}

	defer rows.Close()

	var tickets api.HelpdeskEddyTicketList

	for rows.Next() {
		var id int
		var createdAt time.Time
		var title string
		var description string
		var typeID int
		var priorityId int
		var departmentID int

		err = rows.Scan(&id, &createdAt, &title, &description, &typeID, &priorityId, &departmentID)
		if err != nil {
			return nil, fmt.Errorf("can not red row: %w", err)
		}

		tickets = append(tickets, api.HelpdeskEddyTicket{
			ID:           strconv.Itoa(id),
			Time:         createdAt,
			Title:        title,
			Description:  description,
			TypeId:       typeID,
			PriorityId:   priorityId,
			DepartmentId: departmentID,
		})
	}

	return tickets, nil
}

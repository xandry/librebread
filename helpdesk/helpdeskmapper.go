package helpdesk

import (
	"github.com/jmoiron/sqlx"
)

type Helpdeskmapper struct {
	db *sqlx.DB
}

func NewMapper(db *sqlx.DB) *Helpdeskmapper {
	return &Helpdeskmapper{
		db: db,
	}
}

func (hd *Helpdeskmapper) Create(title, description string, typeID, priorityID, departmentID int) error {
	q := `INSERT INTO helpdesk (title, description, typeID, priorityID, departmentID) VALUES (?,?,?,?,?)`

	_, err := hd.db.Exec(q, title, description, typeID, priorityID, departmentID)
	return err
}

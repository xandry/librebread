package migrations

import "github.com/jmoiron/sqlx"

type m struct {
	name  string
	upSQL string
}

var mList = []m{
	{
		name:  "01_initial",
		upSQL: ``,
	},
}

type Migrator struct {
	db *sqlx.DB
}

func (m *Migrator) createMigrateLog() {
	q := `CREATE TABLE IF NOT EXISTS migrations (name STRING PRIMARY KEY, applied_at DATETIME);`

	m.db.MustExec(q)
}

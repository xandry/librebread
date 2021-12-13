package migrations

import "github.com/jmoiron/sqlx"

type m struct {
	name  string
	upSQL string
}

var mList = []m{
	{
		name: "01_sms",
		upSQL: `CREATE TABLE sms (
			id BLOB NOT NULL,
			time datetime NOT NULL,
			"from" TEXT NOT NULL,
			"to" TEXT NOT NULL,
			"text" TEXT NOT NULL,
			provider TEXT NOT NULL,
			CONSTRAINT sms_PK PRIMARY KEY (id)
		);
		
		CREATE INDEX sms_time_IDX ON sms (time DESC);`,
	},
	{
		name: "02_helpdesk",
		upSQL: `CREATE TABLE helpdesk (
			id TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			type_id INTEGER NOT NULL,
			priority_id INTEGER NOT NULL,
			department_id INTEGER NOT NULL,
			CONSTRAINT helpdesk_PK PRIMARY KEY (id)
		);`,
	},
}

type Migrator struct {
	db *sqlx.DB
}

func (m *Migrator) createMigrateLog() {
	q := `CREATE TABLE IF NOT EXISTS migrations (name STRING PRIMARY KEY, applied_at DATETIME);`

	m.db.MustExec(q)
}

package application

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vasyahuyasa/librebread/api"
	"github.com/vasyahuyasa/librebread/sms"
	"github.com/vasyahuyasa/librebread/web"
)

type Application struct {
	DbPath string
}

func (app *Application) Run() error {
	db, err := app.openSqlite()
	if err != nil {
		return err
	}
	defer db.Close()

	smsMapper := sms.NewSMSMapper(db)

	l := api.NewLibrebread(smsMapper)
	h := api.Handler(l)

	srv := web.NewServer(h)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("web server: %v", err)
	}

	return nil
}

func (app *Application) openSqlite() (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", app.DbPath)
}

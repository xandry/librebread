package main

import (
	"log"

	"github.com/vasyahuyasa/librebread/application"
)

const databasePath = "data.db"

func main() {
	app := application.Application{
		DatabasePath: databasePath,
	}

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

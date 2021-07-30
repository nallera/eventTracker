package main

import (
	"database/sql"
	"eventTracker/cmd/server"
	"eventTracker/internal/db"
	"eventTracker/internal/event"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		panic(fmt.Sprintf("error loading the database: %s", err.Error()))
	}

	env := server.Env{
		EventService: event.EventService{},
		EventDBHandler: db.EventDB{Database: database},
		EventFreqDBHandler: db.EventFreqDB{Database: database},
	}

	server.HandleRequests(env)
}
package main

import (
	"eventTracker/cmd/server"
	"eventTracker/internal/db"
	"eventTracker/internal/event"
)

func main() {
	eventList, eventFreqList := db.DBsStartPoint(true)
	var (
		eventDB = db.EventDB{
			Events: eventList,
			LastID: 4,
		}
		eventFrequenciesDB = db.EventFreqDB{
			EventFrequencies: eventFreqList,
			LastID:           0,
		}
	)

	env := server.Env{
		EventService: event.EventService{
			EventDB:            &eventDB,
			EventFrequenciesDB: &eventFrequenciesDB,
		},
	}

	e := env.EventService.StartDBs()
	if e != nil {
		panic("error starting dbs")
	}

	server.HandleRequests(env)
}
package main

import (
	"eventTracker/cmd/server"
	event "eventTracker/internal"
)

func main() {
	var (
		eventDB = event.EventDB{}
		eventFrequenciesDB = event.EventFreqDB{}
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

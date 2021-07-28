package server

import (
	event "eventTracker/internal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Env struct {
	EventService event.EventService
}

func HandleRequests(env Env) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/events/{name}", env.ReturnEvent).Methods("GET")
	router.HandleFunc("/events", env.ReturnAllEvents).Methods("GET")
	router.HandleFunc("/events", env.ReturnEventsInDateRange).Methods("GET").Queries("start_date", "{[0-9]+}", "end_date", "{[0-9]+}")

	router.HandleFunc("/events/{name}", env.CreateEvent).Methods("POST") //N in body

	router.HandleFunc("/event_frequencies/{name}", env.GetEventFrequency).Methods("GET")
	router.HandleFunc("/event_frequencies", env.GetAllEventsFrequencies).Methods("GET")

	log.Fatal(http.ListenAndServe(":10000", router))
}
package server

import (
	"eventTracker/internal/event"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Env struct {
	EventService event.EventServiceI
}

func HandleRequests(env Env) {
	router := mux.NewRouter().StrictSlash(true)
	sub := router.PathPrefix("/api/v1").Subrouter()

	sub.HandleFunc("/events", env.ReturnEvents).Methods("GET")
	sub.HandleFunc("/events/{name}", env.ReturnEvent).Methods("GET")

	sub.HandleFunc("/events/{name}", env.CreateEvent).Methods("POST") //N in body

	sub.HandleFunc("/event_frequencies/{name}", env.ReturnEventFrequency).Methods("GET")
	sub.HandleFunc("/event_frequencies", env.ReturnAllEventsFrequencies).Methods("GET")

	sub.HandleFunc("/event_frequencies/{name}/hist", env.ReturnEventFrequencyHistogram).Methods("GET")

	log.Fatal(http.ListenAndServe(":10000", router))
}
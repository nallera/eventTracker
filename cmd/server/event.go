package server

import (
	"encoding/json"
	"errors"
	event "eventTracker/internal"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

func (env Env) ReturnAllEvents(w http.ResponseWriter, r *http.Request) {
	retrievedEvents, err := env.EventService.AllEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(retrievedEvents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (env Env) ReturnEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	retrievedEvent, err := env.EventService.EventsByName(name)
	if errors.Is(err, event.ErrEventNotFound) {
		http.Error(w, fmt.Sprintf(err.Error(), name), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(retrievedEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (env Env) ReturnEventsInDateRange(w http.ResponseWriter, r *http.Request) {

}

func (env Env) CreateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	var body event.EventBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		body.Count = 1
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Json decoder error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	err = env.EventService.CreateEvent(name, body.Count, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		println(fmt.Sprintf("error: %v", err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (env Env) ReturnEventFrequency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	retrievedEvent, err := env.EventService.EventFrequencyByName(name)
	if errors.Is(err, event.ErrEventNotFound) {
		http.Error(w, fmt.Sprintf(err.Error(), name), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(retrievedEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (env Env) ReturnAllEventsFrequencies(w http.ResponseWriter, r *http.Request) {
	retrievedEvents, err := env.EventService.AllEventsFrequencies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(retrievedEvents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
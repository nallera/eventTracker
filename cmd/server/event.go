package server

import (
	"encoding/json"
	"errors"
	"eventTracker/internal/model"
	"eventTracker/internal/plotting"
	"fmt"
	"github.com/gorilla/mux"
	"gonum.org/v1/plot/plotter"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (env Env) ReturnEvents(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	println(fmt.Sprintf("query: %v", queryParams))

	var (
		retrievedEvents []model.Event
		err error
	)
	if len(queryParams) > 0 {
		retrievedEvents, err = env.EventService.EventsByDateRange(env.EventDBHandler, queryParams["start_date"][0], queryParams["end_date"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		retrievedEvents, err = env.EventService.AllEvents(env.EventDBHandler)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	retrievedEvent, err := env.EventService.EventsByName(env.EventDBHandler, name)
	if errors.Is(err, model.ErrEventNotFound) {
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

func (env Env) CreateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	var body model.EventBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		body.Count = 1
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Json decoder error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	err = env.EventService.CreateEvent(env.EventDBHandler, env.EventFreqDBHandler, name, body.Count, time.Now())
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

	retrievedEvent, err := env.EventService.EventFrequencyByName(env.EventFreqDBHandler, name)
	if errors.Is(err, model.ErrEventNotFound) {
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
	retrievedEvents, err := env.EventService.AllEventsFrequencies(env.EventFreqDBHandler)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(retrievedEvents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (env Env) ReturnEventFrequencyHistogram(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	retrievedEvent, err := env.EventService.EventFrequencyByName(env.EventFreqDBHandler, name)
	if errors.Is(err, model.ErrEventNotFound) {
		http.Error(w, fmt.Sprintf(err.Error(), name), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//var histImage *image.Image
	var values plotter.Values
	for _, h := range retrievedEvent.HourCount {
		values = append(values, 100 * float64(h) / float64(retrievedEvent.TotalCount))
	}
	writer := plotting.PlotHistogram(values, retrievedEvent.Name)


	imgBytes, err := writer.WriteTo(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(int(imgBytes)))
}
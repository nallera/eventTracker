package server

import (
	"encoding/json"
	"eventTracker/config"
	"eventTracker/internal/db"
	"eventTracker/internal/event"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type Env struct {
	EventService event.EventServiceI
	EventDBHandler db.EventDBHandler
	EventFreqDBHandler db.EventFreqDBHandler
}

func HandleRequests(env Env) {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(AuthMiddleware)

	healthRoute := router.PathPrefix("/health").Subrouter()
	healthRoute.HandleFunc("/ping", pingCheck)

	apiRoute := router.PathPrefix("/api/v1").Subrouter()

	apiRoute.HandleFunc("/events", env.ReturnEvents).Methods("GET")

	apiRoute.HandleFunc("/events/{name}", env.CreateEvent).Methods("POST") //N and date in body

	apiRoute.HandleFunc("/event_history", env.ReturnAllEventsHistory).Methods("GET")

	apiRoute.HandleFunc("/event_frequencies/{name}/hist", env.ReturnEventFrequencyHistogram).Methods("GET")

	adminRoute := router.PathPrefix("/admin/v1").Subrouter()
	adminRoute.HandleFunc("/events/{name}", env.ReturnEvent).Methods("GET")
	adminRoute.HandleFunc("/events/{name}", env.DeleteEvent).Methods("DELETE")
	adminRoute.HandleFunc("/event_frequencies/{name}", env.ReturnEventFrequency).Methods("GET")
	adminRoute.HandleFunc("/event_frequencies", env.ReturnAllEventsFrequencies).Methods("GET")

	log.Fatal(http.ListenAndServe(":10000", router))
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")

		if apiKey == "" {
			w.WriteHeader(http.StatusForbidden)
			err :=json.NewEncoder(w).Encode("Missing auth apiKey")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				println(fmt.Sprintf("error: %v", err.Error()))
				return
			}
			return
		}

		var apiKeysToCheck []string
		if strings.Split(r.URL.Path, "/")[1] == "admin" {
			apiKeysToCheck = config.AdminApiKeys
		} else {
			apiKeysToCheck = config.ApiKeys
		}

		authorized := false
		for _,t := range apiKeysToCheck {
			if apiKey == t {
				authorized = true
			}
		}

		if authorized {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode("Wrong auth apiKey")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				println(fmt.Sprintf("error: %v", err.Error()))
				return
			}
		}
	})
}

func pingCheck(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode("ping ok")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		println(fmt.Sprintf("error: %v", err.Error()))
		return
	}
	return
}

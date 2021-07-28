package internal

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrEventNotFound = errors.New("event %s not found")
	ErrStartEventDB = errors.New("error starting event db")
	ErrStartEventFreqDB = errors.New("error starting event freq db")
)

type Event struct {
	ID 		uint64 `json:"id"`
	Name	string `json:"event"`
	Count	uint64 `json:"count"`
	Date    string `json:"date"`
}

type EventFreq struct {
	ID 		uint64 `json:"id"`
	Name	string `json:"event"`
	TotalCount	uint64 `json:"count"`
	HourCount   [24]uint64 `json:"hour_count"`
}

type EventBody struct {
	Count uint64 `json:"count,omitempty"`
}

type EventService struct {
	EventDB            *EventDB
	EventFrequenciesDB *EventFreqDB
}

func (es EventService) StartDBs() (err error) {
	e := es.EventDB.StartDB()
	if e != nil {
		return ErrStartEventDB
	}

	e = es.EventFrequenciesDB.StartDB()
	if e != nil {
		return ErrStartEventFreqDB
	}

	return nil
}

func (es EventService) EventsByName(name string) (events []Event, err error) {
	events, e := es.EventDB.getEventsByName(name)
	if e != nil {
		return nil, ErrEventNotFound
	}

	return events, nil
}

func (es EventService) EventFrequencyByName(name string) (eventFreq *EventFreq, err error) {
	eventFreq, e := es.EventFrequenciesDB.getEventByName(name)
	if e != nil {
		return nil, ErrEventNotFound
	}

	return eventFreq, nil
}

func (es EventService) AllEvents() (events []Event, err error) {
	events, e := es.EventDB.getEvents()
	if e != nil {
		return []Event{}, nil
	}

	return events, nil
}

func (es EventService) AllEventsFrequencies() (events []EventFreq, err error) {
	events, e := es.EventFrequenciesDB.getEvents()
	if e != nil {
		return []EventFreq{}, nil
	}

	return events, nil
}

func (es EventService) EventByID(ID uint64) (event *Event, err error) {
	event, e := es.EventDB.getEventByID(ID)
	if e != nil {
		return nil, ErrEventNotFound
	}

	return event, nil
}

func (es EventService) CreateEvent(name string, count uint64, date time.Time) (err error) {
	dateYYYYmmdd := date.Format("20060102")

	event, e := es.EventDB.getEventByNameAndDate(name, dateYYYYmmdd)
	if e != nil {
		println(fmt.Sprintf("Creating new event %s", name))

		_, e := es.EventDB.CreateEvent(name, count, dateYYYYmmdd)
		if e != nil {
			return ErrInsertEventDB
		}
	} else {
		println(fmt.Sprintf("Updating event %s", name))

		e := es.EventDB.UpdateEvent(event.ID, count)
		if e != nil {
			return ErrUpdateEventDB
		}
	}
	ev,_ := es.EventDB.getEvents()
	println(fmt.Sprintf("eventDB: %v", ev))

	hour := date.Format("15")
	hourUint, e := strconv.ParseUint(hour, 10, 8)
	if e != nil {
		return ErrParseHour
	}

	eventFreq, e := es.EventFrequenciesDB.getEventByName(name)
	if e != nil {
		println(fmt.Sprintf("Creating new event %s frequency", name))

		_, e = es.EventFrequenciesDB.CreateEvent(name, count, hourUint)
		if e != nil {
			return ErrInsertEventFreqDB
		}
	} else {
		println(fmt.Sprintf("Updating event %s frequency", name))

		e = es.EventFrequenciesDB.UpdateEvent(eventFreq.ID, count, hourUint)
		if e != nil {
			return ErrUpdateEventFreqDB
		}
	}
	evf,_ := es.EventFrequenciesDB.getEvents()
	println(fmt.Sprintf("eventFreqDB: %v", evf))

	return nil
}
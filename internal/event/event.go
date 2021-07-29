package event

import (
	"eventTracker/internal/db"
	"eventTracker/internal/model"
	"fmt"
	"strconv"
	"time"
)

type EventServiceI interface {
	StartDBs() (err error)
 	EventsByName(name string) (events []model.Event, err error)
 	EventsByDateRange(startDate, endDate string) (events []model.Event, err error)
 	EventFrequencyByName(name string) (eventFreq *model.EventFreq, err error)
 	AllEvents() (events []model.Event, err error)
 	AllEventsFrequencies() (events []model.EventFreq, err error)
 	EventByID(ID uint64) (event *model.Event, err error)
 	CreateEvent(name string, count uint64, date time.Time) (err error)
}

type EventService struct {
	EventDB            *db.EventDB
	EventFrequenciesDB *db.EventFreqDB
}

func (es EventService) StartDBs() (err error) {
	e := es.EventDB.StartDB()
	if e != nil {
		return model.ErrStartEventDB
	}

	e = es.EventFrequenciesDB.StartDB()
	if e != nil {
		return model.ErrStartEventFreqDB
	}

	return nil
}

func (es EventService) EventsByName(name string) (events []model.Event, err error) {
	events, e := es.EventDB.GetEventsByName(name)
	if e != nil {
		return nil, model.ErrEventNotFound
	}

	return events, nil
}

func (es EventService) EventsByDateRange(startDate, endDate string) (events []model.Event, err error) {
	events, e := es.EventDB.GetEventsByDateRange(startDate, endDate)
	if e != nil {
		return []model.Event{}, nil
	}

	return events, nil
}

func (es EventService) EventFrequencyByName(name string) (eventFreq *model.EventFreq, err error) {
	eventFreq, e := es.EventFrequenciesDB.GetEventByName(name)
	if e != nil {
		return nil, model.ErrEventNotFound
	}

	return eventFreq, nil
}

func (es EventService) AllEvents() (events []model.Event, err error) {
	events, e := es.EventDB.GetEvents()
	if e != nil {
		return []model.Event{}, nil
	}

	return events, nil
}

func (es EventService) AllEventsFrequencies() (events []model.EventFreq, err error) {
	events, e := es.EventFrequenciesDB.GetEvents()
	if e != nil {
		return []model.EventFreq{}, nil
	}

	return events, nil
}

func (es EventService) EventByID(ID uint64) (event *model.Event, err error) {
	event, e := es.EventDB.GetEventByID(ID)
	if e != nil {
		return nil, model.ErrEventNotFound
	}

	return event, nil
}

func (es EventService) CreateEvent(name string, count uint64, date time.Time) (err error) {
	dateYYYYmmdd := date.Format("20060102")

	event, e := es.EventDB.GetEventByNameAndDate(name, dateYYYYmmdd)
	if e != nil {
		println(fmt.Sprintf("Creating new event %s", name))

		_, e := es.EventDB.CreateEvent(name, count, dateYYYYmmdd)
		if e != nil {
			return model.ErrInsertEventDB
		}
	} else {
		println(fmt.Sprintf("Updating event %s", name))

		e := es.EventDB.UpdateEvent(event.ID, count)
		if e != nil {
			return model.ErrUpdateEventDB
		}
	}
	ev,_ := es.EventDB.GetEvents()
	println(fmt.Sprintf("eventDB: %v", ev))

	hour := date.Format("15")
	hourUint, e := strconv.ParseUint(hour, 10, 8)
	if e != nil {
		return model.ErrParseHour
	}

	eventFreq, e := es.EventFrequenciesDB.GetEventByName(name)
	if e != nil {
		println(fmt.Sprintf("Creating new event %s frequency", name))

		_, e = es.EventFrequenciesDB.CreateEvent(name, count, hourUint)
		if e != nil {
			return model.ErrInsertEventFreqDB
		}
	} else {
		println(fmt.Sprintf("Updating event %s frequency", name))

		e = es.EventFrequenciesDB.UpdateEvent(eventFreq.ID, count, hourUint)
		if e != nil {
			return model.ErrUpdateEventFreqDB
		}
	}
	evf,_ := es.EventFrequenciesDB.GetEvents()
	println(fmt.Sprintf("eventFreqDB: %v", evf))

	return nil
}
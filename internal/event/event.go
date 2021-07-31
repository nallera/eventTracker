package event

import (
	"errors"
	"eventTracker/internal/db"
	"eventTracker/internal/model"
	"fmt"
	"strconv"
	"time"
)

type EventServiceI interface {
 	EventsByName(EventDBHandler db.EventDBHandler, name string) (events []model.Event, err error)
 	EventsByDateRange(EventDBHandler db.EventDBHandler, startDate, endDate string) (events []model.Event, err error)
 	AllEvents(EventDBHandler db.EventDBHandler) (events []model.Event, err error)
 	EventByID(EventDBHandler db.EventDBHandler, ID uint64) (event model.Event, err error)
 	CreateEvent(EventDBHandler db.EventDBHandler, EventDBFreqHandler db.EventFreqDBHandler, name string, count uint64, date time.Time) (err error)
 	DeleteEvent(EventDBHandler db.EventDBHandler, EventDBFreqHandler db.EventFreqDBHandler, name string) (err error)
	EventFrequencyByName(EventDBFreqHandler db.EventFreqDBHandler, name string) (eventFreq model.EventFreq, err error)
	AllEventsFrequencies(EventDBFreqHandler db.EventFreqDBHandler) (events []model.EventFreq, err error)
	AllEventsHistory(EventDBFreqHandler db.EventFreqDBHandler) (events []model.EventHistory, err error)
}

type EventService struct {}

func (es EventService) EventsByName(EventDBHandler db.EventDBHandler, name string) (events []model.Event, err error) {
	events, e := EventDBHandler.GetEventsByName(name)
	if e != nil {
		return nil, model.ErrEventNotFound
	}

	return events, nil
}

func (es EventService) EventsByDateRange(EventDBHandler db.EventDBHandler, startDate, endDate string) (events []model.Event, err error) {
	events, e := EventDBHandler.GetEventsByDateRange(startDate, endDate)
	if e != nil {
		return []model.Event{}, nil
	}

	return events, nil
}

func (es EventService) AllEvents(EventDBHandler db.EventDBHandler) (events []model.Event, err error) {
	events, e := EventDBHandler.GetEvents()
	if e != nil {
		return []model.Event{}, nil
	}

	return events, nil
}

func (es EventService) EventByID(EventDBHandler db.EventDBHandler, ID uint64) (event model.Event, err error) {
	event, e := EventDBHandler.GetEventByID(ID)
	if e != nil {
		return model.Event{}, model.ErrEventNotFound
	}

	return event, nil
}

func (es EventService) CreateEvent(EventDBHandler db.EventDBHandler, EventDBFreqHandler db.EventFreqDBHandler, name string, count uint64, date time.Time) (err error) {
	dateYYYYmmdd := date.Format("2006-01-02")

	event, e := EventDBHandler.GetEventByNameAndDate(name, dateYYYYmmdd)
	if errors.Is(e, model.ErrEventNotFound) {
		println(fmt.Sprintf("Creating new event %s", name))

		e := EventDBHandler.CreateEvent(name, count, dateYYYYmmdd)
		if e != nil {
			return errors.New(fmt.Sprintf(model.ErrInsertEventDB.Error(), e.Error()))
		}
	} else if e != nil {
		return errors.New(fmt.Sprintf("error getting event by name and date: %s", e))
	} else {
		println(fmt.Sprintf("Updating event %s", name))

		e := EventDBHandler.UpdateEvent(event.ID, count)
		if e != nil {
			return errors.New(fmt.Sprintf(model.ErrUpdateEventDB.Error(), e.Error()))
		}
	}

	hour := date.Format("15")
	hourUint, e := strconv.ParseUint(hour, 10, 8)
	if e != nil {
		return model.ErrParseHour
	}

	eventFreq, e := EventDBFreqHandler.GetEventByName(name)
	if errors.Is(e, model.ErrEventNotFound) {
		println(fmt.Sprintf("Creating new event %s frequency", name))

		e = EventDBFreqHandler.CreateEvent(name, count, hourUint)
		if e != nil {
			return errors.New(fmt.Sprintf(model.ErrInsertEventFreqDB.Error(), e.Error()))
		}
	} else if e != nil{
		return errors.New(fmt.Sprintf("error getting event freq by name: %s", e))
	} else {
		println(fmt.Sprintf("Updating event %s frequency", name))

		e = EventDBFreqHandler.UpdateEvent(eventFreq.ID, count, hourUint)
		if e != nil {
			return errors.New(fmt.Sprintf(model.ErrUpdateEventFreqDB.Error(), e.Error()))
		}
	}

	return nil
}

func (es EventService) DeleteEvent(EventDBHandler db.EventDBHandler, EventDBFreqHandler db.EventFreqDBHandler, name string) (err error) {
	IDsToDelete, e := EventDBHandler.GetEventsIDsByName(name)
	if errors.Is(e, model.ErrEventNotFound) {
		return errors.New(fmt.Sprintf(model.ErrDoesntExistEventDB.Error(), name))
	} else if e != nil {
		return errors.New(fmt.Sprintf("error getting event by name and date: %s", e))
	} else {
		println(fmt.Sprintf("Deleting event %s", name))

		e = EventDBHandler.DeleteEvents(IDsToDelete)
		if e != nil {
			return errors.New(fmt.Sprintf(model.ErrDeleteEventDB.Error(), e.Error()))
		}
	}

	eventFreq, e := EventDBFreqHandler.GetEventByName(name)
	if errors.Is(e, model.ErrEventNotFound) {
		return errors.New(fmt.Sprintf(model.ErrDoesntExistEventFreqDB.Error(), name))
	} else if e != nil{
		return errors.New(fmt.Sprintf("error getting event freq by name: %s", e))
	} else {
		println(fmt.Sprintf("Deleting event %s frequency", name))

		e = EventDBFreqHandler.DeleteEvent(eventFreq.ID)
		if e != nil {
			return errors.New(fmt.Sprintf(model.ErrDeleteEventFreqDB.Error(), e.Error()))
		}
	}

	return nil
}

func (es EventService) AllEventsFrequencies(EventDBFreqHandler db.EventFreqDBHandler) (events []model.EventFreq, err error) {
	events, e := EventDBFreqHandler.GetEvents()
	if e != nil {
		return []model.EventFreq{}, nil
	}

	return events, nil
}

func (es EventService) AllEventsHistory(EventDBFreqHandler db.EventFreqDBHandler) (events []model.EventHistory, err error) {
	events, e := EventDBFreqHandler.GetEventsHistory()
	if e != nil {
		return []model.EventHistory{}, nil
	}

	return events, nil
}

func (es EventService) EventFrequencyByName(EventDBFreqHandler db.EventFreqDBHandler, name string) (eventFreq model.EventFreq, err error) {
	eventFreq, e := EventDBFreqHandler.GetEventByName(name)
	if e != nil {
		return model.EventFreq{}, model.ErrEventNotFound
	}

	return eventFreq, nil
}
package internal

import (
	"errors"
)

var (
	ErrInsertEventDB = errors.New("error inserting new event in event db")
	ErrInsertEventFreqDB = errors.New("error inserting new event in event freq db")
	ErrUpdateEventDB = errors.New("error updating new event in event db")
	ErrUpdateEventFreqDB = errors.New("error updating new event in event freq db")
	ErrParseHour = errors.New("error parsing hour into int")
)

type EventDB struct {
	Events  []*Event
	LastID  uint64
}

type EventFreqDB struct {
	EventFrequencies  []*EventFreq
	LastID            uint64
}

func (db *EventDB) StartDB() (err error) {
	db.LastID = 0
	return nil
}

func (db *EventDB) getEvents() (events []Event, err error) {
	for _, ev := range db.Events {
		events = append(events, *ev)
	}
	return events, nil
}

func (db *EventDB) getEventsByName(name string) (events []Event, err error) {
	for _, e := range db.Events {
		if e.Name == name {
			events = append(events, *e)
		}
	}

	if events == nil {
		return nil, ErrEventNotFound
	}

	return events, nil
}

func (db *EventDB) getEventByNameAndDate(name, date string) (event *Event, err error) {
	for _, e := range db.Events {
		if e.Name == name && e.Date == date{
			return e, nil
		}
	}

	return nil, ErrEventNotFound
}

func (db *EventDB) getEventByID(ID uint64) (event *Event, err error) {
	for _, e := range db.Events {
		if e.ID == ID {
			return e, nil
		}
	}

	return nil, ErrEventNotFound
}

func (db *EventDB) CreateEvent(name string, count uint64, date string) (event *Event, err error) {
	newEvent := Event{
		ID:    db.LastID + 1,
		Name:  name,
		Count: count,
		Date:  date,
	}

	db.Events = append(db.Events, &newEvent)

	event, e := db.getEventByID(db.LastID + 1)
	if e != nil {
		return nil, ErrInsertEventDB
	}

	db.LastID += 1

	return event, nil
}

func (db *EventDB) UpdateEvent(ID, count uint64) (err error) {
	event, e := db.getEventByID(ID)
	if e != nil {
		return ErrUpdateEventDB
	}

	event.Count += count

	return nil
}

func (db *EventFreqDB) StartDB() (err error) {
	db.LastID = 0

	return nil
}

func (db *EventFreqDB) getEvents() (events []EventFreq, err error) {
	for _, ev := range db.EventFrequencies {
		events = append(events, *ev)
	}
	return events, nil
}

func (db *EventFreqDB) getEventByID(ID uint64) (event *EventFreq, err error) {
	for _, e := range db.EventFrequencies {
		if e.ID == ID {
			return e, nil
		}
	}

	return nil, ErrEventNotFound
}

func (db *EventFreqDB) getEventByName(name string) (event *EventFreq, err error) {
	for _, e := range db.EventFrequencies {
		if e.Name == name {
			return e, nil
		}
	}

	return nil, ErrEventNotFound
}

func (db *EventFreqDB) CreateEvent(name string, count uint64, hour uint64) (event *EventFreq, err error) {
	var hourCount [24]uint64
	hourCount[hour] = count

	newEventFreq := EventFreq{
		ID:    db.LastID + 1,
		Name:  name,
		TotalCount: count,
		HourCount:  hourCount,
	}

	db.EventFrequencies = append(db.EventFrequencies, &newEventFreq)

	event, e := db.getEventByID(db.LastID + 1)
	if e != nil {
		return nil, ErrInsertEventFreqDB
	}

	db.LastID += 1

	return event, nil
}

func (db *EventFreqDB) UpdateEvent(ID, count, hour uint64) (err error) {
	event, e := db.getEventByID(ID)
	if e != nil {
		return ErrUpdateEventFreqDB
	}

	event.TotalCount += count
	event.HourCount[hour] += count

	return nil
}
package db

import (
	"eventTracker/internal/model"
	"time"
)

type EventDB struct {
	Events  []*model.Event
	LastID  uint64
}

type EventFreqDB struct {
	EventFrequencies  []*model.EventFreq
	LastID            uint64
}

func (db *EventDB) StartDB() (err error) {
	db.LastID = 0
	return nil
}

func (db *EventDB) GetEvents() (events []model.Event, err error) {
	for _, ev := range db.Events {
		events = append(events, *ev)
	}
	return events, nil
}

func (db *EventDB) GetEventsByName(name string) (retrievedEvents []model.Event, err error) {
	for _, ev := range db.Events {
		if ev.Name == name {
			retrievedEvents = append(retrievedEvents, *ev)
		}
	}

	if retrievedEvents == nil {
		return nil, model.ErrEventNotFound
	}

	return retrievedEvents, nil
}

func (db *EventDB) GetEventsByDateRange(startDate, endDate string) (retrievedEvents []model.Event, err error) {
	startDateTime, e := time.Parse("20060102", startDate)
	if e != nil {
		return nil, model.ErrParseDate
	}
	endDateTime, e := time.Parse("20060102", endDate)
	if e != nil {
		return nil, model.ErrParseDate
	}

	for _, ev := range db.Events {
		eventDate, e := time.Parse("20060102", ev.Date)
		if e != nil {
			return nil, model.ErrParseDate
		}
		if eventDate.After(startDateTime) && eventDate.Before(endDateTime) {
			retrievedEvents = append(retrievedEvents, *ev)
		}
	}

	if retrievedEvents == nil {
		return nil, model.ErrEventNotFound
	}

	return retrievedEvents, nil
}

func (db *EventDB) GetEventByNameAndDate(name, date string) (retrievedEvent *model.Event, err error) {
	for _, ev := range db.Events {
		if ev.Name == name && ev.Date == date{
			return ev, nil
		}
	}

	return nil, model.ErrEventNotFound
}

func (db *EventDB) GetEventByID(ID uint64) (retrievedEvent *model.Event, err error) {
	for _, ev := range db.Events {
		if ev.ID == ID {
			return ev, nil
		}
	}

	return nil, model.ErrEventNotFound
}

func (db *EventDB) CreateEvent(name string, count uint64, date string) (retrievedEvents *model.Event, err error) {
	newEvent := model.Event{
		ID:    db.LastID + 1,
		Name:  name,
		Count: count,
		Date:  date,
	}

	db.Events = append(db.Events, &newEvent)

	retrievedEvents, e := db.GetEventByID(db.LastID + 1)
	if e != nil {
		return nil, model.ErrInsertEventDB
	}

	db.LastID += 1

	return retrievedEvents, nil
}

func (db *EventDB) UpdateEvent(ID, count uint64) (err error) {
	event, ev := db.GetEventByID(ID)
	if ev != nil {
		return model.ErrUpdateEventDB
	}

	event.Count += count

	return nil
}

func (db *EventFreqDB) StartDB() (err error) {
	db.LastID = 0

	return nil
}

func (db *EventFreqDB) GetEvents() (retrievedEvents []model.EventFreq, err error) {
	for _, ev := range db.EventFrequencies {
		retrievedEvents = append(retrievedEvents, *ev)
	}
	return retrievedEvents, nil
}

func (db *EventFreqDB) GetEventByID(ID uint64) (retrievedEvent *model.EventFreq, err error) {
	for _, ev := range db.EventFrequencies {
		if ev.ID == ID {
			return ev, nil
		}
	}

	return nil, model.ErrEventNotFound
}

func (db *EventFreqDB) GetEventByName(name string) (retrievedEvent *model.EventFreq, err error) {
	for _, ev := range db.EventFrequencies {
		if ev.Name == name {
			return ev, nil
		}
	}

	return nil, model.ErrEventNotFound
}

func (db *EventFreqDB) CreateEvent(name string, count uint64, hour uint64) (retrievedEvent *model.EventFreq, err error) {
	var hourCount [24]uint64
	hourCount[hour] = count

	newEventFreq := model.EventFreq{
		ID:    db.LastID + 1,
		Name:  name,
		TotalCount: count,
		HourCount:  hourCount,
	}

	db.EventFrequencies = append(db.EventFrequencies, &newEventFreq)

	retrievedEvent, e := db.GetEventByID(db.LastID + 1)
	if e != nil {
		return nil, model.ErrInsertEventFreqDB
	}

	db.LastID += 1

	return retrievedEvent, nil
}

func (db *EventFreqDB) UpdateEvent(ID, count, hour uint64) (err error) {
	event, e := db.GetEventByID(ID)
	if e != nil {
		return model.ErrUpdateEventFreqDB
	}

	event.TotalCount += count
	event.HourCount[hour] += count

	return nil
}

func DBsStartPoint(debug bool) ([]*model.Event, []*model.EventFreq) {
	if debug {
		eventList := []*model.Event{
			{1, "login1", 3, "20210203"},
			{1, "login1", 6, "20210603"},
			{1, "login1", 2, "20211207"},
			{2, "login2", 54, "20210103"},
			{2, "login2", 43, "20210203"},
			{2, "login2", 32, "20210223"},
			{2, "login2", 12, "20211025"},
			{3, "logout1", 13, "20210809"},
			{3, "logout1", 8, "20210811"},
			{4, "logout2", 1, "20211213"},
			{4, "logout2", 7, "20211223"},
			{4, "logout2", 2, "20210805"},
		}
		eventFreqList := []*model.EventFreq{
			{1, "login1", 11, [24]uint64{0, 0, 0, 1, 2, 0, 3, 0, 3, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}},
			{2, "login2", 54 + 43 + 32 + 12, [24]uint64{0, 0, 23, 0, 12, 11, 0, 10, 25, 10, 10, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 10, 10, 0}},
			{3, "logout1", 21, [24]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 10, 1}},
			{4, "logout2", 10, [24]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0}},
		}
		return eventList, eventFreqList
	}
	return nil, nil
}

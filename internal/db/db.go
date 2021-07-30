package db

import (
	"database/sql"
	"encoding/json"
	"eventTracker/internal/model"
	"fmt"
)

type EventDBHandler interface {
	GetEvents() (events []model.Event, err error)
	GetEventsByName(name string) (retrievedEvents []model.Event, err error)
	GetEventsByDateRange(startDate, endDate string) (retrievedEvents []model.Event, err error)
	GetEventByNameAndDate(name, date string) (retrievedEvent model.Event, err error)
	GetEventByID(ID uint64) (retrievedEvent model.Event, err error)
	CreateEvent(name string, count uint64, date string) (err error)
	UpdateEvent(ID, count uint64) (err error)
}

type EventFreqDBHandler interface {
	GetEvents() (retrievedEvents []model.EventFreq, err error)
	GetEventByID(ID uint64) (retrievedEvent model.EventFreq, err error)
	GetEventByName(name string) (retrievedEvent model.EventFreq, err error)
	CreateEvent(name string, count uint64, hour uint64) (err error)
	UpdateEvent(ID, count, hour uint64) (err error)
}

type EventDB struct {
	Database *sql.DB
}

type EventFreqDB struct {
	Database *sql.DB
}

func (db EventDB) GetEvents() (retrievedEvents []model.Event, err error) {
	rows, e := db.Database.Query("SELECT * FROM eventDB")
	if e != nil {
		return nil, e
	}

	for rows.Next() {
		var event model.Event

		e = rows.Scan(&event.ID, &event.Date, &event.Name, &event.Count)
		if e != nil {
			return nil, e
		}

		retrievedEvents = append(retrievedEvents, event)
	}

	e = rows.Close()
	if e != nil {
		return nil, e
	}

	return retrievedEvents, nil
}

func (db EventDB) GetEventsByName(name string) (retrievedEvents []model.Event, err error) {
	rows, e := db.Database.Query(fmt.Sprintf("SELECT * FROM eventDB WHERE name LIKE \"%s\"", name))
	if e != nil {
		return nil, e
	}

	for rows.Next() {
		var event model.Event

		e = rows.Scan(&event.ID, &event.Date, &event.Name, &event.Count)
		if e != nil {
			return nil, e
		}

		retrievedEvents = append(retrievedEvents, event)
	}

	e = rows.Close()
	if e != nil {
		return nil, e
	}

	if retrievedEvents == nil {
		return nil, model.ErrEventNotFound
	}

	return retrievedEvents, nil
}

func (db EventDB) GetEventsByDateRange(startDate, endDate string) (retrievedEvents []model.Event, err error) {
	rows, e := db.Database.Query(fmt.Sprintf("SELECT * FROM eventDB WHERE date BETWEEN \"%s\" and \"%s\"", startDate, endDate))
	if e != nil {
		return nil, e
	}

	for rows.Next() {
		var event model.Event

		e = rows.Scan(&event.ID, &event.Date, &event.Name, &event.Count)
		if e != nil {
			return nil, e
		}

		retrievedEvents = append(retrievedEvents, event)
	}

	e = rows.Close()
	if e != nil {
		return nil, e
	}

	if retrievedEvents == nil {
		return nil, model.ErrEventNotFound
	}

	return retrievedEvents, nil
}

func (db EventDB) GetEventByNameAndDate(name, date string) (retrievedEvent model.Event, err error) {
	rows, e := db.Database.Query(fmt.Sprintf("SELECT * FROM eventDB WHERE name LIKE \"%s\" AND date LIKE \"%s\"", name, date))
	if e != nil {
		return model.Event{}, e
	}

	for rows.Next() {
		e = rows.Scan(&retrievedEvent.ID, &retrievedEvent.Date, &retrievedEvent.Name, &retrievedEvent.Count)
		if e != nil {
			return model.Event{}, e
		}
	}

	e = rows.Close()
	if e != nil {
		return model.Event{}, e
	}

	if retrievedEvent.ID == 0 {
		return model.Event{}, model.ErrEventNotFound
	}
	return retrievedEvent, nil
}

func (db EventDB) GetEventByID(ID uint64) (retrievedEvent model.Event, err error) {
	rows, e := db.Database.Query(fmt.Sprintf("SELECT * FROM eventDB WHERE ID=%d", ID))
	if e != nil {
		return model.Event{}, e
	}

	for rows.Next() {
		e = rows.Scan(&retrievedEvent.ID, &retrievedEvent.Date, &retrievedEvent.Name, &retrievedEvent.Count)
		if e != nil {
			return model.Event{}, e
		}
	}

	e = rows.Close()
	if e != nil {
		return model.Event{}, e
	}

	if retrievedEvent.ID == 0 {
		return model.Event{}, model.ErrEventNotFound
	}
	return retrievedEvent, nil
}

func (db EventDB) CreateEvent(name string, count uint64, date string) (err error) {
	stmt, e := db.Database.Prepare("INSERT into eventDB (date, name, count) VALUES (?, ?, ?)")
	if e != nil {
		return e
	}

	_, e = stmt.Exec(date, name, count)
	if e != nil {
		return e
	}

	return nil
}

func (db EventDB) UpdateEvent(ID, count uint64) (err error) {
	event,e := db.GetEventByID(ID)
	if e != nil {
		return e
	}

	stmt, e := db.Database.Prepare("UPDATE eventDB SET count=? WHERE id=?")
	if e != nil {
		return e
	}

	_, e = stmt.Exec(event.Count + count, ID)
	if e != nil {
		return e
	}

	return nil
}

func (db EventFreqDB) GetEvents() (retrievedEvents []model.EventFreq, err error) {
	rows, e := db.Database.Query("SELECT * FROM eventFreqDB")
	if e != nil {
		return nil, e
	}

	for rows.Next() {
		var (
			event model.EventFreq
			HourCountString string
		)

		e = rows.Scan(&event.ID, &event.Name, &event.TotalCount, &HourCountString)
		if e != nil {
			return nil, e
		}

		e = json.Unmarshal([]byte(HourCountString), &event.HourCount)
		if e != nil {
			return nil, e
		}

		retrievedEvents = append(retrievedEvents, event)
	}

	e = rows.Close()
	if e != nil {
		return nil, e
	}

	return retrievedEvents, nil
}

func (db EventFreqDB) GetEventByID(ID uint64) (retrievedEvent model.EventFreq, err error) {
	rows, e := db.Database.Query(fmt.Sprintf("SELECT * FROM eventFreqDB WHERE id=%d", ID))
	if e != nil {
		return model.EventFreq{}, e
	}

	for rows.Next() {
		var HourCountString string

		e = rows.Scan(&retrievedEvent.ID, &retrievedEvent.Name, &retrievedEvent.TotalCount, &HourCountString)
		if e != nil {
			return model.EventFreq{}, e
		}

		e = json.Unmarshal([]byte(HourCountString), &retrievedEvent.HourCount)
		if e != nil {
			return model.EventFreq{}, e
		}
	}

	e = rows.Close()
	if e != nil {
		return model.EventFreq{}, e
	}

	if retrievedEvent.ID == 0 {
		return model.EventFreq{}, model.ErrEventNotFound
	}
	return retrievedEvent, nil
}

func (db EventFreqDB) GetEventByName(name string) (retrievedEvent model.EventFreq, err error) {
	rows, e := db.Database.Query(fmt.Sprintf("SELECT * FROM eventFreqDB WHERE name LIKE \"%s\"", name))
	if e != nil {
		return model.EventFreq{}, e
	}

	for rows.Next() {
		var HourCountString string

		e = rows.Scan(&retrievedEvent.ID, &retrievedEvent.Name, &retrievedEvent.TotalCount, &HourCountString)
		if e != nil {
			return model.EventFreq{}, e
		}

		e = json.Unmarshal([]byte(HourCountString), &retrievedEvent.HourCount)
		if e != nil {
			return model.EventFreq{}, e
		}
	}

	e = rows.Close()
	if e != nil {
		return model.EventFreq{}, e
	}

	if retrievedEvent.ID == 0 {
		return model.EventFreq{}, model.ErrEventNotFound
	}
	return retrievedEvent, nil
}

func (db EventFreqDB) CreateEvent(name string, count uint64, hour uint64) (err error) {
	var hourCount [24]uint64
	hourCount[hour] = count

	stmt, e := db.Database.Prepare("INSERT into eventFreqDB (name, count, hour_count) VALUES (?, ?, ?)")
	if e != nil {
		return e
	}

	hourCountBytes, e := json.Marshal(hourCount)
	if e != nil {
		return e
	}
	hourCountString := string(hourCountBytes)

	_, e = stmt.Exec(name, count, hourCountString)
	if e != nil {
		return e
	}

	return nil
}

func (db EventFreqDB) UpdateEvent(ID, count, hour uint64) (err error) {
	event, e := db.GetEventByID(ID)
	if e != nil {
		return e
	}

	event.HourCount[hour] += count
	newTotalCount := event.TotalCount + count

	stmt, e := db.Database.Prepare("UPDATE eventFreqDB SET count=?, hour_count=? where ID=?")
	if e != nil {
		return e
	}

	hourCountBytes, e := json.Marshal(event.HourCount)
	if e != nil {
		return e
	}
	hourCountString := string(hourCountBytes)

	_, e = stmt.Exec(newTotalCount, hourCountString, ID)
	if e != nil {
		return e
	}

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

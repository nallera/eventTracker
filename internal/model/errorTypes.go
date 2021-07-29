package model

import "errors"

var (
	ErrEventNotFound    = errors.New("event %s not found")
	ErrStartEventDB     = errors.New("error starting event db")
	ErrStartEventFreqDB = errors.New("error starting event freq db")
)

var (
	ErrInsertEventDB     = errors.New("error inserting new event in event db")
	ErrInsertEventFreqDB = errors.New("error inserting new event in event freq db")
	ErrUpdateEventDB     = errors.New("error updating new event in event db")
	ErrUpdateEventFreqDB = errors.New("error updating new event in event freq db")
	ErrParseHour         = errors.New("error parsing hour into int")
	ErrParseDate         = errors.New("error parsing string into date")
)


package model

import "errors"

var (
	ErrEventNotFound    = errors.New("event %s not found")
	ErrInsertEventDB     = errors.New("error inserting new event in event db: %s")
	ErrInsertEventFreqDB = errors.New("error inserting new event in event freq db: %s")
	ErrUpdateEventDB     = errors.New("error updating new event in event db: %s")
	ErrUpdateEventFreqDB = errors.New("error updating new event in event freq db: %s")
	ErrParseHour         = errors.New("error parsing hour into int")
)


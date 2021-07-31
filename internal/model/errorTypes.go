package model

import "errors"

var (
	ErrEventNotFound          = errors.New("event %s not found")
	ErrInsertEventDB          = errors.New("error inserting new event in event db: %s")
	ErrInsertEventFreqDB      = errors.New("error inserting new event in event freq db: %s")
	ErrUpdateEventDB          = errors.New("error updating new event in event db: %s")
	ErrUpdateEventFreqDB      = errors.New("error updating new event in event freq db: %s")
	ErrDeleteEventDB          = errors.New("error deleting new event in event db: %s")
	ErrDeleteEventFreqDB      = errors.New("error deleting new event in event freq db: %s")
	ErrParseHour              = errors.New("error parsing hour into int")
	ErrDoesntExistEventDB     = errors.New("error trying to delete non existing event %s")
	ErrDoesntExistEventFreqDB = errors.New("error trying to delete non existing event freq %s")
)


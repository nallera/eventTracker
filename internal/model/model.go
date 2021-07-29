package model

type Event struct {
	ID    uint64 `json:"-"`
	Name  string `json:"event"`
	Count uint64 `json:"count"`
	Date  string `json:"date"`
}

type EventFreq struct {
	ID         uint64     `json:"-"`
	Name       string     `json:"event"`
	TotalCount uint64     `json:"count"`
	HourCount  [24]uint64 `json:"hour_count"`
}

type EventBody struct {
	Count uint64 `json:"count,omitempty"`
}


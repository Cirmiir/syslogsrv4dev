package storages

import (
	"time"
)

type Lookup struct {
	Key  interface{}
	Name string
}

type LogEvent struct {
	Host     string
	Message  string
	Severity int
	Facility int
	Date     time.Time
	AppName  string
}

type Criteria struct {
	Page  int
	Limit int
}

type FetchCriteria struct {
	Criteria
	StartDate time.Time
	EndDate   time.Time
	Filters   map[Filter]interface{}
}

type Filter int

const (
	Severity Filter = iota
	Facility
	Application
)

type LogStorage interface {
	SaveEvent(event LogEvent) (bool, error)
	FetchEvents(criteria FetchCriteria, page int, limit int) ([]LogEvent, int)
	FetchFilter(filter Filter) ([]Lookup, error)
}
type MemoryStorage struct {
	Logs []LogEvent
}

func (storage *MemoryStorage) SaveEvent(event LogEvent) (bool, error) {
	storage.Logs = append(storage.Logs, event)
	return true, nil
}

func (filter Filter) GetColumnForSQL() string {
	switch filter {
	case Severity:
		return "severity"
	case Facility:
		return "facility"
	case Application:
		return "appname"
	}
	return ""
}

func (filter Filter) GetTitle() string {
	switch filter {
	case Severity:
		return "Severity"
	case Facility:
		return "Facility"
	case Application:
		return "Application Name"
	}
	return ""
}

func (filter Filter) GetBinding() string {
	switch filter {
	case Severity:
		return "Severity"
	case Facility:
		return "Facility"
	case Application:
		return "AppName"
	}
	return ""
}

func (storage *MemoryStorage) FetchEvents(criteria FetchCriteria, page int, limit int) ([]LogEvent, int) {
	length := len(storage.Logs)
	start := page * limit
	if start > length {
		start = length
	}
	end := (page + 1) * limit
	if end > length {
		end = length
	}
	return storage.Logs[start:end], len(storage.Logs)
}

package memorystorage

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu                sync.RWMutex
	eventsByID        map[string]*storage.Event
	eventsByDateStart map[string]*storage.Event
	eventsByDay       map[int64][]storage.Event
}

var (
	ErrEventDuplicateID = errors.New("duplicate event id in storage")
	ErrEventNotExist    = errors.New("event not found in storage")
	ErrDateBusy         = errors.New("event date start is busy")
)

func (s *Storage) CreateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.NewString()

	if _, ok := s.eventsByID[id]; ok {
		return ErrEventDuplicateID
	}

	e := event

	if _, ok := s.eventsByDateStart[e.DateStart]; ok {
		return ErrDateBusy
	}

	s.createEvent(id, e)

	return nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.eventsByID[id]; !ok {
		return ErrEventNotExist
	}

	e := event

	if _, ok := s.eventsByDateStart[e.DateStart]; ok {
		return ErrDateBusy
	}

	s.deleteEvent(id, e)
	s.createEvent(id, e)

	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.eventsByID[id]
	if !ok {
		return ErrEventNotExist
	}

	s.deleteEvent(id, *e)

	return nil
}

func (s *Storage) ListEventDay(date string) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	unixTime := t.Unix()

	events := make([]storage.Event, 0)

	slice, ok := s.eventsByDay[unixTime]
	if !ok {
		return events, nil
	}

	events = append(events, slice...)

	sort.Slice(events, func(i, j int) bool {
		return events[i].DateStartUnix() < events[j].DateStartUnix()
	})

	return events, nil
}

func (s *Storage) ListEventWeek(date string) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]storage.Event, 0)

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	startTimeUnix := t.Unix()

	endTimeUnix := t.AddDate(0, 0, 7).Unix()

	for unixTime, slice := range s.eventsByDay {
		if unixTime >= startTimeUnix && unixTime < endTimeUnix {
			events = append(events, slice...)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].DateStartUnix() < events[j].DateStartUnix()
	})

	return events, nil
}

func (s *Storage) ListEventMonth(date string) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]storage.Event, 0)

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	startTimeUnix := t.Unix()

	endTimeUnix := t.AddDate(0, 1, 0).Unix()

	for unixTime, slice := range s.eventsByDay {
		if unixTime >= startTimeUnix && unixTime < endTimeUnix {
			events = append(events, slice...)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].DateStartUnix() < events[j].DateStartUnix()
	})

	return events, nil
}

func (s *Storage) Event(id string) (storage.Event, error) {
	val, ok := s.eventsByID[id]
	if !ok {
		return storage.Event{}, ErrEventNotExist
	}

	return *val, nil
}

func (s *Storage) createEvent(id string, e storage.Event) {
	unixTime := e.DateStartDayUnix()
	s.eventsByDay[unixTime] = append(s.eventsByDay[unixTime], e)

	s.eventsByID[id] = &e
	s.eventsByDateStart[e.DateStart] = &e
}

func (s *Storage) deleteEvent(id string, e storage.Event) {
	delete(s.eventsByID, id)
	delete(s.eventsByDateStart, e.DateStart)

	dayTimeUnix := e.DateStartDayUnix()

	slice := s.eventsByDay[dayTimeUnix]

	for k, v := range slice {
		if v.ID != id {
			continue
		}
		s.eventsByDay[dayTimeUnix] = append(slice[:k], slice[k+1:]...)
		break
	}
}

func New() *Storage {
	return &Storage{
		eventsByID:        make(map[string]*storage.Event),
		eventsByDateStart: make(map[string]*storage.Event),
		eventsByDay:       make(map[int64][]storage.Event),
	}
}

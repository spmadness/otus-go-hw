package internalhttp

import (
	"encoding/json"
	"net/http"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/server"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type UpdateRequest struct {
	ID    string        `json:"id"`
	Event storage.Event `json:"event"`
}

type ListHandlerFunc func(date string) ([]storage.Event, error)

func createEvent(w http.ResponseWriter, r *http.Request, s app.Storager) (interface{}, error) {
	e := storage.Event{}

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err
	}

	err := server.ValidateCreateEvent(e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = s.CreateEvent(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func updateEvent(w http.ResponseWriter, r *http.Request, s app.Storager) (interface{}, error) {
	ur := UpdateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		return nil, err
	}

	id := ur.ID
	update := ur.Event

	e, err := s.GetEvent(id)
	if err != nil {
		return nil, err
	}

	updateEventFields(&e, &update)

	err = server.ValidateUpdateEvent(e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = s.UpdateEvent(e.ID, e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func deleteEvent(w http.ResponseWriter, r *http.Request, s app.Storager) (interface{}, error) {
	e := storage.Event{}

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err
	}

	err := server.ValidateDeleteEvent(e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = s.DeleteEvent(e.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func listEventDay(w http.ResponseWriter, r *http.Request, s app.Storager) (interface{}, error) {
	events, err := listEvent(w, r, s.ListEventDay)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return events, nil
}

func listEventWeek(w http.ResponseWriter, r *http.Request, s app.Storager) (interface{}, error) {
	events, err := listEvent(w, r, s.ListEventWeek)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return events, nil
}

func listEventMonth(w http.ResponseWriter, r *http.Request, s app.Storager) (interface{}, error) {
	events, err := listEvent(w, r, s.ListEventMonth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return events, nil
}

func listEvent(w http.ResponseWriter, r *http.Request, f ListHandlerFunc) ([]storage.Event, error) {
	lm := storage.ListEventValidation{}

	if err := json.NewDecoder(r.Body).Decode(&lm); err != nil {
		return nil, err
	}

	err := server.ValidateListEvent(lm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	return f(lm.DateStart)
}

func updateEventFields(e *storage.Event, changed *storage.Event) {
	if changed.Title != "" && changed.Title != e.Title {
		e.Title = changed.Title
	}
	if changed.DateStart != "" && changed.DateStart != e.DateStart {
		e.DateStart = changed.DateStart
	}
	if changed.DateEnd != "" && changed.DateEnd != e.DateEnd {
		e.DateEnd = changed.DateEnd
	}
	if changed.Description != "" && changed.Description != e.Description {
		e.Description = changed.Description
	}
	if changed.UserID != "" && changed.UserID != e.UserID {
		e.UserID = changed.UserID
	}
	if changed.DatePost != "" && changed.DatePost != e.DatePost {
		e.DatePost = changed.DatePost
	}
}

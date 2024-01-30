package storage

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEventDuplicateID = errors.New("duplicate event id in storage")
	ErrEventNotExist    = errors.New("event not found in storage")
	ErrDateBusy         = errors.New("event date start is busy")
)

type Event struct {
	ID          string `json:"id" validate:"required,uuid"`
	Title       string `json:"title" validate:"required"`
	DateStart   string `json:"dateStart" validate:"required,datetime=2006-01-02 15:04:05"`
	DateEnd     string `json:"dateEnd" validate:"required,datetime=2006-01-02 15:04:05" `
	Description string `json:"description"`
	UserID      string `json:"userId" validate:"required,uuid"`
	DatePost    string `json:"datePost" validate:"datetime=2006-01-02 15:04:05"`
}

type ListEventValidation struct {
	DateStart string `json:"dateStart" validate:"required,datetime=2006-01-02"`
}

func (e *Event) DateStartUnix() int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", e.DateStart)
	return t.Unix()
}

func (e *Event) DateStartDayUnix() int64 {
	s := strings.Split(e.DateStart, " ")
	t, _ := time.Parse("2006-01-02", s[0])
	return t.Unix()
}

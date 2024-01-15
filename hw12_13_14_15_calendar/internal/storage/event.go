package storage

import (
	"strings"
	"time"
)

type Event struct {
	ID          string
	Title       string
	DateStart   string
	DateEnd     string
	Description string
	UserID      int
	DatePost    string
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

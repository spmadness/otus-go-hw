package app

import (
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	DBModeSQL      = "sql"
	DBModeInMemory = "in-memory"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface{}

type Storage interface {
	CreateEvent(event storage.Event) error
	UpdateEvent(id string, event storage.Event) error
	DeleteEvent(id string) error
	ListEventDay(date string) ([]storage.Event, error)
	ListEventWeek(date string) ([]storage.Event, error)
	ListEventMonth(date string) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func NewStorage(mode string, dsn string) Storage {
	var s Storage
	if mode == DBModeSQL {
		s = sqlstorage.New(dsn)
	}
	if mode == DBModeInMemory {
		s = memorystorage.New()
	}
	return s
}

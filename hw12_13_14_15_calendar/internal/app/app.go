package app

import (
	"context"

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
	storage Storager
}

type Logger interface{}

type StorageEvent interface {
	CreateEvent(event storage.Event) error
	UpdateEvent(id string, event storage.Event) error
	DeleteEvent(id string) error
	GetEvent(id string) (storage.Event, error)
	ListEventDay(date string) ([]storage.Event, error)
	ListEventWeek(date string) ([]storage.Event, error)
	ListEventMonth(date string) ([]storage.Event, error)
}

type StorageConnector interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
}

type Storager interface {
	StorageEvent
	StorageConnector
}

func (a *App) GetStorage() Storager {
	return a.storage
}

func New(logger Logger, storage Storager) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func NewStorage(mode string, dsn string) Storager {
	var s Storager
	if mode == DBModeSQL {
		s = sqlstorage.New(dsn)
	}
	if mode == DBModeInMemory {
		s = memorystorage.New()
	}
	return s
}

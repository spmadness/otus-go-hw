package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // postgres driver
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	dsn  string
	Conn *sql.DB
}

const QueryTimeout = time.Second * 3

const selectFields = "select id, title, date_start, date_end, description, user_id, date_post"

func (s *Storage) CreateEvent(event storage.Event) error {
	query := "insert into events (title, date_start, date_end, description, user_id, date_post) " +
		"values ($1, $2 ,$3 ,$4 ,$5 ,$6)"

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	_, err := s.Conn.ExecContext(ctx,
		query, event.Title, event.DateStart, event.DateEnd, event.Description, event.UserID, event.DatePost)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) error {
	query := "update events " +
		"set title = $2, date_start = $3, date_end = $4, description = $5, user_id = $6, date_post = $7 where id = $1"

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	result, err := s.Conn.ExecContext(ctx,
		query, id, event.Title, event.DateStart, event.DateEnd, event.Description, event.UserID, event.DatePost)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("event with id: %s was not found in DB", id)
	}

	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	query := "delete from events where id = $1"

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	result, err := s.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("event with id: %s was not found in DB", id)
	}
	return nil
}

func (s *Storage) ListEventDay(date string) ([]storage.Event, error) {
	var events []storage.Event

	query := selectFields + " from events where DATE(date_start) = DATE($1)"

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	rows, err := s.Conn.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		err = rows.Scan(&e.ID, &e.Title, &e.DateStart, &e.DateEnd, &e.Description, &e.UserID, &e.DatePost)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (s *Storage) ListEventWeek(date string) ([]storage.Event, error) {
	var events []storage.Event

	query := selectFields + " from events " +
		"where DATE(date_start) >= $1 and DATE(date_start) < DATE($1) + INTERVAL '7 DAY'"

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	rows, err := s.Conn.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		err = rows.Scan(&e.ID, &e.Title, &e.DateStart, &e.DateEnd, &e.Description, &e.UserID, &e.DatePost)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (s *Storage) ListEventMonth(date string) ([]storage.Event, error) {
	var events []storage.Event

	query := selectFields + " from events " +
		"where DATE(date_start) >= $1 and DATE(date_start) < DATE($1) + INTERVAL '1 MONTH'"

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	rows, err := s.Conn.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		err = rows.Scan(&e.ID, &e.Title, &e.DateStart, &e.DateEnd, &e.Description, &e.UserID, &e.DatePost)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	var e storage.Event

	query := selectFields + " from events where id = $1"

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	row := s.Conn.QueryRowContext(ctx, query, id)
	err := row.Scan(&e.ID, &e.Title, &e.DateStart, &e.DateEnd, &e.Description, &e.UserID, &e.DatePost)
	if err == sql.ErrNoRows {
		return e, storage.ErrEventNotExist
	}
	if err != nil {
		return e, err
	}

	return e, nil
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Open(ctx context.Context) error {
	var err error
	s.Conn, err = sql.Open("pgx", s.dsn)
	if err != nil {
		return err
	}
	err = s.Conn.PingContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	<-ctx.Done()

	err := s.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

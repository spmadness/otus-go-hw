package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Scheduler struct {
	storage StorageScheduler
	broker  Broker
	logger  Logger
}
type Notification struct {
	ID        string
	Title     string
	DateStart string
	UserID    string
}

type Broker interface {
	Open() error
	Close() error
	SetQueue(queueName string) error
	SendMessage(ctx context.Context, n Notification) error
	ConsumeMessage(queueName string) (<-chan amqp.Delivery, error)
}

func (s *Scheduler) ProcessNotifications(ctx context.Context, pollTime int, outdatedEventDays int) {
	ticker := time.NewTicker(time.Duration(pollTime) * time.Second)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			err := s.sendNotifications()
			if err != nil {
				s.logger.Error(err.Error())
			}

			err = s.removeOutdatedEvents(outdatedEventDays)
			if err != nil {
				s.logger.Error(err.Error())
			}
		}
	}
}

func (s *Scheduler) sendNotifications() error {
	events, err := s.storage.ListEventWithNotification()
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	s.logger.Info(fmt.Sprintf("sending notifications: %d", len(events)))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, event := range events {
		n := Notification{
			ID:        event.ID,
			Title:     event.Title,
			DateStart: event.DateStart,
			UserID:    uuid.NewString(),
		}

		err = s.broker.SendMessage(ctx, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) removeOutdatedEvents(days int) error {
	currentTime := time.Now()
	date := currentTime.AddDate(0, 0, -days).Format("2006-01-02 15:04:05")

	err := s.storage.DeleteEventsBeforeDate(date)
	if err != nil {
		return err
	}

	return nil
}

func NewScheduler(storage StorageScheduler, broker Broker, logger Logger) *Scheduler {
	return &Scheduler{
		storage: storage,
		broker:  broker,
		logger:  logger,
	}
}

package broker

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/app"
)

type Broker struct {
	dsn        string
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func New(dsn string) *Broker {
	return &Broker{
		dsn: dsn,
	}
}

func (b *Broker) Open() error {
	var err error
	b.connection, err = amqp.Dial(b.dsn)

	if err != nil {
		return err
	}

	b.channel, err = b.connection.Channel()
	if err != nil {
		return err
	}
	return nil
}

func (b *Broker) Close() error {
	return b.connection.Close()
}

func (b *Broker) SetQueue(queueName string) error {
	var err error

	b.queue, err = b.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (b *Broker) SendMessage(ctx context.Context, n app.Notification) error {
	jData, err := json.Marshal(&n)
	if err != nil {
		return errors.Errorf("JSON Marshal error: %v", err)
	}

	err = b.channel.PublishWithContext(
		ctx,
		"",
		b.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jData,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (b *Broker) ConsumeMessage(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := b.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/queue"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

const QUEUE = "arpro"

var _ queue.Queue = &Client{}

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func New() (*Client, error) {
	cfg := config.Get()
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/",
			cfg.RabbitMQ.Username,
			cfg.RabbitMQ.Password,
			cfg.RabbitMQ.Host,
			cfg.RabbitMQ.Port,
		),
	)

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	q, err := ch.QueueDeclare(QUEUE, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, ch: ch, q: q}, nil
}

func (c *Client) Publish(msgBody []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.ch.PublishWithContext(
		ctx,
		"",
		c.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgBody,
		},
	)

	if err != nil {
		return errors.Wrap(err, "failed from Publish(msgBody []byte)")
	}

	return nil
}

func (c *Client) Consume() (<-chan amqp.Delivery, error) {
	err := c.ch.Qos(
		2,
		0,
		false,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed from Consume() (<-chan amqp.Delivery, error)")
	}

	msgs, err := c.ch.Consume(
		c.q.Name,
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

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

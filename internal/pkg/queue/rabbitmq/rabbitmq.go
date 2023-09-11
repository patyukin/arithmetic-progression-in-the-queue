package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/queue"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

const QUEUE = "arpro"

var _ queue.Queue = &Client{}

type Client struct {
	conn *amqp.Connection
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

	return &Client{conn: conn}, nil
}

func (c *Client) Publish(msgBody []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch, err := c.conn.Channel()

	q, err := ch.QueueDeclare(QUEUE, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgBody,
		})

	return nil
}

func (c *Client) Consume() (<-chan amqp.Delivery, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		QUEUE,
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

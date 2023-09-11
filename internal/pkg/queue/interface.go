package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Close() error
	Publish(msgBody []byte) error
	Consume() (<-chan amqp.Delivery, error)
}

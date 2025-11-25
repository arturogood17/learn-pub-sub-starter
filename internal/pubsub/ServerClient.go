package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Transient SimpleQueueType = iota
	Durable
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	byteVal, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        byteVal,
	}); err != nil {
		return err
	}
	return nil
}

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	var b bool
	if queueType == Durable {
		b = true
	}

	q, err := ch.QueueDeclare(queueName, b, !b, !b, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	if err = ch.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, q, nil

}

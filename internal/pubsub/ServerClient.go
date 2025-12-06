package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int
type AckType int

const (
	Transient SimpleQueueType = iota
	Durable
)

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
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

	q, err := ch.QueueDeclare(queueName, queueType != Transient,
		queueType == Transient, queueType == Transient, false, amqp.Table{"x-dead-letter-exchange": "peril_dlx"})
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	if err = ch.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, q, nil

}

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	c, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() error {
		for delivery := range c {
			var data T
			if err := json.Unmarshal(delivery.Body, &data); err != nil {
				return err
			}
			ack := handler(data)
			switch ack {
			case Ack:
				if err := delivery.Ack(false); err != nil {
					return err
				}
				log.Println("Action: ack")
			case NackRequeue:
				if err := delivery.Nack(false, true); err != nil {
					return err
				}
				log.Println("Action: NackRequeue")
			case NackDiscard:
				if err := delivery.Nack(false, false); err != nil {
					return err
				}
				log.Println("Action: NackDiscard")
			}
		}
		return nil
	}()

	return nil
}

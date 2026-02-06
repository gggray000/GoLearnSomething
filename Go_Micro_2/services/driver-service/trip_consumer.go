package main

import (
	"context"
	"log"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type TripConsumer struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ) *TripConsumer {
	return &TripConsumer{
		rabbitmq: rabbitmq,
	}
}

func (t *TripConsumer) Listen() error {
	return t.rabbitmq.ConsumeMessages("hello", func(ctx context.Context, msg amqp091.Delivery) error {
		log.Printf("Driver received message: %v", msg)
		return nil
	})
}

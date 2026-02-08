package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type TripConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  *Service
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ, service *Service) *TripConsumer {
	return &TripConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (t *TripConsumer) Listen() error {
	return t.rabbitmq.ConsumeMessages(
		messaging.FindAvailableDriversQueue,
		func(ctx context.Context, msg amqp091.Delivery) error {
			var tripEvent contracts.AmqpMessage
			if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return err
			}

			var payload messaging.TripEventData
			if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return err
			}
			log.Printf("Driver received message: %+v", payload)

			switch msg.RoutingKey {
			case contracts.TripEventCreated, contracts.TripEventDriverNotInterested:
				return t.handleFindAndNotifyDrivers(ctx, payload)
			}

			log.Printf("Unknown trip event: %+v", payload)

			return nil
		})
}

func (t *TripConsumer) handleFindAndNotifyDrivers(ctx context.Context, payload messaging.TripEventData) error {
	suitableDriverIDs := t.service.FindAvailableDrivers(payload.Trip.SelectedFare.PackageSlug)

	log.Printf("Found suitable %v driver ", len(suitableDriverIDs))
	if len(suitableDriverIDs) == 0 {
		// Notify the driver that no drivers are available
		if err := t.rabbitmq.PublishMessage(
			ctx,
			contracts.TripEventNoDriversFound,
			contracts.AmqpMessage{
				OwnerID: payload.Trip.UserID,
			},
		); err != nil {
			log.Printf("Failed to publish message to exchange %v", err)
			return err
		}
		return nil
	} else {
		// Notify driver for the potential trip
		randomIndex := rand.Intn(len(suitableDriverIDs))
		suitableDriverID := suitableDriverIDs[randomIndex]
		marshalledEvent, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if err := t.rabbitmq.PublishMessage(
			ctx,
			contracts.DriverCmdTripRequest,
			contracts.AmqpMessage{
				OwnerID: suitableDriverID,
				Data:    marshalledEvent,
			},
		); err != nil {
			log.Printf("Failed to publish message to exchange %v", err)
			return err
		}
	}
	return nil
}

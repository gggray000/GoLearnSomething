package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	pb_d "ride-sharing/shared/proto/driver"

	"github.com/rabbitmq/amqp091-go"
)

type DriverConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  domain.TripService
}

func NewDriverConsumer(rabbitmq *messaging.RabbitMQ, service domain.TripService) *DriverConsumer {
	return &DriverConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (d *DriverConsumer) Listen() error {
	log.Printf("Trip-service consuming queue=%q", messaging.DriverTripResponseQueue)
	log.Printf("Driver request queue=%q", messaging.DriverCmdTripRequestQueue)
	return d.rabbitmq.ConsumeMessages(
		messaging.DriverTripResponseQueue,
		func(ctx context.Context, msg amqp091.Delivery) error {
			var message contracts.AmqpMessage
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return err
			}

			var payload messaging.DriverTripResponse
			if err := json.Unmarshal(message.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return err
			}
			log.Printf("Driver response received message: %+v", payload)
			log.Printf("routing=%s raw=%s", msg.RoutingKey, string(msg.Body))

			switch msg.RoutingKey {
			case contracts.DriverCmdTripAccept:
				if err := d.handleTripAccepted(ctx, payload.TripID, payload.Driver); err != nil {
					log.Printf("Failed to handle trip accept: %v", err)
					return err
				}
			case contracts.DriverCmdTripDecline:
				if err := d.handleTripDeclined(ctx, payload.TripID, payload.RiderID); err != nil {
					log.Printf("Failed to handle trip decline: %v", err)
					return err
				}
				return nil
			}

			log.Printf("Unknown trip event: %+v", payload)

			return nil
		})
}

func (d *DriverConsumer) handleTripAccepted(ctx context.Context, tripID string, driver *pb_d.Driver) error {
	trip, err := d.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	if trip == nil {
		return fmt.Errorf("Trip was not found %s", tripID)
	}

	if err := d.service.UpdateTrip(ctx, tripID, "accepted", driver); err != nil {
		log.Printf("Failed to update the trip: %v", err)
		return err
	}

	trip, err = d.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	log.Printf("Driver Info: %+v", trip.Driver)

	marshalledTrip, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	// Notify rider that a driver has been assigned
	if err := d.rabbitmq.PublishMessage(
		ctx,
		contracts.TripEventDriverAssigned,
		contracts.AmqpMessage{
			OwnerID: trip.UserID,
			Data:    marshalledTrip,
		}); err != nil {
		return err
	}

	marshalledPayload, err := json.Marshal(messaging.PaymentTripRequestData{
		TripID:   tripID,
		UserID:   trip.UserID,
		DriverID: driver.Id,
		Amount:   trip.RideFare.TotalPriceInCents,
		Currency: "EUR",
	})

	if err := d.rabbitmq.PublishMessage(ctx, contracts.PaymentCmdCreateSession,
		contracts.AmqpMessage{
			OwnerID: trip.UserID,
			Data:    marshalledPayload,
		},
	); err != nil {
		return err
	}

	return nil
}

func (d *DriverConsumer) handleTripDeclined(ctx context.Context, tripID string, riderID string) error {
	trip, err := d.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	newPayload := messaging.TripEventData{
		Trip: trip.ToProto(),
	}

	marshalledPayload, err := json.Marshal(newPayload)
	if err != nil {
		return err
	}

	if err := d.rabbitmq.PublishMessage(
		ctx,
		contracts.TripEventDriverNotInterested,
		contracts.AmqpMessage{
			OwnerID: riderID,
			Data:    marshalledPayload,
		}); err != nil {
		return err
	}
	return nil
}

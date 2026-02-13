package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/retry"
	"ride-sharing/shared/tracing"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TripExchange = "trip"
	DeadLetterExchange = "dlx"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQConnection(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
	}

	log.Println("Starting RabbitMQ connection")

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Failed to create channel: %v", err)
	}

	rmq := &RabbitMQ{
		conn:    conn,
		Channel: ch,
	}

	if err := rmq.setupExchangesAndQueues(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("Failed to setup exchanges and queues: %v", err)
	}

	return rmq, nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
	if r.Channel != nil {
		r.Channel.Close()
	}
}

func (r *RabbitMQ) setupExchangesAndQueues() error {

	if err := r.setupDeadLetterExchange(); err != nil {
		return err
	}

	err := r.Channel.ExchangeDeclare(
		TripExchange, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("Failed to declare exchange: %s: %v", TripExchange, err)
	}

	if err := r.declareAndBindQueue(
		FindAvailableDriversQueue,
		[]string{
			contracts.TripEventCreated,
			contracts.TripEventDriverNotInterested,
		},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		DriverCmdTripRequestQueue,
		[]string{
			contracts.DriverCmdTripRequest,
		},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		DriverTripResponseQueue,
		[]string{
			contracts.DriverCmdTripAccept,
			contracts.DriverCmdTripDecline,
		},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		NotifyDriverNoDriversFoundQueue,
		[]string{
			contracts.TripEventNoDriversFound,
		},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		NotifyDriverAssignedQueue,
		[]string{
			contracts.TripEventDriverAssigned,
		},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		PaymentTripRequestQueue,
		[]string{
			contracts.PaymentCmdCreateSession,
		},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		NotifyPaymentSessionCreatedQueue,
		[]string{
			contracts.PaymentEventSessionCreated,
		},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		NotifyPaymentSuccessQueue,
		[]string{
			contracts.PaymentEventSuccess,
		},
		TripExchange,
	); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) setupDeadLetterExchange() error {
	err := r.Channel.ExchangeDeclare(
		DeadLetterExchange, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("Failed to declare exchange: %s: %v", TripExchange, err)
	}

	q, err := r.Channel.QueueDeclare(
		DeadLetterQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := r.Channel.QueueBind(
			q.Name,   // queue name
			"#",      // wildcard routing key to catch all messages
			DeadLetterExchange, // exchange
			false,
			nil,
		); err != nil {
			return fmt.Errorf("Failed to bind queue to %s: %v", q.Name, err)
		}
	return nil
}

func (r *RabbitMQ) declareAndBindQueue(queueName string, messageTypes []string, exchange string) error {

	// Add DLQ config
	args := amqp.Table{
		"x-dead-letter-exchange": DeadLetterExchange,
	}

	q, err := r.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range messageTypes {
		if err := r.Channel.QueueBind(
			q.Name,   // queue name
			msg,      // routing key
			exchange, // exchange
			false,
			nil,
		); err != nil {
			return fmt.Errorf("Failed to bind queue to %s: %v", q.Name, err)
		}
	}
	return nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {
	log.Printf("Publishing message with routing key: %s", routingKey)

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("Failed to marshal message: %v", err)
	}

	msg := amqp.Publishing{
			ContentType:  "application/json", // "text/plain" before video No.85
			Body:         jsonMsg,
			DeliveryMode: amqp.Persistent,
		}
	
	return tracing.TracedPublisher(ctx, TripExchange, routingKey, msg, r.publishWithTracing)
}

func (r *RabbitMQ) publishWithTracing(ctx context.Context, exchange, routingKey string, msg amqp.Publishing) error {
	return r.Channel.PublishWithContext(
			ctx,
			exchange, 	  // exchange
			routingKey,   // routing key
			false,        // mandatory
			false,        // immediate
			msg,
		)
}

type MessageHandler func(context.Context, amqp.Delivery) error

func (r *RabbitMQ) ConsumeMessages(queueName string, handler MessageHandler) error {
	// Fair dispatch
	err := r.Channel.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("Failed to set QoS: %v", err)
	}

	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := tracing.TracedConsumer(msg, func(ctx context.Context, d amqp.Delivery) error {
				log.Printf("Received a message: %s", msg.Body)

				cfg := retry.DefaultConfig()
				err := retry.WithBackoff(ctx, cfg, func() error{
					return handler(ctx, d)
				})
				if err != nil {
					log.Printf("Message processing failed after %d retries for message ID: %s, err: %v",
					cfg.MaxRetries, d.MessageId ,err )
					
					headers := amqp.Table{}
					if d.Headers != nil {
						headers = d.Headers
					}

					headers["x-death-reason"] = err.Error()
					headers["x-origin-exchange"]= d.Exchange
					headers["x-origin-routing-key"] = d.RoutingKey
					headers["x-retry-count"] = cfg.MaxRetries
					d.Headers = headers

					_ = d.Reject(false)
					return err
				}

				if ackErr := msg.Ack(false); ackErr != nil {
					log.Printf("Error: Failed to Ack message: %v. Message body: %s", ackErr, msg.Body)
				}
				return nil
			}); err != nil {
				log.Printf("Failed to consume message: %v", err)
			}
		}
	}()

	return nil
}

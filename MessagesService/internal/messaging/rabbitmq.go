package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQConn struct {
	conn *amqp.Connection
}

func NewRabbitMQ(url string) (*RabbitMQConn, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	log.Println("Connected Successfully to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	// declare internal queue for message_creation_requested event
	_, err = ch.QueueDeclare(
		"message_creation_requested", // name
		true,                         // durable
		false,                        // delete when unused
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQConn{conn: conn}, nil
}

func (r *RabbitMQConn) PublishEvent(event string, eventData interface{}) error {
	body, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		event,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish event %s: %v", event, err)
	} else {
		log.Println("Published event", event)
	}
	return err
}

func (r *RabbitMQConn) ConsumeMessages(ctx context.Context, queueName string, handler func(amqp.Delivery)) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if ok {
					handler(d)
				} else {
					return
				}
			case <-ctx.Done():
				log.Printf("Stopping consumer for queue: %s", queueName)
				return
			}
		}

	}()

	log.Printf("Waiting for messages from queue: %s", queueName)
	<-ctx.Done()

	return nil
}

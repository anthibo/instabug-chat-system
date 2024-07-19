package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"message_service/internal/events"

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

	for _, eventQueue := range events.EventQueues {
		_, err := ch.QueueDeclare(
			eventQueue.Name, // name
			true,            // durable
			false,           // delete when unused
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	return &RabbitMQConn{conn: conn}, nil
}

func (r *RabbitMQConn) PublishEvent(event events.EventQueue, eventData interface{}) error {
	body, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	ch, err := r.conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel")
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		"",
		event.Name,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish event %s: %v", event, err)
	} else {
		log.Println("Published event", event, "with event data", eventData)
	}
	return err
}

func (r *RabbitMQConn) ConsumeMessages(ctx context.Context, queueName string, handler func(amqp.Delivery)) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		fmt.Println("Error consuming messages: ", err)
		ch.Close()
		return err
	}

	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if ok {
					handler(d)
				} else {
					log.Printf("Consumer for queue: %s closed", queueName)
					return
				}
			case <-ctx.Done():
				log.Printf("Stopping consumer for queue: %s", queueName)
				ch.Close()
				return
			}
		}
	}()

	log.Printf("Waiting for messages from queue: %s", queueName)
	return nil
}

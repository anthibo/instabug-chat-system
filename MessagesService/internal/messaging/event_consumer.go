package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	QueueName string
	Handler   func(amqp.Delivery)
}

type ConsumerManager struct {
	consumer  *RabbitMQConn
	consumers map[string]func(amqp.Delivery)
}

func NewConsumerManager(consumer *RabbitMQConn) *ConsumerManager {
	return &ConsumerManager{
		consumer:  consumer,
		consumers: make(map[string]func(amqp.Delivery)),
	}
}

func (cm *ConsumerManager) AddConsumer(queueName string, handler func(amqp.Delivery)) {
	cm.consumers[queueName] = handler
}

func (cm *ConsumerManager) StartConsumers(ctx context.Context) error {
	fmt.Println("Starting event consumers...")
	for queueName, handler := range cm.consumers {
		go func(queue string, h func(amqp.Delivery)) {
			fmt.Println("Starting consumer for queue: ", queue)
			err := cm.consumer.ConsumeMessages(ctx, queue, h)
			if err != nil {
				log.Fatalf("Failed to start consumer for queue %s: %v", queue, err)
			}
		}(queueName, handler)
	}
	return nil
}

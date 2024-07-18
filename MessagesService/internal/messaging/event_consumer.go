package messaging

import (
	"context"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	QueueName string
	Handler   func(amqp.Delivery)
}

type ConsumerManager struct {
	consumer  *RabbitMQConn
	consumers []Consumer
}

func NewConsumerManager(consumer *RabbitMQConn) *ConsumerManager {
	return &ConsumerManager{
		consumer:  consumer,
		consumers: []Consumer{},
	}
}

func (cm *ConsumerManager) AddConsumer(queueName string, handler func(amqp.Delivery)) {
	consumer := Consumer{
		QueueName: queueName,
		Handler:   handler,
	}
	cm.consumers = append(cm.consumers, consumer)
}

func (cm *ConsumerManager) StartConsumers(ctx context.Context) error {
	for _, consumer := range cm.consumers {
		go func(consumer Consumer) {
			err := cm.consumer.ConsumeMessages(ctx, consumer.QueueName, consumer.Handler)
			if err != nil {
				log.Printf("Failed to start consumer for queue %s: %v", consumer.QueueName, err)
			}
		}(consumer)
	}
	return nil
}

package handlers

import (
	"context"
	"encoding/json"
	"log"
	"message_service/internal/models"
	"message_service/internal/services"

	"github.com/streadway/amqp"
)

type EventHandler struct {
	MessageService *services.MessageService
}

type IEventHHandlers interface {
	HandleMessageCreated(d amqp.Delivery)
}

func NewEventHandler(messageService *services.MessageService) *EventHandler {
	return &EventHandler{
		MessageService: messageService,
	}
}

func (eventHandler *EventHandler) HandleMessageCreated(d amqp.Delivery) {
	messageCreationRequestedEvent := models.MessageCreationRequestedEvent{}
	err := json.Unmarshal(d.Body, &messageCreationRequestedEvent)
	log.Printf("\n Received message created event: %v", messageCreationRequestedEvent)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		d.Nack(false, true)
		return
	}

	ctx := context.Background()

	err = eventHandler.MessageService.CreateMessage(ctx, &messageCreationRequestedEvent)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		d.Nack(false, true)
		return
	}

	d.Ack(false)
}

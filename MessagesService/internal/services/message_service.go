package services

import (
	"context"
	"fmt"
	"log"

	"message_service/internal/events"
	"message_service/internal/messaging"
	"message_service/internal/models"
	"message_service/internal/repositories"
)

type MessageService struct {
	repositories.MessageRepository
	repositories.ChatRepository
	messaging.EventPublisherManager
}

func (s *MessageService) CreateMessage(ctx context.Context, messageCreationRequestedEvent *models.MessageCreationRequestedEvent) error {
	message := &models.Message{
		ChatID: messageCreationRequestedEvent.ChatID,
		Body:   messageCreationRequestedEvent.Body,
		Number: messageCreationRequestedEvent.MessageNumber,
	}
	if err := s.MessageRepository.CreateMessage(ctx, message); err != nil {
		return err
	}
	log.Printf("Created message with ID %d for chat %d", message.ID, messageCreationRequestedEvent.ChatID)

	messageCreatedEventData := &models.MessageCreatedEvent{
		MessageID: message.ID,
		Number:    message.Number,
		ChatID:    message.ChatID,
		Body:      message.Body,
	}

	// TODO: add retry logic here for more resiliency
	log.Println("Publishing Message Created Event")
	messageCreatedEventQueue, exists := events.GetEventQueue(events.MessageCreatedQueue)
	if !exists {
		fmt.Println("event queue not found")
		return nil
	}
	if err := s.EventPublisherManager.PublishEvent(messageCreatedEventQueue, messageCreatedEventData); err != nil {
		log.Printf("Failed to publish event %s: %v", messageCreatedEventQueue.Name, err)
		return err
	}

	return nil
}

func (s *MessageService) GetChatId(ctx context.Context, applicationToken string, chatNumber int) (int, error) {
	chat, err := s.ChatRepository.GetChatByChatNumber(ctx, applicationToken, chatNumber)
	if err != nil {
		log.Printf("Failed to get chat ID from chat number %d: %v", chatNumber, err)
		return -1, err
	}
	return chat.ID, nil
}

package services

import (
	"context"
	"log"
	"strconv"

	"message_service/internal/messaging"
	"message_service/internal/models"
	"message_service/internal/repositories"
)

type MessageService struct {
	repositories.MessageRepository
	repositories.ChatRepository
	messaging.EventPublisherManager
}

func (s *MessageService) CreateMessage(ctx context.Context, message *models.Message, applicationToken string, chatNumber int) error {
	chat, err := s.ChatRepository.GetChatByChatNumber(ctx, applicationToken, chatNumber)
	if err != nil {
		log.Printf("Failed to get chat ID from chat number %d: %v", chatNumber, err)
		return err
	}
	message.ChatID = chat.ID

	if err := s.MessageRepository.CreateMessage(ctx, message); err != nil {
		return err
	}
	log.Printf("Created message with ID %d for chat %d", message.ID, message.ChatID)

	// TODO: Implement Better Event Handling Approach
	event := "message_created"
	eventData := []byte(`{"id":` + strconv.Itoa(message.ID) + `,"chat_id":` + strconv.Itoa(message.ChatID) + `,"body":"` + message.Body + `","created_at":"` + message.CreatedAt + `"}`)

	// TODO: add retry logic here for more resiliency
	log.Println("Publishing Event Created Event")
	if err := s.EventPublisherManager.PublishEvent(event, eventData); err != nil {
		log.Printf("Failed to publish event %s: %v", event, err)
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

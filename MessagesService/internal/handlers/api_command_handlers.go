package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"message_service/internal/caching"
	"message_service/internal/events"
	"message_service/internal/helpers"
	"message_service/internal/messaging"
	"message_service/internal/models"
	"message_service/internal/services"
)

// For the C part (Commands) in CQRS Pattern, we have the Command Handlers.
// As for the Query Handlers, they are implemented in the RoR API.
type ApiCmdHandlers struct {
	MessageService        *services.MessageService
	Cache                 *caching.RedisCacheManager
	EventPublisherManager *messaging.RabbitMQConn
}

func (cmdHandler *ApiCmdHandlers) CreateMessageCmdHandler(w http.ResponseWriter, r *http.Request) {
	var message models.CreateMessageBody
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Print("create message body: ", message)
	if message.Body == "" {
		helpers.ErrorJSON(w, errors.New("message body is required"), http.StatusBadRequest)
		return
	}

	applicationToken, err := helpers.GetParameterFromURLPath(r, "token")
	if err != nil {
		helpers.ErrorJSON(w, errors.New("application token not set properly"), http.StatusBadRequest)
		return
	}
	chatNumberStr, err := helpers.GetParameterFromURLPath(r, "number")
	if err != nil {
		helpers.ErrorJSON(w, errors.New("chat number not set properly"), http.StatusBadRequest)
		return
	}
	chatNumber, _ := strconv.Atoi(chatNumberStr)

	log.Println("Chat Number: ", chatNumber)

	chatId, err := cmdHandler.MessageService.GetChatId(r.Context(), applicationToken, chatNumber)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if chatId == -1 {
		helpers.ErrorJSON(w, errors.New("chat not found"), http.StatusNotFound)
		return
	}
	fmt.Println("Chat ID: ", chatId)

	// cache key: chat-<chat_id>-latestMessageNo
	cacheKey := fmt.Sprintf("chat-%d-latestMessageNo", chatId)
	latestMessageNumber, err := cmdHandler.Cache.Get(cacheKey)
	if err != nil {
		fmt.Println("error getting latest message number from cache: ", err)
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Println("Latest Message Number: ", latestMessageNumber)

	var latestMessageNumberInt int
	if latestMessageNumber == nil {
		latestMessageNumberInt = 0
	} else {
		latestMessageNumberStr, ok := latestMessageNumber.(string)
		if !ok {
			fmt.Println("error converting cache value to int: ", err)
			helpers.ErrorJSON(w, errors.New("invalid message number format"), http.StatusInternalServerError)
			return
		}
		latestMessageNumberInt, err = strconv.Atoi(latestMessageNumberStr)
		if err != nil {
			fmt.Println("error converting cache value to int: ", err)
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	var messageCreationRequestedEventData = models.MessageCreationRequestedEvent{
		ChatID:           chatId,
		Body:             message.Body,
		MessageNumber:    latestMessageNumberInt + 1,
		ChatNumber:       chatNumber,
		ApplicationToken: applicationToken,
	}
	fmt.Println("Publishing message creation requested event with data: ", messageCreationRequestedEventData)
	var messageCreationRequestedEventQueue = events.EventQueues[events.MessageCreationRequestedQueue]
	if err := cmdHandler.EventPublisherManager.PublishEvent(messageCreationRequestedEventQueue, messageCreationRequestedEventData); err != nil {
		log.Printf("Failed to publish event %s: %v", messageCreationRequestedEventQueue.Name, err)
	}

	fmt.Printf("Updating the cache with latestMessageNumber: %d for chatId: %d", messageCreationRequestedEventData.MessageNumber, messageCreationRequestedEventData.ChatID)
	cmdHandler.Cache.Set(cacheKey, messageCreationRequestedEventData.MessageNumber, 0)

	messageResponse := models.MessageResponse{
		Body:          message.Body,
		MessageNumber: latestMessageNumberInt + 1,
		ChatNumber:    chatNumber,
	}
	payload := helpers.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Created Message Successfully for Chat Number %d", chatNumber),
		Data:    messageResponse,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)
}

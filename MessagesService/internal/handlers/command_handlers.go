package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"message_service/internal/caching"
	"message_service/internal/helpers"
	"message_service/internal/messaging"
	"message_service/internal/models"
	"message_service/internal/services"
)

type ApiCmdHandlers struct {
	MessageService *services.MessageService
	Cache          *caching.RedisCacheManager
	eventPublisher *messaging.RabbitMQConn
}

func (cmdHandler *ApiCmdHandlers) CreateMessageCmdHandler(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Print("message: ", message)
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

	// TODO: add validation for correctness of both application token associated with the chat_number
	chatId, err := cmdHandler.MessageService.GetChatId(r.Context(), applicationToken, chatNumber)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if chatId == 0 {
		helpers.ErrorJSON(w, errors.New("chat not found"), http.StatusNotFound)
		return
	}

	// cache key: chat-<chat_id>-latestMessageNo
	cacheKey := fmt.Sprintf("chat-%d-latestMessageNo", chatId)
	latestMessageNumber, err := cmdHandler.Cache.Get(cacheKey)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	latestMessageNumberStr, ok := latestMessageNumber.(string)
	if !ok {
		helpers.ErrorJSON(w, errors.New("invalid message number format"), http.StatusInternalServerError)
		return
	}
	latestMessageNumberInt, err := strconv.Atoi(latestMessageNumberStr)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// TODO: Raise event to create message and process in background
	// (Event raising logic would be implemented here)
	message.Number = latestMessageNumberInt + 1
	message.ChatID = chatId
	// convert message to byte stream
	eventData := []byte(`{"number":` + strconv.Itoa(message.Number) + `,"chat_id":` + strconv.Itoa(message.ChatID) + `,"body":"` + message.Body + `"}`)
	cmdHandler.eventPublisher.PublishEvent("message_created", eventData)
	// if err := cmdHandler.MessageService.CreateMessage(r.Context(), &message, applicationToken, chatNumber); err != nil {
	// 	helpers.ErrorJSON(w, err, http.StatusInternalServerError)
	// 	return
	// }
	// log.Printf("Created message with ID %d for chat %d", message.ID, message.ChatID)
	cmdHandler.Cache.Set(cacheKey, message.Number, 0)

	payload := helpers.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Created Message Successfully for Chat %d", message.ChatID),
		Data:    message,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)
}

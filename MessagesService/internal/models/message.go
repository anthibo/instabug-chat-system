package models

type Message struct {
	ID        int
	Number    int    `json:"number"`
	ChatID    int    `json:"chat_id"`
	Body      string `json:"body" validate:"required"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateMessageBody struct {
	Body string `json:"body" validate:"required"`
}

type MessageResponse struct {
	MessageNumber int `json:"messageNumber"`
	ChatNumber    int `json:"chatNumber"`
	Body          string
}

type MessageCreationRequestedEvent struct {
	MessageNumber    int    `json:"message_number"`
	ChatID           int    `json:"chat_id"`
	ChatNumber       int    `json:"chat_number"`
	Body             string `json:"body"`
	ApplicationToken string `json:"application_token"`
}

type MessageCreatedEvent struct {
	MessageNumber    int    `json:"message_number"`
	ApplicationToken string `json:"application_token"`
	ChatNumber       int    `json:"chat_number"`
	Body             string `json:"body"`
}

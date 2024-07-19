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
	MessageNumber int    `json:"messageNumber"`
	ChatID        int    `json:"chatID"`
	Body          string `json:"body"`
}

type MessageCreatedEvent struct {
	MessageID int    `json:"id"`
	Number    int    `json:"number"`
	ChatID    int    `json:"chat_id"`
	Body      string `json:"body"`
}

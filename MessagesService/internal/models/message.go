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

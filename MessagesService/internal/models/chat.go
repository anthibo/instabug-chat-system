package models

type Chat struct {
	ID            int    `json:"id"`
	ChatNumber    int    `json:"chat_number"`
	ApplicationID int    `json:"application_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

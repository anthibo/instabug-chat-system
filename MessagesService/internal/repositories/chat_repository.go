package repositories

import (
	"context"
	"database/sql"
	"message_service/internal/models"
)

type ChatRepository interface {
	GetChatByChatNumber(ctx context.Context, applicationToken string, chatNumber int) (*models.Chat, error)
}

type MySQLChatRepository struct {
	DB *sql.DB
}

func NewMySQLChatRepository(db *sql.DB) *MySQLMessageRepository {
	return &MySQLMessageRepository{DB: db}
}

func (r *MySQLMessageRepository) GetChatByChatNumber(ctx context.Context, applicationToken string, chatNumber int) (*models.Chat, error) {
	query := `
    SELECT chats.id, chats.number
    FROM chats
    INNER JOIN applications ON applications.id = chats.application_id
    WHERE chats.number = ?
    AND applications.token = ?
    `

	row := singleRowQueryWrapper(ctx, r.DB, query, chatNumber, applicationToken)
	chat := &models.Chat{}
	err := row.Scan(&chat.ID, &chat.ChatNumber)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

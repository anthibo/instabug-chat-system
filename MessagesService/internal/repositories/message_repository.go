package repositories

import (
	"context"
	"database/sql"

	"message_service/internal/models"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *models.Message) error
}

type MySQLMessageRepository struct {
	DB *sql.DB
}

func NewMySQLMessageRepository(db *sql.DB) *MySQLMessageRepository {
	return &MySQLMessageRepository{DB: db}
}

func (r *MySQLMessageRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query := "INSERT INTO messages (chat_id, body, number, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)"
	result, err := commandQueryWrapper(ctx, tx, query, message.ChatID, message.Body, message.Number)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	message.ID = int(id)

	return nil
}

package repositories

import (
	"context"
	"database/sql"
	"log"

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

	msgNumber, err := getMessageNumberTx(ctx, tx, message.ChatID)
	if err != nil {
		return err
	}

	message.Number = msgNumber

	query := "INSERT INTO messages (chat_id, body, number, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)"
	result, err := commandQueryWrapper(ctx, tx, query, message.ChatID, message.Body, message.Number)
	log.Println("query: ", query)
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

func getMessageNumberTx(ctx context.Context, tx *sql.Tx, chatID int) (int, error) {
	query := "SELECT COALESCE(MAX(number), 0) FROM messages WHERE chat_id = ?"

	log.Println("chatId: ", chatID)

	row := singleRowQueryWrapper(ctx, tx, query, chatID)

	var number int
	err := row.Scan(&number)
	if err != nil {
		return 0, err
	}

	return number + 1, nil
}

// func (r *MySQLMessageRepository) UpdateMessage(ctx context.Context, msgNo int, message *models.Message) error {
// 	query := "UPDATE messages SET body = ? WHERE number = ? and chat_id = ?"
// 	result, err := r.DB.ExecContext(ctx, query, message.Body, msgNo, message.ChatID)
// 	if err != nil {
// 		log.Println("Error updating message: ", err)
// 		return err
// 	}

// 	affectedRows, err := result.RowsAffected()
// 	if err != nil {
// 		log.Println("Error updating message: ", err)
// 		return err
// 	}

// 	if affectedRows == 0 {
// 		log.Println("No message found to update")
// 		return nil
// 	}

// 	return nil
// }

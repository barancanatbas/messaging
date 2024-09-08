package message

import (
	"database/sql"
	"github.com/barancanatbas/messaging/internal/domain/entity"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUnsentMessages(lastMessageId int) (*sql.Rows, error) {
	var query string
	var args []interface{}

	if lastMessageId > 0 {
		query = `SELECT id, content, phone_number, sent_at 
		         FROM messages 
		         WHERE sent_at IS NULL AND id > ? 
		         ORDER BY id ASC`
		args = append(args, lastMessageId)
	} else {
		query = `SELECT id, content, phone_number, sent_at 
		         FROM messages 
		         WHERE sent_at IS NULL 
		         ORDER BY id ASC`
	}

	return r.db.Query(query, args...)
}

func (r *Repository) Scan(rows *sql.Rows, msg *entity.Message) error {
	return rows.Scan(&msg.ID, &msg.Content, &msg.PhoneNumber, &msg.SentAt)
}

func (r *Repository) MarkMessageAsSent(id int, uuid string, sentAt time.Time) error {
	query := `UPDATE messages SET sent_at = ?, uuid = ?, status = ? WHERE id = ?`
	_, err := r.db.Exec(query, sentAt, uuid, entity.StatusDelivered, id)
	return err
}

func (r *Repository) GetSentMessages() ([]entity.Message, error) {
	query := `SELECT id, content, phone_number, sent_at, status, uuid FROM messages WHERE sent_at IS NOT NULL`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.PhoneNumber, &msg.SentAt, &msg.Status, &msg.UUID); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) CreateMessage(msg *entity.Message) error {
	query := `INSERT INTO messages (content, phone_number, status) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, msg.Content, msg.PhoneNumber, msg.Status)
	return err
}

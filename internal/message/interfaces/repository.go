package interfaces

import (
	"database/sql"
	"github.com/barancanatbas/messaging/internal/domain/entity"
	"time"
)

type MessageRepository interface {
	MarkMessageAsSent(id int, uuid string, sentAt time.Time) error
	GetSentMessages() ([]entity.Message, error)
	Scan(rows *sql.Rows, msg *entity.Message) error
	GetUnsentMessages(lastMessageId int) (*sql.Rows, error)
	CreateMessage(msg *entity.Message) error
}

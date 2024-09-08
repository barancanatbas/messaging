package entity

import "time"

type MessageStatus string

const (
	StatusPending   MessageStatus = "PENDING"
	StatusFailed    MessageStatus = "FAILED"
	StatusDelivered MessageStatus = "DELIVERED"
)

type Message struct {
	ID          int
	Content     string
	PhoneNumber string
	Status      MessageStatus
	SentAt      *time.Time
	UUID        string
}

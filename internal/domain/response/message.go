package response

import "time"

type Message struct {
	ID          string     `json:"id"`
	Content     string     `json:"content"`
	PhoneNumber string     `json:"phone_number"`
	Status      string     `json:"status"`
	Sent        bool       `json:"sent"`
	SentAt      *time.Time `json:"sent_at"`
}

package response

import "time"

type MessageResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
	SentAt    time.Time
}

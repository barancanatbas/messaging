package interfaces

import (
	"github.com/barancanatbas/messaging/internal/domain/request"
	"github.com/barancanatbas/messaging/internal/domain/response"
)

type DeliveryService interface {
	SendMessage(data *request.MessageRequest) (*response.MessageResponse, error)
}

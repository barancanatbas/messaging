package delivery

import (
	"encoding/json"
	"github.com/barancanatbas/messaging/internal/delivery/interfaces"
	"github.com/barancanatbas/messaging/internal/domain/request"
	"github.com/barancanatbas/messaging/internal/domain/response"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Service struct {
	httpClient interfaces.HttpClient
}

func NewDeliveryService(httpClient interfaces.HttpClient) *Service {
	return &Service{
		httpClient: httpClient,
	}
}

func (s *Service) SendMessage(data *request.MessageRequest) (*response.MessageResponse, error) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message request")
		return nil, err
	}

	responseBody, err := s.httpClient.Send(http.MethodPost, requestBody)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send message to %s", data.To)
		return nil, err
	}

	var messageResponse response.MessageResponse
	err = json.Unmarshal(responseBody, &messageResponse)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal message response")
		return nil, err
	}

	messageResponse.SentAt = time.Now()

	return &messageResponse, nil
}

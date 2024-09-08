package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/barancanatbas/messaging/internal/domain/request"
	"github.com/barancanatbas/messaging/internal/domain/response"
	"github.com/barancanatbas/messaging/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendMessage_Success(t *testing.T) {
	mockHttpClient := new(mocks.HttpClient)

	messageResponse := &response.MessageResponse{
		MessageID: "123",
	}
	responseBody, _ := json.Marshal(messageResponse)

	mockHttpClient.On("Send", http.MethodPost, mock.AnythingOfType("[]uint8")).Return(responseBody, nil)

	deliveryService := NewDeliveryService(mockHttpClient)

	messageRequest := &request.MessageRequest{
		To:      "+123456789",
		Content: "Hello!",
	}

	resp, err := deliveryService.SendMessage(messageRequest)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "123", resp.MessageID)

	mockHttpClient.AssertExpectations(t)
}

func TestSendMessage_HttpError(t *testing.T) {
	mockHttpClient := new(mocks.HttpClient)

	mockHttpClient.On("Send", http.MethodPost, mock.AnythingOfType("[]uint8")).Return(nil, errors.New("HTTP error"))

	deliveryService := NewDeliveryService(mockHttpClient)

	messageRequest := &request.MessageRequest{
		To:      "+123456789",
		Content: "Hello!",
	}

	resp, err := deliveryService.SendMessage(messageRequest)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "HTTP error", err.Error())

	mockHttpClient.AssertExpectations(t)
}

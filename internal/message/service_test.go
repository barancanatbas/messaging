package message

import (
	"github.com/barancanatbas/messaging/config"
	"github.com/barancanatbas/messaging/internal/domain/entity"
	"github.com/barancanatbas/messaging/internal/domain/request"
	"github.com/barancanatbas/messaging/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	mockRepo := new(mocks.MessageRepository)
	mockCache := new(mocks.Cache)
	mockDeliveryService := new(mocks.DeliveryService)

	service := NewMessageService(config.AppConfig{MessageInterval: 10}, mockRepo, mockCache, mockDeliveryService)

	req := request.CreateMessageRequest{
		PhoneNumber: "+123456789",
		Content:     "Hello World!",
	}

	expectedMessage := entity.Message{
		PhoneNumber: req.PhoneNumber,
		Content:     req.Content,
		Status:      entity.StatusPending,
	}

	mockRepo.On("CreateMessage", mock.AnythingOfType("*entity.Message")).Return(nil).Run(func(args mock.Arguments) {
		msg := args.Get(0).(*entity.Message)
		assert.Equal(t, expectedMessage.PhoneNumber, msg.PhoneNumber)
		assert.Equal(t, expectedMessage.Content, msg.Content)
		assert.Equal(t, entity.StatusPending, msg.Status)
	})

	err := service.CreateMessage(req)

	mockRepo.AssertExpectations(t)

	assert.NoError(t, err)
}

func TestCreateMessage_Error(t *testing.T) {
	mockRepo := new(mocks.MessageRepository)
	mockCache := new(mocks.Cache)
	mockDeliveryService := new(mocks.DeliveryService)

	service := NewMessageService(config.AppConfig{MessageInterval: 10}, mockRepo, mockCache, mockDeliveryService)

	req := request.CreateMessageRequest{
		PhoneNumber: "+123456789",
		Content:     "Hello World!",
	}

	mockRepo.On("CreateMessage", mock.AnythingOfType("*entity.Message")).Return(assert.AnError)

	err := service.CreateMessage(req)

	mockRepo.AssertExpectations(t)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestGetSentMessages(t *testing.T) {
	mockCache := new(mocks.Cache)

	sentMessages := []entity.Message{
		{ID: 1, PhoneNumber: "+123456789", Content: "Test message 1"},
		{ID: 2, PhoneNumber: "+987654321", Content: "Test message 2"},
	}

	mockCache.On("LRange", SentMessageCacheKey, int64(0), int64(-1), mock.AnythingOfType("*[]entity.Message")).
		Run(func(args mock.Arguments) {
			result := args.Get(3).(*[]entity.Message)
			*result = sentMessages
		}).
		Return(nil)

	service := Service{
		cache: mockCache,
	}

	messages, err := service.GetSentMessages()

	assert.NoError(t, err)
	assert.Equal(t, sentMessages, messages)

	mockCache.AssertExpectations(t)
}

func TestGetSentMessages_Error(t *testing.T) {
	mockCache := new(mocks.Cache)

	mockError := assert.AnError

	mockCache.On("LRange", SentMessageCacheKey, int64(0), int64(-1), mock.AnythingOfType("*[]entity.Message")).
		Return(mockError)

	service := Service{
		cache: mockCache,
	}

	messages, err := service.GetSentMessages()

	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.Equal(t, mockError, err)

	mockCache.AssertExpectations(t)
}

func TestStopAutomaticSending_ChannelEmpty(t *testing.T) {
	stopChan := make(chan bool, 1)

	service := Service{
		stopChan: stopChan,
	}

	service.StopAutomaticSending()

	select {
	case val := <-stopChan:
		assert.True(t, val)
	default:
		t.Error("Expected true to be sent to stopChan, but nothing was sent.")
	}
}

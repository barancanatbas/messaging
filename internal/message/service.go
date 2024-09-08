package message

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/barancanatbas/messaging/internal/domain/request"
	"github.com/barancanatbas/messaging/internal/domain/response"
	"github.com/barancanatbas/messaging/internal/message/interfaces"
	"sync"
	"time"

	"github.com/barancanatbas/messaging/config"
	"github.com/barancanatbas/messaging/internal/domain/entity"
	"github.com/rs/zerolog/log"
)

const SentMessageCacheKey = "sent_messages"

type Service struct {
	repository      interfaces.MessageRepository
	cache           interfaces.Cache
	ticker          *time.Ticker
	stopChan        chan bool
	wg              sync.WaitGroup
	deliveryService interfaces.DeliveryService
	messageInterval time.Duration
}

func NewMessageService(config config.AppConfig, repository interfaces.MessageRepository, cache interfaces.Cache, deliveryService interfaces.DeliveryService) *Service {
	service := &Service{
		repository:      repository,
		cache:           cache,
		stopChan:        make(chan bool),
		deliveryService: deliveryService,
		messageInterval: time.Duration(config.MessageInterval) * time.Minute,
	}

	if err := service.StartAutomaticSending(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start automatic sending")
	}

	return service
}

func (s *Service) CreateMessage(msg request.CreateMessageRequest) error {
	var message entity.Message

	message.PhoneNumber = msg.PhoneNumber
	message.Content = msg.Content
	message.Status = entity.StatusPending

	return s.repository.CreateMessage(&message)
}

func (s *Service) GetSentMessages() ([]entity.Message, error) {
	messages, err := s.repository.GetSentMessages()

	return messages, err
}

func (s *Service) StartAutomaticSending() error {
	s.ticker = time.NewTicker(s.messageInterval)
	go s.runAutomaticSending()
	return nil
}

func (s *Service) StopAutomaticSending() {
	select {
	case s.stopChan <- true:
	default:
		log.Info().Msg("Automatic sending already stopped")
	}
}

func (s *Service) runAutomaticSending() {
	var cancelFunc context.CancelFunc

	for {
		select {
		case <-s.ticker.C:
			s.handleSendMessages(&cancelFunc)
		case stop := <-s.stopChan:
			if stop {
				s.handleStopSending(cancelFunc)
				return
			}
		}
	}
}

func (s *Service) handleSendMessages(cancelFunc *context.CancelFunc) {
	if *cancelFunc != nil {
		(*cancelFunc)()
	}

	ctx, newCancelFunc := context.WithCancel(context.Background())
	*cancelFunc = newCancelFunc

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.sendMessages(ctx)
	}()
}

func (s *Service) handleStopSending(cancelFunc context.CancelFunc) {
	log.Info().Msg("Stopping automatic message sending...")

	if cancelFunc != nil {
		cancelFunc()
	}

	s.wg.Wait()
	s.ticker.Stop()
}

func (s *Service) sendMessages(ctx context.Context) {
	log.Info().Msg("Processing unsent messages...")

	lastMessageID := s.getLastMessageID()

	rows, err := s.repository.GetUnsentMessages(lastMessageID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve unsent messages")
		return
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-s.stopChan:
			return
		case <-ctx.Done():
			log.Info().Msg("sendMessages operation canceled")
			return
		default:
			s.processMessage(rows)
		}
	}

	log.Info().Msg("All unsent messages processed")
}

func (s *Service) processMessage(rows *sql.Rows) {
	var msg entity.Message
	if err := s.repository.Scan(rows, &msg); err != nil {
		log.Error().Err(err).Msg("Failed to scan message")
		return
	}

	if err := s.sendMessage(&msg); err != nil {
		log.Error().Err(err).Msgf("Failed to send message to %s", msg.PhoneNumber)
		return
	}

	log.Info().Msgf("Message sent to %s", msg.PhoneNumber)
}

func (s *Service) sendMessage(msg *entity.Message) error {
	resp, err := s.deliveryService.SendMessage(&request.MessageRequest{
		To:      msg.PhoneNumber,
		Content: msg.Content,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
		return err
	}

	s.cacheSentMessage(resp)

	if err := s.repository.MarkMessageAsSent(msg.ID, resp.MessageID, resp.SentAt); err != nil {
		log.Error().Err(err).Msg("Failed to mark message as sent")
		return err
	}

	return nil
}

func (s *Service) cacheSentMessage(resp *response.MessageResponse) {
	if s.cache == nil {
		return
	}

	cacheValue := response.SentMessageCache{
		ID:     resp.MessageID,
		SentAt: resp.SentAt.Format(time.RFC3339),
	}

	cacheValueJSON, err := json.Marshal(cacheValue)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message struct to JSON")
		return
	}

	if err := s.cache.LPush(SentMessageCacheKey, cacheValueJSON); err != nil {
		log.Error().Err(err).Msg("Failed to cache last sent message")
	}
}

func (s *Service) getLastMessageID() int {
	if s.cache == nil {
		return 0
	}

	var lastMessage entity.Message
	err := s.cache.LIndex(SentMessageCacheKey, 0, &lastMessage)
	if err != nil {
		return 0
	}

	return lastMessage.ID
}

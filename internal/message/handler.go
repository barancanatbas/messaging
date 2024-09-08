package message

import (
	"github.com/barancanatbas/messaging/internal/domain/request"
	"github.com/barancanatbas/messaging/internal/domain/response"
	"github.com/barancanatbas/messaging/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	ActionStart = "start"
	ActionStop  = "stop"
)

type Handler struct {
	app       *fiber.App
	service   *Service
	validator *validator.Service
}

func NewMessageHandler(app *fiber.App, service *Service, validator *validator.Service) *Handler {
	handler := &Handler{
		app:       app,
		service:   service,
		validator: validator,
	}

	handler.RegisterRoutes()

	return handler
}

func (h *Handler) RegisterRoutes() {
	h.app.Post("/send-messages", h.StartMessageSending)
	h.app.Get("/sent-messages", h.GetSentMessages)
	h.app.Post("/message", h.CreateMessage)
}

// StartMessageSending starts or stops automatic message sending
// @Summary Start or stop automatic message sending
// @Description Start or stop automatic message sending based on the query parameter 'action'.
// @Tags messages
// @Param action query string true "Action to perform (start or stop)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /send-messages [post]
func (h *Handler) StartMessageSending(c *fiber.Ctx) error {
	action := c.Query("action")

	switch action {
	case ActionStart:
		if err := h.service.StartAutomaticSending(); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(response.Error("Failed to start automatic message sending", nil))
		}
		return c.Status(http.StatusOK).JSON(response.Success("Automatic message sending started", nil))

	case ActionStop:
		h.service.StopAutomaticSending()
		return c.Status(http.StatusOK).JSON(response.Success("Automatic message sending stopped", nil))

	default:
		return c.Status(http.StatusBadRequest).JSON(response.Error("Invalid action, use ?action=start or ?action=stop", nil))
	}
}

// GetSentMessages retrieves sent messages
// @Summary Get sent messages
// @Description Retrieve all sent messages from the system
// @Tags messages
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /sent-messages [get]
func (h *Handler) GetSentMessages(c *fiber.Ctx) error {
	messages, err := h.service.GetSentMessages()
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve sent messages")
		return c.Status(http.StatusInternalServerError).JSON(response.Error("Failed to retrieve sent messages", nil))
	}

	return c.Status(http.StatusOK).JSON(response.Success("Sent messages retrieved successfully", messages))
}

// CreateMessage creates a new message
// @Summary Create a new message
// @Description Create and send a new message
// @Tags messages
// @Param message body request.CreateMessageRequest true "Create message request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /message [post]
func (h *Handler) CreateMessage(c *fiber.Ctx) error {
	var message request.CreateMessageRequest
	if err := h.validator.ParseAndValidate(c, &message); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Error(err.Error(), nil))
	}

	err := h.service.CreateMessage(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create SMS")
		return c.Status(http.StatusInternalServerError).JSON(response.Error("Failed to create SMS", nil))
	}

	return c.Status(http.StatusCreated).JSON(response.Success("SMS created successfully", nil))
}

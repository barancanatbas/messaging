package request

type CreateMessageRequest struct {
	Content     string `json:"content" validate:"required,max=255"`
	PhoneNumber string `json:"phone_number" validate:"required,max=20"`
}

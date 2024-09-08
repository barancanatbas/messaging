package request

type MessageRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

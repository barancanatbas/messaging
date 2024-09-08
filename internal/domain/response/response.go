package response

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(message string, data interface{}) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func Error(message string, data interface{}) Response {
	return Response{
		Status:  "error",
		Message: message,
		Data:    data,
	}
}

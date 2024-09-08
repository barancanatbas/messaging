package interfaces

type HttpClient interface {
	Send(method string, body []byte) ([]byte, error)
}

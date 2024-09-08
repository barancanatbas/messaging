package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/barancanatbas/messaging/config"
	"github.com/rs/zerolog/log"
)

type HttpClient struct {
	BaseURL string
	AuthKey string
	client  *http.Client
}

func NewHttpClient(cfg config.HttpClientConfig) *HttpClient {
	return &HttpClient{
		BaseURL: cfg.BaseURL,
		AuthKey: cfg.AuthKey,
		client:  &http.Client{},
	}
}

func (c *HttpClient) Send(method string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, c.BaseURL, bytes.NewBuffer(body))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", c.AuthKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send request")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}

package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RealClient реальная реализация Telegram клиента
type RealClient struct {
	httpClient *http.Client
}

// NewRealClient создает новый реальный Telegram клиент
func NewRealClient() *RealClient {
	return &RealClient{
		httpClient: &http.Client{},
	}
}

// SendMessage отправляет сообщение через Telegram Bot API
func (r *RealClient) SendMessage(token, chatID, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	return nil
}

package telegram

import (
	"fmt"
	"log"
)

// MockClient мок-реализация Telegram клиента для разработки и тестирования
type MockClient struct {
	enabled bool
}

// NewMockClient создает новый мок-клиент
func NewMockClient() *MockClient {
	return &MockClient{
		enabled: true,
	}
}

// SendMessage имитирует отправку сообщения в Telegram
// В мок-режиме просто логирует сообщение
func (m *MockClient) SendMessage(token, chatID, message string) error {
	if !m.enabled {
		return fmt.Errorf("mock client is disabled")
	}

	// Логируем вместо реальной отправки
	log.Printf("[MOCK TELEGRAM] Sending message to chat_id=%s: %s", chatID, message)
	log.Printf("[MOCK TELEGRAM] Token: %s...%s", token[:min(10, len(token))], token[max(0, len(token)-5):])

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

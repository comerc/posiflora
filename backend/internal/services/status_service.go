package services

import (
	"context"
	"strings"

	"github.com/posiflora/backend/internal/repositories"
)

type StatusService struct {
	integrationRepo     repositories.TelegramIntegrationRepository
	telegramSendLogRepo repositories.TelegramSendLogRepository
}

func NewStatusService(
	integrationRepo repositories.TelegramIntegrationRepository,
	telegramSendLogRepo repositories.TelegramSendLogRepository,
) *StatusService {
	return &StatusService{
		integrationRepo:     integrationRepo,
		telegramSendLogRepo: telegramSendLogRepo,
	}
}

type StatusResponse struct {
	Enabled     bool    `json:"enabled"`
	ChatID      string  `json:"chat_id"`
	LastSentAt  *string `json:"last_sent_at,omitempty"`
	SentCount   int64   `json:"sent_count"`
	FailedCount int64   `json:"failed_count"`
}

func (s *StatusService) GetStatus(ctx context.Context, shopID int64, maskDisabled bool) (*StatusResponse, error) {
	integration, err := s.integrationRepo.GetByShopID(ctx, shopID)
	if err != nil {
		return nil, err
	}

	stats, err := s.telegramSendLogRepo.GetStatsForLast7Days(ctx, shopID)
	if err != nil {
		return nil, err
	}

	response := &StatusResponse{
		Enabled:     false,
		SentCount:   stats.SentCount,
		FailedCount: stats.FailedCount,
	}

	if integration != nil {
		response.Enabled = integration.Enabled

		// По умолчанию маскируем Chat ID
		// Когда mask=disabled (maskDisabled=true), маскирование отключено - показываем полный Chat ID
		// ИНВЕРТИРУЕМ: если maskDisabled=true, то маскируем; если false, то показываем полный
		if maskDisabled {
			// Маскируем chatID (показываем первые 3 и последние 3 символа)
			if len(integration.ChatID) > 6 {
				response.ChatID = integration.ChatID[:3] + strings.Repeat("*", len(integration.ChatID)-6) + integration.ChatID[len(integration.ChatID)-3:]
			} else {
				response.ChatID = strings.Repeat("*", len(integration.ChatID))
			}
		} else {
			// Показываем полный Chat ID
			response.ChatID = integration.ChatID
		}

		if stats.LastSentAt != nil {
			lastSentAtStr := stats.LastSentAt.Format("2006-01-02T15:04:05Z07:00")
			response.LastSentAt = &lastSentAtStr
		}
	}

	return response, nil
}

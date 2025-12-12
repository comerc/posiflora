package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/posiflora/backend/internal/models"
	"github.com/posiflora/backend/internal/repositories"
)

type TelegramIntegrationService struct {
	integrationRepo repositories.TelegramIntegrationRepository
	shopRepo        repositories.ShopRepository
}

func NewTelegramIntegrationService(
	integrationRepo repositories.TelegramIntegrationRepository,
	shopRepo repositories.ShopRepository,
) *TelegramIntegrationService {
	return &TelegramIntegrationService{
		integrationRepo: integrationRepo,
		shopRepo:        shopRepo,
	}
}

type ConnectRequest struct {
	BotToken string `json:"botToken" binding:"required"`
	ChatID   string `json:"chatId" binding:"required"`
	Enabled  bool   `json:"enabled"`
}

func (s *TelegramIntegrationService) Connect(ctx context.Context, shopID int64, req ConnectRequest) (*models.TelegramIntegration, error) {
	if req.BotToken == "" || req.ChatID == "" {
		return nil, errors.New("botToken and chatId are required")
	}

	// Создаем магазин, если его нет
	_, err := s.shopRepo.GetOrCreate(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create shop: %w", err)
	}

	now := time.Now()
	integration := &models.TelegramIntegration{
		ShopID:    shopID,
		BotToken:  req.BotToken,
		ChatID:    req.ChatID,
		Enabled:   req.Enabled,
		UpdatedAt: now,
	}

	existing, err := s.integrationRepo.GetByShopID(ctx, shopID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		integration.ID = existing.ID
		integration.CreatedAt = existing.CreatedAt
	} else {
		integration.CreatedAt = now
	}

	err = s.integrationRepo.Upsert(ctx, integration)
	if err != nil {
		return nil, err
	}

	return integration, nil
}

func (s *TelegramIntegrationService) GetStatus(ctx context.Context, shopID int64) (*models.TelegramIntegration, error) {
	return s.integrationRepo.GetByShopID(ctx, shopID)
}

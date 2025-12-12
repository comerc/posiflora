package services

import (
	"context"
	"fmt"
	"time"

	"github.com/posiflora/backend/internal/models"
	"github.com/posiflora/backend/internal/repositories"
	"github.com/posiflora/backend/internal/telegram"
)

type OrderService struct {
	orderRepo           repositories.OrderRepository
	integrationRepo     repositories.TelegramIntegrationRepository
	telegramSendLogRepo repositories.TelegramSendLogRepository
	telegramClient      telegram.Client
	shopRepo            repositories.ShopRepository
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	integrationRepo repositories.TelegramIntegrationRepository,
	telegramSendLogRepo repositories.TelegramSendLogRepository,
	telegramClient telegram.Client,
	shopRepo repositories.ShopRepository,
) *OrderService {
	return &OrderService{
		orderRepo:           orderRepo,
		integrationRepo:     integrationRepo,
		telegramSendLogRepo: telegramSendLogRepo,
		telegramClient:      telegramClient,
		shopRepo:            shopRepo,
	}
}

type CreateOrderRequest struct {
	Number       string  `json:"number" binding:"required"`
	Total        float64 `json:"total" binding:"required"`
	CustomerName string  `json:"customerName" binding:"required"`
}

type CreateOrderResponse struct {
	Order      *models.Order `json:"order"`
	SendStatus string        `json:"send_status"` // sent, failed, skipped
}

func (s *OrderService) CreateOrder(ctx context.Context, shopID int64, req CreateOrderRequest) (*CreateOrderResponse, error) {
	// Создаем магазин, если его нет
	_, err := s.shopRepo.GetOrCreate(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create shop: %w", err)
	}

	// 1. Создаем заказ
	order := &models.Order{
		ShopID:       shopID,
		Number:       req.Number,
		Total:        req.Total,
		CustomerName: req.CustomerName,
		CreatedAt:    time.Now(),
	}

	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 2. Проверяем наличие активной Telegram-интеграции
	integration, err := s.integrationRepo.GetByShopID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	// 3. Проверяем идемпотентность
	existingLog, err := s.telegramSendLogRepo.GetByShopIDAndOrderID(ctx, shopID, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing log: %w", err)
	}

	if existingLog != nil {
		// Уже отправляли, пропускаем
		return &CreateOrderResponse{
			Order:      order,
			SendStatus: "skipped",
		}, nil
	}

	// 4. Если интеграция включена - отправляем уведомление
	if integration != nil && integration.Enabled {
		message := fmt.Sprintf("Новый заказ %s на сумму %.2f ₽, клиент %s", order.Number, order.Total, order.CustomerName)

		sendLog := &models.TelegramSendLog{
			ShopID:  shopID,
			OrderID: order.ID,
			Message: message,
			Status:  models.StatusSent,
			SentAt:  time.Now(),
		}

		err = s.telegramClient.SendMessage(integration.BotToken, integration.ChatID, message)
		if err != nil {
			// Ошибка не должна ронять создание заказа
			errorMsg := err.Error()
			sendLog.Status = models.StatusFailed
			sendLog.Error = &errorMsg
		}

		// Логируем попытку отправки
		err = s.telegramSendLogRepo.Create(ctx, sendLog)
		if err != nil {
			// Даже если не удалось записать лог, заказ уже создан
			// В продакшене здесь можно добавить дополнительное логирование
		}

		status := "sent"
		if sendLog.Status == models.StatusFailed {
			status = "failed"
		}

		return &CreateOrderResponse{
			Order:      order,
			SendStatus: status,
		}, nil
	}

	// Интеграция не настроена или отключена
	return &CreateOrderResponse{
		Order:      order,
		SendStatus: "skipped",
	}, nil
}

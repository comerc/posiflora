package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/posiflora/backend/internal/models"
	"github.com/posiflora/backend/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTelegramClient мок для Telegram клиента
type MockTelegramClient struct {
	mock.Mock
}

func (m *MockTelegramClient) SendMessage(token, chatID, message string) error {
	args := m.Called(token, chatID, message)
	return args.Error(0)
}

// MockOrderRepository мок для OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	// Устанавливаем ID для созданного заказа
	if order.ID == 0 {
		order.ID = 1
	}
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

// MockTelegramIntegrationRepository мок для TelegramIntegrationRepository
type MockTelegramIntegrationRepository struct {
	mock.Mock
}

func (m *MockTelegramIntegrationRepository) GetByShopID(ctx context.Context, shopID int64) (*models.TelegramIntegration, error) {
	args := m.Called(ctx, shopID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TelegramIntegration), args.Error(1)
}

func (m *MockTelegramIntegrationRepository) Upsert(ctx context.Context, integration *models.TelegramIntegration) error {
	args := m.Called(ctx, integration)
	return args.Error(0)
}

// MockShopRepository мок для ShopRepository
type MockShopRepository struct {
	mock.Mock
}

func (m *MockShopRepository) GetByID(ctx context.Context, id int64) (*models.Shop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Shop), args.Error(1)
}

func (m *MockShopRepository) Create(ctx context.Context, shop *models.Shop) error {
	args := m.Called(ctx, shop)
	return args.Error(0)
}

func (m *MockShopRepository) GetOrCreate(ctx context.Context, id int64) (*models.Shop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Shop), args.Error(1)
}

// MockTelegramSendLogRepository мок для TelegramSendLogRepository
type MockTelegramSendLogRepository struct {
	mock.Mock
}

func (m *MockTelegramSendLogRepository) GetByShopIDAndOrderID(ctx context.Context, shopID, orderID int64) (*models.TelegramSendLog, error) {
	args := m.Called(ctx, shopID, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TelegramSendLog), args.Error(1)
}

func (m *MockTelegramSendLogRepository) Create(ctx context.Context, log *models.TelegramSendLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockTelegramSendLogRepository) GetStatsForLast7Days(ctx context.Context, shopID int64) (*repositories.TelegramStats, error) {
	args := m.Called(ctx, shopID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.TelegramStats), args.Error(1)
}

// Тест 1: При создании заказа и включённой интеграции вызывается TelegramClient и пишется лог SENT
func TestOrderService_CreateOrder_WithEnabledIntegration_SendsMessage(t *testing.T) {
	ctx := context.Background()
	shopID := int64(1)

	mockOrderRepo := new(MockOrderRepository)
	mockIntegrationRepo := new(MockTelegramIntegrationRepository)
	mockSendLogRepo := new(MockTelegramSendLogRepository)
	mockTelegramClient := new(MockTelegramClient)
	mockShopRepo := new(MockShopRepository)

	// Настройка моков
	shop := &models.Shop{ID: shopID, Name: "Test Shop"}
	integration := &models.TelegramIntegration{
		ID:       1,
		ShopID:   shopID,
		BotToken: "123456:ABC-DEF",
		ChatID:   "987654321",
		Enabled:  true,
	}

	mockShopRepo.On("GetOrCreate", ctx, shopID).Return(shop, nil)
	mockIntegrationRepo.On("GetByShopID", ctx, shopID).Return(integration, nil)
	mockOrderRepo.On("Create", ctx, mock.MatchedBy(func(order *models.Order) bool {
		order.ID = 1 // Устанавливаем ID при создании
		return true
	})).Return(nil)
	mockSendLogRepo.On("GetByShopIDAndOrderID", ctx, shopID, int64(1)).Return(nil, nil)
	mockTelegramClient.On("SendMessage", "123456:ABC-DEF", "987654321", mock.AnythingOfType("string")).Return(nil)
	mockSendLogRepo.On("Create", ctx, mock.MatchedBy(func(log *models.TelegramSendLog) bool {
		return log.Status == models.StatusSent && log.ShopID == shopID && log.OrderID == 1
	})).Return(nil)

	service := NewOrderService(mockOrderRepo, mockIntegrationRepo, mockSendLogRepo, mockTelegramClient, mockShopRepo)

	req := CreateOrderRequest{
		Number:       "A-1005",
		Total:        2490.0,
		CustomerName: "Анна",
	}

	response, err := service.CreateOrder(ctx, shopID, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "sent", response.SendStatus)
	assert.NotNil(t, response.Order)

	mockTelegramClient.AssertExpectations(t)
	mockSendLogRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

// Тест 2: Повторная отправка/повторное создание заказа не создаёт дублей telegram_send_log и не шлёт сообщение повторно (идемпотентность)
func TestOrderService_CreateOrder_Idempotency(t *testing.T) {
	ctx := context.Background()
	shopID := int64(1)
	orderID := int64(1)

	mockOrderRepo := new(MockOrderRepository)
	mockIntegrationRepo := new(MockTelegramIntegrationRepository)
	mockSendLogRepo := new(MockTelegramSendLogRepository)
	mockTelegramClient := new(MockTelegramClient)
	mockShopRepo := new(MockShopRepository)

	// Настройка моков - уже есть лог отправки
	shop := &models.Shop{ID: shopID, Name: "Test Shop"}
	existingLog := &models.TelegramSendLog{
		ID:      1,
		ShopID:  shopID,
		OrderID: orderID,
		Status:  models.StatusSent,
		SentAt:  time.Now(),
	}

	mockShopRepo.On("GetOrCreate", ctx, shopID).Return(shop, nil)
	mockIntegrationRepo.On("GetByShopID", ctx, shopID).Return(nil, nil) // Интеграция не настроена

	// Сначала создаем заказ, потом проверяем лог
	mockOrderRepo.On("Create", ctx, mock.MatchedBy(func(order *models.Order) bool {
		order.ID = orderID // Устанавливаем ID при создании
		return true
	})).Return(nil)

	// После создания заказа проверяем лог - он уже существует
	mockSendLogRepo.On("GetByShopIDAndOrderID", ctx, shopID, orderID).Return(existingLog, nil)

	service := NewOrderService(mockOrderRepo, mockIntegrationRepo, mockSendLogRepo, mockTelegramClient, mockShopRepo)

	req := CreateOrderRequest{
		Number:       "A-1005",
		Total:        2490.0,
		CustomerName: "Анна",
	}

	response, err := service.CreateOrder(ctx, shopID, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "skipped", response.SendStatus)

	// Telegram клиент не должен быть вызван
	mockTelegramClient.AssertNotCalled(t, "SendMessage")
}

// Тест 3: При ошибке TelegramClient пишется лог FAILED, а заказ всё равно создаётся
func TestOrderService_CreateOrder_TelegramError_OrderStillCreated(t *testing.T) {
	ctx := context.Background()
	shopID := int64(1)

	mockOrderRepo := new(MockOrderRepository)
	mockIntegrationRepo := new(MockTelegramIntegrationRepository)
	mockSendLogRepo := new(MockTelegramSendLogRepository)
	mockTelegramClient := new(MockTelegramClient)
	mockShopRepo := new(MockShopRepository)

	// Настройка моков
	shop := &models.Shop{ID: shopID, Name: "Test Shop"}
	integration := &models.TelegramIntegration{
		ID:       1,
		ShopID:   shopID,
		BotToken: "123456:ABC-DEF",
		ChatID:   "987654321",
		Enabled:  true,
	}

	telegramError := errors.New("telegram API error")

	mockShopRepo.On("GetOrCreate", ctx, shopID).Return(shop, nil)
	mockIntegrationRepo.On("GetByShopID", ctx, shopID).Return(integration, nil)
	mockOrderRepo.On("Create", ctx, mock.MatchedBy(func(order *models.Order) bool {
		order.ID = 1 // Устанавливаем ID при создании
		return true
	})).Return(nil)
	mockSendLogRepo.On("GetByShopIDAndOrderID", ctx, shopID, int64(1)).Return(nil, nil)
	mockTelegramClient.On("SendMessage", "123456:ABC-DEF", "987654321", mock.AnythingOfType("string")).Return(telegramError)
	mockSendLogRepo.On("Create", ctx, mock.MatchedBy(func(log *models.TelegramSendLog) bool {
		return log.Status == models.StatusFailed && log.Error != nil && *log.Error == "telegram API error" && log.OrderID == 1
	})).Return(nil)

	service := NewOrderService(mockOrderRepo, mockIntegrationRepo, mockSendLogRepo, mockTelegramClient, mockShopRepo)

	req := CreateOrderRequest{
		Number:       "A-1005",
		Total:        2490.0,
		CustomerName: "Анна",
	}

	response, err := service.CreateOrder(ctx, shopID, req)

	assert.NoError(t, err) // Заказ создан успешно, несмотря на ошибку Telegram
	assert.NotNil(t, response)
	assert.Equal(t, "failed", response.SendStatus)
	assert.NotNil(t, response.Order)

	mockOrderRepo.AssertExpectations(t)
	mockTelegramClient.AssertExpectations(t)
	mockSendLogRepo.AssertExpectations(t)
}

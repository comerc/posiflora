package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/posiflora/backend/config"
	"github.com/posiflora/backend/internal/handlers"
	"github.com/posiflora/backend/internal/repositories"
	"github.com/posiflora/backend/internal/services"
	"github.com/posiflora/backend/internal/telegram"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к БД
	sqldb, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, pgdialect.New())

	// Включаем отладочный режим в dev
	if os.Getenv("DEBUG") == "true" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	// Проверка подключения
	ctx := context.Background()
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализация Telegram клиента
	telegramClient := telegram.NewClient(cfg.Telegram.MockMode)

	// Инициализация репозиториев
	shopRepo := repositories.NewShopRepository(db)
	integrationRepo := repositories.NewTelegramIntegrationRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	telegramSendLogRepo := repositories.NewTelegramSendLogRepository(db)

	// Инициализация сервисов
	telegramIntegrationService := services.NewTelegramIntegrationService(integrationRepo, shopRepo)
	orderService := services.NewOrderService(orderRepo, integrationRepo, telegramSendLogRepo, telegramClient, shopRepo)
	statusService := services.NewStatusService(integrationRepo, telegramSendLogRepo)

	// Инициализация хендлеров
	telegramHandler := handlers.NewTelegramHandler(telegramIntegrationService, statusService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Настройка роутера
	router := handlers.SetupRouter(cfg, telegramHandler, orderHandler)

	// Запуск сервера
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

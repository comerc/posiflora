package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/posiflora/backend/internal/models"
	"github.com/uptrace/bun"
)

type TelegramSendLogRepository interface {
	GetByShopIDAndOrderID(ctx context.Context, shopID, orderID int64) (*models.TelegramSendLog, error)
	Create(ctx context.Context, log *models.TelegramSendLog) error
	GetStatsForLast7Days(ctx context.Context, shopID int64) (*TelegramStats, error)
}

type TelegramStats struct {
	SentCount   int64      `json:"sent_count"`
	FailedCount int64      `json:"failed_count"`
	LastSentAt  *time.Time `json:"last_sent_at"`
}

type telegramSendLogRepository struct {
	db *bun.DB
}

func NewTelegramSendLogRepository(db *bun.DB) TelegramSendLogRepository {
	return &telegramSendLogRepository{db: db}
}

func (r *telegramSendLogRepository) GetByShopIDAndOrderID(ctx context.Context, shopID, orderID int64) (*models.TelegramSendLog, error) {
	var log models.TelegramSendLog
	err := r.db.NewSelect().
		Model(&log).
		Where("shop_id = ? AND order_id = ?", shopID, orderID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (r *telegramSendLogRepository) Create(ctx context.Context, log *models.TelegramSendLog) error {
	_, err := r.db.NewInsert().
		Model(log).
		Exec(ctx)
	return err
}

func (r *telegramSendLogRepository) GetStatsForLast7Days(ctx context.Context, shopID int64) (*TelegramStats, error) {
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	var stats TelegramStats

	// Подсчет отправленных
	sentCount, err := r.db.NewSelect().
		Model((*models.TelegramSendLog)(nil)).
		Where("shop_id = ? AND status = ? AND sent_at >= ?", shopID, models.StatusSent, sevenDaysAgo).
		Count(ctx)
	if err != nil {
		return nil, err
	}
	stats.SentCount = int64(sentCount)

	// Подсчет неудачных
	failedCount, err := r.db.NewSelect().
		Model((*models.TelegramSendLog)(nil)).
		Where("shop_id = ? AND status = ? AND sent_at >= ?", shopID, models.StatusFailed, sevenDaysAgo).
		Count(ctx)
	if err != nil {
		return nil, err
	}
	stats.FailedCount = int64(failedCount)

	// Последняя отправка
	var lastLog models.TelegramSendLog
	err = r.db.NewSelect().
		Model(&lastLog).
		Where("shop_id = ? AND sent_at >= ?", shopID, sevenDaysAgo).
		Order("sent_at DESC").
		Limit(1).
		Scan(ctx)
	if err == nil {
		stats.LastSentAt = &lastLog.SentAt
	}

	return &stats, nil
}

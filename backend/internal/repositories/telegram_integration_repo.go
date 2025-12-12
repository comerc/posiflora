package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/posiflora/backend/internal/models"
	"github.com/uptrace/bun"
)

type TelegramIntegrationRepository interface {
	GetByShopID(ctx context.Context, shopID int64) (*models.TelegramIntegration, error)
	Upsert(ctx context.Context, integration *models.TelegramIntegration) error
}

type telegramIntegrationRepository struct {
	db *bun.DB
}

func NewTelegramIntegrationRepository(db *bun.DB) TelegramIntegrationRepository {
	return &telegramIntegrationRepository{db: db}
}

func (r *telegramIntegrationRepository) GetByShopID(ctx context.Context, shopID int64) (*models.TelegramIntegration, error) {
	var integration models.TelegramIntegration
	err := r.db.NewSelect().
		Model(&integration).
		Where("shop_id = ?", shopID).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &integration, nil
}

func (r *telegramIntegrationRepository) Upsert(ctx context.Context, integration *models.TelegramIntegration) error {
	existing, err := r.GetByShopID(ctx, integration.ShopID)
	if err != nil {
		return err
	}

	if existing != nil {
		// Обновляем существующую запись
		integration.ID = existing.ID
		integration.CreatedAt = existing.CreatedAt
		_, err = r.db.NewUpdate().
			Model(integration).
			Where("id = ?", existing.ID).
			Exec(ctx)
		return err
	}

	// Создаем новую запись
	_, err = r.db.NewInsert().
		Model(integration).
		Exec(ctx)
	return err
}

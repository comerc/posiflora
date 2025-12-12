package repositories

import (
	"context"
	"fmt"

	"github.com/posiflora/backend/internal/models"
	"github.com/uptrace/bun"
)

type ShopRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Shop, error)
	Create(ctx context.Context, shop *models.Shop) error
	GetOrCreate(ctx context.Context, id int64) (*models.Shop, error)
}

type shopRepository struct {
	db *bun.DB
}

func NewShopRepository(db *bun.DB) ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) GetByID(ctx context.Context, id int64) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.NewSelect().
		Model(&shop).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepository) Create(ctx context.Context, shop *models.Shop) error {
	_, err := r.db.NewInsert().
		Model(shop).
		Exec(ctx)
	return err
}

func (r *shopRepository) GetOrCreate(ctx context.Context, id int64) (*models.Shop, error) {
	shop, err := r.GetByID(ctx, id)
	if err == nil && shop != nil {
		return shop, nil
	}

	// Магазин не найден, создаем новый
	shop = &models.Shop{
		ID:   id,
		Name: fmt.Sprintf("Магазин #%d", id),
	}
	err = r.Create(ctx, shop)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

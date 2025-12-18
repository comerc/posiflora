package repositories

import (
	"context"

	"github.com/posiflora/backend/internal/models"
	"github.com/uptrace/bun"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id int64) (*models.Order, error)
	GetByShopIDAndNumber(ctx context.Context, shopID int64, number string) (*models.Order, error)
}

type orderRepository struct {
	db *bun.DB
}

func NewOrderRepository(db *bun.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	_, err := r.db.NewInsert().
		Model(order).
		Returning("*").
		Exec(ctx)
	return err
}

func (r *orderRepository) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	var order models.Order
	err := r.db.NewSelect().
		Model(&order).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByShopIDAndNumber(ctx context.Context, shopID int64, number string) (*models.Order, error) {
	var order models.Order
	err := r.db.NewSelect().
		Model(&order).
		Where("shop_id = ? AND number = ?", shopID, number).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

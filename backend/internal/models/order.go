package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders"`

	ID           int64     `bun:"id,pk,autoincrement" json:"id"`
	ShopID       int64     `bun:"shop_id,notnull" json:"shop_id"`
	Number       string    `bun:"number,notnull" json:"number"`
	Total        float64   `bun:"total,notnull" json:"total"`
	CustomerName string    `bun:"customer_name,notnull" json:"customer_name"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
}

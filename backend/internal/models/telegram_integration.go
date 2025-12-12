package models

import (
	"time"

	"github.com/uptrace/bun"
)

type TelegramIntegration struct {
	bun.BaseModel `bun:"table:telegram_integrations"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	ShopID    int64     `bun:"shop_id,notnull,unique" json:"shop_id"`
	BotToken  string    `bun:"bot_token,notnull" json:"bot_token"`
	ChatID    string    `bun:"chat_id,notnull" json:"chat_id"`
	Enabled   bool      `bun:"enabled,notnull,default:false" json:"enabled"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

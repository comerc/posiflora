package models

import (
	"time"

	"github.com/uptrace/bun"
)

type TelegramSendLogStatus string

const (
	StatusSent   TelegramSendLogStatus = "SENT"
	StatusFailed TelegramSendLogStatus = "FAILED"
)

type TelegramSendLog struct {
	bun.BaseModel `bun:"table:telegram_send_log"`

	ID      int64                 `bun:"id,pk,autoincrement" json:"id"`
	ShopID  int64                 `bun:"shop_id,notnull" json:"shop_id"`
	OrderID int64                 `bun:"order_id,notnull" json:"order_id"`
	Message string                `bun:"message,notnull" json:"message"`
	Status  TelegramSendLogStatus `bun:"status,notnull" json:"status"`
	Error   *string               `bun:"error" json:"error,omitempty"`
	SentAt  time.Time             `bun:"sent_at,notnull,default:current_timestamp" json:"sent_at"`
}

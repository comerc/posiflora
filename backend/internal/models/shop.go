package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Shop struct {
	bun.BaseModel `bun:"table:shops"`

	ID   int64  `bun:"id,pk,autoincrement" json:"id"`
	Name string `bun:"name,notnull" json:"name"`

	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

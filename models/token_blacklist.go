package models

import (
	"time"

	"github.com/uptrace/bun"
)

type TokenBlacklist struct {
	bun.BaseModel `bun:"table:token_blacklist,alias:tb"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Token     string    `bun:"token,unique,notnull" json:"token"`
	ExpiresAt time.Time `bun:"expires_at,notnull" json:"expires_at"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}

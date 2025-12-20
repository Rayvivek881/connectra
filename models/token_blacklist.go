package models

import (
	"time"

	"github.com/uptrace/bun"
)

// TokenBlacklist stores blacklisted refresh tokens
type TokenBlacklist struct {
	bun.BaseModel `bun:"table:token_blacklist,alias:tb"`

	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	Token     string     `bun:"token,unique,notnull,type:text" json:"token"`
	UserID    *string    `bun:"user_id,type:text" json:"user_id,omitempty"` // For audit trail
	ExpiresAt time.Time  `bun:"expires_at,notnull,type:timestamptz" json:"expires_at"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
}

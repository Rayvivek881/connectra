package models

import (
	"time"

	"github.com/uptrace/bun"
)

// FeatureUsage tracks feature usage per user and billing period
type FeatureUsage struct {
	bun.BaseModel `bun:"table:feature_usage,alias:fu"`

	ID          int64       `bun:"id,pk,autoincrement" json:"id"`
	UserID      string      `bun:"user_id,notnull,type:text" json:"user_id"`
	Feature     FeatureType `bun:"feature,notnull,type:varchar(50)" json:"feature"`
	Used        int         `bun:"used,notnull,default:0" json:"used"`
	Limit       int         `bun:"limit,notnull,default:0" json:"limit"` // "limit" is a reserved keyword but works with quotes
	PeriodStart time.Time   `bun:"period_start,notnull,default:current_timestamp,type:timestamptz" json:"period_start"`
	PeriodEnd   *time.Time  `bun:"period_end,type:timestamptz" json:"period_end,omitempty"`
	CreatedAt   time.Time   `bun:"created_at,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt   *time.Time  `bun:"updated_at,type:timestamptz" json:"updated_at,omitempty"`

	// Relationship
	User *User `bun:"rel:belongs-to,join:user_id=uuid" json:"user,omitempty"`
}

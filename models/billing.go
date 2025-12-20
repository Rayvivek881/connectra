package models

import (
	"time"

	"github.com/uptrace/bun"
)

// SubscriptionPlan represents a subscription plan tier
type SubscriptionPlan struct {
	bun.BaseModel `bun:"table:subscription_plans,alias:sp"`

	Tier      string    `bun:"tier,pk,type:varchar(50)" json:"tier"`
	Name      string    `bun:"name,notnull,type:varchar(255)" json:"name"`
	Category  string    `bun:"category,notnull,type:varchar(50)" json:"category"` // STARTER, PROFESSIONAL, BUSINESS, ENTERPRISE
	IsActive  bool      `bun:"is_active,notnull,default:true" json:"is_active"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at,type:timestamptz" json:"updated_at,omitempty"`

	// Relationships
	Periods []SubscriptionPlanPeriod `bun:"rel:has-many,join:tier=plan_tier" json:"periods,omitempty"`
}

// SubscriptionPlanPeriod represents pricing for a specific billing period
type SubscriptionPlanPeriod struct {
	bun.BaseModel `bun:"table:subscription_plan_periods,alias:spp"`

	ID                int        `bun:"id,pk,autoincrement" json:"id"`
	PlanTier          string     `bun:"plan_tier,notnull,type:varchar(50)" json:"plan_tier"`
	Period            string     `bun:"period,notnull,type:varchar(20)" json:"period"` // monthly, quarterly, yearly
	Credits           int        `bun:"credits,notnull" json:"credits"`
	RatePerCredit     float64    `bun:"rate_per_credit,notnull,type:numeric(10,6)" json:"rate_per_credit"`
	Price             float64    `bun:"price,notnull,type:numeric(10,2)" json:"price"`
	SavingsAmount     *float64   `bun:"savings_amount,type:numeric(10,2)" json:"savings_amount,omitempty"`
	SavingsPercentage *int       `bun:"savings_percentage" json:"savings_percentage,omitempty"`
	CreatedAt         time.Time  `bun:"created_at,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt         *time.Time `bun:"updated_at,type:timestamptz" json:"updated_at,omitempty"`

	// Relationships
	Plan *SubscriptionPlan `bun:"rel:belongs-to,join:plan_tier=tier" json:"plan,omitempty"`
}

// AddonPackage represents an addon credit package
type AddonPackage struct {
	bun.BaseModel `bun:"table:addon_packages,alias:ap"`

	ID            string     `bun:"id,pk,type:varchar(50)" json:"id"`
	Name          string     `bun:"name,notnull,type:varchar(255)" json:"name"`
	Credits       int        `bun:"credits,notnull" json:"credits"`
	RatePerCredit float64    `bun:"rate_per_credit,notnull,type:numeric(10,6)" json:"rate_per_credit"`
	Price         float64    `bun:"price,notnull,type:numeric(10,2)" json:"price"`
	IsActive      bool       `bun:"is_active,notnull,default:true" json:"is_active"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt     *time.Time `bun:"updated_at,type:timestamptz" json:"updated_at,omitempty"`
}


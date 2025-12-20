package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

// NotificationPreferences represents user notification settings
type NotificationPreferences struct {
	WeeklyReports  *bool `json:"weeklyReports,omitempty"`
	NewLeadAlerts  *bool `json:"newLeadAlerts,omitempty"`
}

// Value implements driver.Valuer for JSONB storage
func (n NotificationPreferences) Value() (driver.Value, error) {
	return json.Marshal(n)
}

// Scan implements sql.Scanner for JSONB retrieval
func (n *NotificationPreferences) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, n)
}

// UserProfile represents additional user information and preferences
type UserProfile struct {
	bun.BaseModel `bun:"table:user_profiles,alias:up"`

	ID                   int64                  `bun:"id,pk,autoincrement" json:"id"`
	UserID               string                 `bun:"user_id,unique,notnull,type:text" json:"user_id"`
	JobTitle             *string                `bun:"job_title,type:varchar(255)" json:"job_title,omitempty"`
	Bio                  *string                `bun:"bio,type:text" json:"bio,omitempty"`
	Timezone             *string                `bun:"timezone,type:varchar(100)" json:"timezone,omitempty"`
	AvatarURL            *string                `bun:"avatar_url,type:text" json:"avatar_url,omitempty"`
	Notifications        *NotificationPreferences `bun:"notifications,type:jsonb" json:"notifications,omitempty"`
	Role                 *string                `bun:"role,type:varchar(50),default:'Member'" json:"role,omitempty"`
	
	// Billing fields
	Credits              int                    `bun:"credits,notnull,default:0" json:"credits"`
	SubscriptionPlan     *string                `bun:"subscription_plan,type:varchar(50),default:'free'" json:"subscription_plan,omitempty"`
	SubscriptionPeriod   *string                `bun:"subscription_period,type:varchar(20),default:'monthly'" json:"subscription_period,omitempty"`
	SubscriptionStatus   *string                `bun:"subscription_status,type:varchar(50),default:'active'" json:"subscription_status,omitempty"`
	SubscriptionStartedAt *time.Time            `bun:"subscription_started_at,type:timestamptz" json:"subscription_started_at,omitempty"`
	SubscriptionEndsAt   *time.Time            `bun:"subscription_ends_at,type:timestamptz" json:"subscription_ends_at,omitempty"`
	
	CreatedAt            time.Time              `bun:"created_at,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt            *time.Time             `bun:"updated_at,type:timestamptz" json:"updated_at,omitempty"`

	// Relationship
	User *User `bun:"rel:belongs-to,join:user_id=uuid" json:"user,omitempty"`
}


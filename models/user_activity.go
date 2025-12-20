package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

// JSONBMap is a helper type for JSONB fields
type JSONBMap map[string]interface{}

// Value implements driver.Valuer for JSONB storage
func (j JSONBMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner for JSONB retrieval
func (j *JSONBMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// UserActivity tracks LinkedIn and email service activities
type UserActivity struct {
	bun.BaseModel `bun:"table:user_activities,alias:ua"`

	ID            int64               `bun:"id,pk,autoincrement" json:"id"`
	UserID        string              `bun:"user_id,notnull,type:text" json:"user_id"`
	ServiceType   ActivityServiceType `bun:"service_type,notnull,type:varchar(50)" json:"service_type"`
	ActionType    ActivityActionType  `bun:"action_type,notnull,type:varchar(50)" json:"action_type"`
	Status        ActivityStatus      `bun:"status,notnull,type:varchar(50)" json:"status"`
	RequestParams *JSONBMap           `bun:"request_params,type:jsonb" json:"request_params,omitempty"`
	ResultCount   int                 `bun:"result_count,notnull,default:0" json:"result_count"`
	ResultSummary *JSONBMap           `bun:"result_summary,type:jsonb" json:"result_summary,omitempty"`
	ErrorMessage  *string             `bun:"error_message,type:text" json:"error_message,omitempty"`
	IPAddress     *string             `bun:"ip_address,type:varchar(45)" json:"ip_address,omitempty"` // IPv6 max length
	UserAgent     *string             `bun:"user_agent,type:text" json:"user_agent,omitempty"`
	CreatedAt     time.Time           `bun:"created_at,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`

	// Relationship
	User *User `bun:"rel:belongs-to,join:user_id=uuid" json:"user,omitempty"`
}

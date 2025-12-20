package models

import (
	"time"

	"github.com/uptrace/bun"
)

// User represents the core user authentication model
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID            string     `bun:"id,pk,type:text" json:"id"`
	UUID          string     `bun:"uuid,unique,notnull,type:text" json:"uuid"`
	Email         string     `bun:"email,unique,notnull,type:varchar(255)" json:"email"`
	HashedPassword string   `bun:"hashed_password,notnull,type:text" json:"-"`
	Name          *string    `bun:"name,type:varchar(255)" json:"name,omitempty"`
	IsActive      bool       `bun:"is_active,notnull,default:true" json:"is_active"`
	LastSignInAt  *time.Time `bun:"last_sign_in_at,type:timestamptz" json:"last_sign_in_at,omitempty"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt     *time.Time `bun:"updated_at,type:timestamptz" json:"updated_at,omitempty"`

	// Relationships
	Profile  *UserProfile    `bun:"rel:has-one,join:uuid=user_id" json:"profile,omitempty"`
	History  []UserHistory    `bun:"rel:has-many,join:uuid=user_id" json:"history,omitempty"`
	Activities []UserActivity `bun:"rel:has-many,join:uuid=user_id" json:"activities,omitempty"`
}

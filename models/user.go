package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Role string

const (
	RoleSuperAdmin Role = "SuperAdmin"
	RoleAdmin      Role = "Admin"
	RoleProUser    Role = "ProUser"
	RoleFreeUser   Role = "FreeUser"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	Password  string    `bun:"password,notnull" json:"-"` // Password should not be exposed in JSON
	Role      Role      `bun:"role,notnull" json:"role"`
	Credits   int       `bun:"credits,default:0" json:"credits"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

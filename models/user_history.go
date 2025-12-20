package models

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

// UserHistory tracks user registration and login events with IP geolocation data
type UserHistory struct {
	bun.BaseModel `bun:"table:user_history,alias:uh"`

	ID            int64                  `bun:"id,pk,autoincrement" json:"id"`
	UserID        string                 `bun:"user_id,notnull,type:text" json:"user_id"`
	EventType     UserHistoryEventType   `bun:"event_type,notnull,type:varchar(50)" json:"event_type"`
	
	// IP and device info
	IP            *string                `bun:"ip,type:varchar(45)" json:"ip,omitempty"` // IPv6 max length
	Device        *string                `bun:"device,type:text" json:"device,omitempty"` // User-Agent string
	
	// Geolocation fields
	Continent     *string                `bun:"continent,type:varchar(50)" json:"continent,omitempty"`
	ContinentCode *string                `bun:"continent_code,type:varchar(2)" json:"continent_code,omitempty"`
	Country       *string                `bun:"country,type:varchar(100)" json:"country,omitempty"`
	CountryCode   *string                `bun:"country_code,type:varchar(2)" json:"country_code,omitempty"`
	Region        *string                `bun:"region,type:varchar(10)" json:"region,omitempty"`
	RegionName    *string                `bun:"region_name,type:varchar(100)" json:"region_name,omitempty"`
	City          *string                `bun:"city,type:varchar(100)" json:"city,omitempty"`
	District      *string                `bun:"district,type:varchar(100)" json:"district,omitempty"`
	Zip           *string                `bun:"zip,type:varchar(20)" json:"zip,omitempty"`
	Lat           *sql.NullFloat64       `bun:"lat,type:numeric(10,7)" json:"lat,omitempty"`
	Lon           *sql.NullFloat64       `bun:"lon,type:numeric(10,7)" json:"lon,omitempty"`
	Timezone      *string                `bun:"timezone,type:varchar(100)" json:"timezone,omitempty"`
	Currency      *string                `bun:"currency,type:varchar(10)" json:"currency,omitempty"`
	ISP           *string                `bun:"isp,type:varchar(255)" json:"isp,omitempty"`
	Org           *string                `bun:"org,type:varchar(255)" json:"org,omitempty"`
	ASName        *string                `bun:"asname,type:varchar(255)" json:"asname,omitempty"`
	Reverse       *string                `bun:"reverse,type:varchar(255)" json:"reverse,omitempty"`
	Proxy         *bool                  `bun:"proxy,default:false" json:"proxy,omitempty"`
	Hosting       *bool                  `bun:"hosting,default:false" json:"hosting,omitempty"`
	
	CreatedAt     time.Time              `bun:"created_at,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`

	// Relationship
	User          *User                  `bun:"rel:belongs-to,join:user_id=uuid" json:"user,omitempty"`
}


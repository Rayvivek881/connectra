package models

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type ModelJobsData struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:jobs_data,alias:jd"`

	Id      uint64          `bun:"id,pk,autoincrement" json:"id"`
	UUID    string          `bun:"uuid,notnull,unique" json:"uuid"`
	JobType string          `bun:"job_type,notnull" json:"job_type"`
	Data    json.RawMessage `bun:"data,type:jsonb" json:"data"`
	Status  string          `bun:"status,notnull" json:"status"`

	RetryAfter       *time.Time `bun:"retry_after,nullzero" json:"retry_after"`
	RetryInterval    int        `bun:"retry_interval,notnull,default:60" json:"retry_interval"`
	RemainingRetries int        `bun:"remaining_retries,notnull,default:3" json:"remaining_retries"`
	RuntimeErrors    []string   `bun:"runtime_errors,array" json:"runtime_errors"`

	CreatedAt time.Time `bun:"created_at,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,default:current_timestamp" json:"updated_at"`
}

func (m *ModelJobsData) SetDB(db *bun.DB) *ModelJobsData {
	m.db = db
	return m
}

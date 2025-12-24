package models

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type ModelJobs struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:jobs,alias:j"`

	Id      uint64          `bun:"id,pk,autoincrement" json:"id"`
	UUID    string          `bun:"uuid,notnull,unique" json:"uuid"`
	JobType string          `bun:"job_type,notnull" json:"job_type"`
	Data    json.RawMessage `bun:"data,type:jsonb" json:"data"`
	Status  string          `bun:"status,notnull" json:"status"`

	RetryCount    int        `bun:"retry_count,notnull,default:0" json:"retry_count"`
	RuntimeErrors []string   `bun:"runtime_errors,type:text[]" json:"runtime_errors"`
	RunAfter      *time.Time `bun:"run_after,nullzero" json:"run_after"`

	CreatedAt *time.Time `bun:"created_at,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at,nullzero,default:current_timestamp" json:"updated_at"`
}

func (m *ModelJobs) SetDB(db *bun.DB) *ModelJobs {
	m.db = db
	return m
}

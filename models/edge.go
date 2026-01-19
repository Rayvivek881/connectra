package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ModelEdges struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:edge,alias:e"`

	Id     uint64 `bun:"id,pk,autoincrement" json:"id"`
	Source string `bun:"source,notnull" json:"source"`
	Target string `bun:"target,notnull" json:"target"`

	CreatedAt *time.Time `bun:"created_at,nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at,nullzero,default:current_timestamp" json:"updated_at"`
}

func (m *ModelEdges) SetDB(db *bun.DB) *ModelEdges {
	m.db = db
	return m
}

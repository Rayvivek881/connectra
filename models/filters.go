package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ModelFilter struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:filters,alias:cf"`

	Id           uint64     `bun:"id,pk,autoincrement" json:"id"`
	Key          string     `bun:"key" json:"key"`
	Service      string     `bun:"service" json:"service"`
	DisplayName  string     `bun:"display_name,notnull" json:"display_name"`
	DirectDrived bool       `bun:"direct_drived,nullzero" json:"direct_drived"`
	DeletedAt    *time.Time `bun:"deleted_at,nullzero" json:"deleted_at"`
}

func (m *ModelFilter) SetDB(db *bun.DB) *ModelFilter {
	m.db = db
	return m
}

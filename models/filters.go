package models

import (
	"github.com/uptrace/bun"
)

type ModelFilter struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:filters,alias:cf"`

	Id           uint64 `bun:"id,pk,autoincrement" json:"id"`
	Key          string `bun:"key" json:"key"`
	Service      string `bun:"service" json:"service"`
	FilterType   string `bun:"filter_type,notnull" json:"filter_type"`
	DisplayName  string `bun:"display_name,notnull" json:"display_name"`
	DirectDerived bool   `bun:"direct_derived,nullzero" json:"direct_derived"`
	Active       bool   `bun:"active,notnull,default:true" json:"active"`

	DeletedAt bun.NullTime `bun:"deleted_at,nullzero" json:"deleted_at,omitempty"`
}

func (m *ModelFilter) SetDB(db *bun.DB) *ModelFilter {
	m.db = db
	return m
}

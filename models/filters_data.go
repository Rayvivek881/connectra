package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ModelFilterData struct {
	db            *bun.DB
	bun.BaseModel `bun:"table:filters_data,alias:cfd"`

	Id   uint64 `bun:"id,pk,autoincrement" json:"id"`
	UUID string `bun:"uuid,notnull,unique" json:"uuid"` // filter_key + service + value

	FilterKey    string     `bun:"filter_key,notnull" json:"filter_key"`
	Service      string     `bun:"service,notnull" json:"service"`
	DisplayValue string     `bun:"display_value,notnull" json:"display_value"`
	Value        string     `bun:"value,nullzero" json:"value"`
	DeletedAt    *time.Time `bun:"deleted_at,nullzero" json:"deleted_at"`
}

func (m *ModelFilterData) SetDB(db *bun.DB) *ModelFilterData {
	m.db = db
	return m
}

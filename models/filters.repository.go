package models

import (
	"context"

	"github.com/uptrace/bun"
)

type FiltersStruct struct {
	PgDbClient *bun.DB
}

func FiltersRepository(db *bun.DB) FiltersSvcRepo {
	return &FiltersStruct{
		PgDbClient: db,
	}
}

type FiltersSvcRepo interface {
	GetFiltersByService(service string) ([]*ModelFilter, error)
	GetFilterByKeyAndService(service, key string) (ModelFilter, error)
}

func (t *FiltersStruct) GetFiltersByService(service string) ([]*ModelFilter, error) {
	var filters []*ModelFilter
	err := t.PgDbClient.NewSelect().Model(&filters).Where("service = ?", service).Scan(context.Background())
	return filters, err
}

func (t *FiltersStruct) GetFilterByKeyAndService(service, key string) (ModelFilter, error) {
	var filter ModelFilter
	err := t.PgDbClient.NewSelect().Model(&filter).Where("service = ? AND key = ?", service, key).Scan(context.Background())

	return filter, err
}

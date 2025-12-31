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
	GetTempFilters() ([]*ModelFilter, error)
	GetFiltersByService(service string) ([]*ModelFilter, error)
	GetFilterByKeyAndService(service, key string) (ModelFilter, error)
	UpdateActiveStatus(key, service string, status bool) error
}

func (t *FiltersStruct) GetTempFilters() ([]*ModelFilter, error) {
	var filters []*ModelFilter
	err := t.PgDbClient.NewSelect().Model(&filters).Where("direct_derived = false").Scan(context.Background())
	return filters, err
}

func (t *FiltersStruct) GetFiltersByService(service string) ([]*ModelFilter, error) {
	var filters []*ModelFilter
	err := t.PgDbClient.NewSelect().Model(&filters).Where("active = true AND deleted_at IS NULL").
		Where("service = ?", service).Scan(context.Background())
	return filters, err
}

func (t *FiltersStruct) GetFilterByKeyAndService(service, key string) (ModelFilter, error) {
	var filter ModelFilter
	err := t.PgDbClient.NewSelect().Model(&filter).Where("active = true AND deleted_at IS NULL").
		Where("service = ? AND key = ?", service, key).Scan(context.Background())

	return filter, err
}

func (t *FiltersStruct) UpdateActiveStatus(key, service string, status bool) error {
	_, err := t.PgDbClient.NewUpdate().Model(&ModelFilter{}).
		Set("active = ?", status).
		Where("key = ? AND service = ?", key, service).
		Exec(context.Background())
	return err
}


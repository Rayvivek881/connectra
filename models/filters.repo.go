package models

import (
	"context"
	"time"
	"vivek-ray/utilities"

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
	// Try to get from cache
	cacheKey := "temp_filters:all"
	if cached, found := utilities.FilterMetadataCache.Get(cacheKey); found {
		if filters, ok := cached.([]*ModelFilter); ok {
			return filters, nil
		}
	}
	
	// Fetch from database
	var filters []*ModelFilter
	err := t.PgDbClient.NewSelect().Model(&filters).Where("direct_derived = false").Scan(context.Background())
	if err != nil {
		return nil, err
	}
	
	// Store in cache with longer TTL (temp filters change infrequently)
	utilities.FilterMetadataCache.SetWithTTL(cacheKey, filters, 15*time.Minute)
	
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


package models

import (
	"context"
	"vivek-ray/constants"
	"vivek-ray/utilities"

	"github.com/uptrace/bun"
)

type FiltersDataQuery struct {
	Service    string `json:"service,omitempty"`
	FilterKey  string `json:"filter_key"`
	SearchText string `json:"search_text,omitempty"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type FiltersDataStruct struct {
	PgDbClient *bun.DB
}

func FiltersDataRepository(db *bun.DB) FiltersDataSvcRepo {
	return &FiltersDataStruct{
		PgDbClient: db,
	}
}

type FiltersDataSvcRepo interface {
	GetFiltersByQuery(query FiltersDataQuery) ([]*ModelFilterData, error)
	BulkUpsert(filtersData []*ModelFilterData) error
}

func (t *FiltersDataStruct) GetFiltersByQuery(query FiltersDataQuery) ([]*ModelFilterData, error) {
	var filtersData []*ModelFilterData

	queryBuilder := t.PgDbClient.NewSelect().Model(&filtersData).Where("service = ?", query.Service).Where("filter_key = ?", query.FilterKey)
	if query.SearchText != "" {
		queryBuilder = queryBuilder.Where("display_value ILIKE ?", "%"+query.SearchText+"%")
	}
	query.Limit = utilities.InlineIf(query.Limit > 0, query.Limit, constants.DefaultPageSize).(int)
	if query.Page > 0 {
		queryBuilder = queryBuilder.Offset((query.Page - 1) * query.Limit)
	}
	err := queryBuilder.Limit(query.Limit).Scan(context.Background())
	return filtersData, err
}

func (t *FiltersDataStruct) BulkUpsert(filtersData []*ModelFilterData) error {
	_, err := t.PgDbClient.NewInsert().
		Model(&filtersData).
		On("CONFLICT(uuid) DO NOTHING").
		Exec(context.Background())
	
	// Invalidate cache for affected filter keys
	if err == nil && len(filtersData) > 0 {
		// Group by service and filter_key to minimize cache invalidation
		services := make(map[string]bool)
		for _, fd := range filtersData {
			if fd.Service != "" {
				services[fd.Service] = true
			}
		}
		
		// Invalidate all filter data cache for affected services
		// (More granular invalidation would require tracking specific filter_key)
		for service := range services {
			utilities.FilterDataCache.InvalidateByPrefix("filter_data:" + service + ":")
		}
	}
	
	return err
}


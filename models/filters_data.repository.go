package models

import (
	"context"

	"github.com/uptrace/bun"
)

type FiltersDataQuery struct {
	Service    string `json:"service"`
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
}

func (t *FiltersDataStruct) GetFiltersByQuery(query FiltersDataQuery) ([]*ModelFilterData, error) {
	var filtersData []*ModelFilterData

	queryBuilder := t.PgDbClient.NewSelect().Model(&filtersData).Where("service = ?", query.Service).Where("filter_key = ?", query.FilterKey)
	if query.SearchText != "" {
		queryBuilder = queryBuilder.Where("display_value ILIKE ?", "%"+query.SearchText+"%")
	}
	if query.Page > 0 && query.Limit > 0 {
		queryBuilder = queryBuilder.Offset((query.Page - 1) * query.Limit).Limit(query.Limit)
	}
	err := queryBuilder.Scan(context.Background())
	return filtersData, err
}

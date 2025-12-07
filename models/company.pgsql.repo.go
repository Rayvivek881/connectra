package models

import (
	"context"

	"github.com/uptrace/bun"
)

type PgCompanyStruct struct {
	PgDbClient *bun.DB
}

func PgCompanyRepository(db *bun.DB) PgCompanySvcRepo {
	return &PgCompanyStruct{
		PgDbClient: db,
	}
}

type PgCompanyFilters struct {
	Uuids             []string
	Names             []string
	NormalizedDomains []string

	Page  int
	Limit int
}

func (f *PgCompanyFilters) ToQuery(query *bun.SelectQuery) *bun.SelectQuery {
	if len(f.Uuids) > 0 {
		query.Where("uuid IN (?)", bun.In(f.Uuids))
	}
	if len(f.Names) > 0 {
		query.Where("name IN (?)", bun.In(f.Names))
	}
	if len(f.NormalizedDomains) > 0 {
		query.Where("normalized_domain IN (?)", bun.In(f.NormalizedDomains))
	}
	return query
}

type PgCompanySvcRepo interface {
	GetFiltersByQuery(query FiltersDataQuery) ([]*PgCompany, error)
	ListByFilters(filters PgCompanyFilters) ([]*PgCompany, error)
}

func (t *PgCompanyStruct) GetFiltersByQuery(query FiltersDataQuery) ([]*PgCompany, error) {
	var companies []*PgCompany

	queryBuilder := t.PgDbClient.NewSelect().Model(&companies).Where("? ILIKE ?", query.FilterKey, "%"+query.SearchText+"%")
	if query.Page > 0 && query.Limit > 0 {
		queryBuilder = queryBuilder.Offset((query.Page - 1) * query.Limit).Limit(query.Limit)
	}
	err := queryBuilder.Scan(context.Background())
	return companies, err
}

func (t *PgCompanyStruct) ListByFilters(filters PgCompanyFilters) ([]*PgCompany, error) {
	var companies []*PgCompany

	queryBuilder := t.PgDbClient.NewSelect().Model(&companies)
	filters.ToQuery(queryBuilder)

	if filters.Page > 0 && filters.Limit > 0 {
		queryBuilder = queryBuilder.Offset((filters.Page - 1) * filters.Limit).Limit(filters.Limit)
	}
	err := queryBuilder.Scan(context.Background())
	return companies, err
}

package models

import (
	"context"

	"github.com/uptrace/bun"
)

type PgContactStruct struct {
	PgDbClient *bun.DB
}

func PgContactRepository(db *bun.DB) PgContactSvcRepo {
	return &PgContactStruct{
		PgDbClient: db,
	}
}

type PgContactFilters struct {
	Uuids        []string
	CompanyIds   []string
	Emails       []string
	MobilePhones []string

	Page  int
	Limit int
}

func (f *PgContactFilters) ToQuery(query *bun.SelectQuery) *bun.SelectQuery {
	if len(f.Uuids) > 0 {
		query.Where("uuid IN (?)", bun.In(f.Uuids))
	}
	if len(f.CompanyIds) > 0 {
		query.Where("company_id IN (?)", bun.In(f.CompanyIds))
	}
	if len(f.Emails) > 0 {
		query.Where("email IN (?)", bun.In(f.Emails))
	}
	if len(f.MobilePhones) > 0 {
		query.Where("mobile_phone IN (?)", bun.In(f.MobilePhones))
	}
	return query
}

type PgContactSvcRepo interface {
	GetFiltersByQuery(query FiltersDataQuery) ([]*PgContact, error)
	ListByFilters(filters PgContactFilters) ([]*PgContact, error)
}

func (t *PgContactStruct) GetFiltersByQuery(query FiltersDataQuery) ([]*PgContact, error) {
	var contacts []*PgContact

	queryBuilder := t.PgDbClient.NewSelect().Model(&contacts).Where("? ILIKE ?", query.FilterKey, "%"+query.SearchText+"%")
	if query.Page > 0 && query.Limit > 0 {
		queryBuilder = queryBuilder.Offset((query.Page - 1) * query.Limit).Limit(query.Limit)
	}
	err := queryBuilder.Scan(context.Background())
	return contacts, err
}

func (t *PgContactStruct) ListByFilters(filters PgContactFilters) ([]*PgContact, error) {
	var contacts []*PgContact

	queryBuilder := t.PgDbClient.NewSelect().Model(&contacts)
	filters.ToQuery(queryBuilder)

	if filters.Page > 0 && filters.Limit > 0 {
		queryBuilder = queryBuilder.Offset((filters.Page - 1) * filters.Limit).Limit(filters.Limit)
	}
	err := queryBuilder.Scan(context.Background())
	return contacts, err
}

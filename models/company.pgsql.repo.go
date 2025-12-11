package models

import (
	"context"
	"vivek-ray/constants"
	"vivek-ray/utilities"

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

func (f *PgCompanyFilters) IsEmpty() bool {
	if len(f.Uuids) > 0 {
		return false
	}
	if len(f.Names) > 0 {
		return false
	}
	if len(f.NormalizedDomains) > 0 {
		return false
	}

	return true
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
	BulkUpsert(companies []*PgCompany) (int64, error)
}

func (t *PgCompanyStruct) GetFiltersByQuery(query FiltersDataQuery) ([]*PgCompany, error) {
	var companies []*PgCompany

	// fetch only filter column
	queryBuilder := t.PgDbClient.NewSelect().Model(&companies).
		Where("? ILIKE ?", bun.Ident(query.FilterKey), "%"+query.SearchText+"%").Column(query.FilterKey).Distinct()

	query.Limit = utilities.InlineIf(query.Limit > 0, query.Limit, constants.DefaultPageSize).(int)
	if query.Page > 0 {
		queryBuilder = queryBuilder.Offset((query.Page - 1) * query.Limit)
	}
	err := queryBuilder.Limit(query.Limit).Scan(context.Background())
	return companies, err
}

func (t *PgCompanyStruct) ListByFilters(filters PgCompanyFilters) ([]*PgCompany, error) {
	companies := make([]*PgCompany, 0)
	if filters.IsEmpty() {
		return companies, nil
	}

	queryBuilder := t.PgDbClient.NewSelect().Model(&companies)
	filters.ToQuery(queryBuilder)

	if filters.Page > 0 && filters.Limit > 0 {
		queryBuilder = queryBuilder.Offset((filters.Page - 1) * filters.Limit).Limit(filters.Limit)
	}
	err := queryBuilder.Scan(context.Background())
	return companies, err
}

func (t *PgCompanyStruct) BulkUpsert(companies []*PgCompany) (int64, error) {
	_, err := t.PgDbClient.NewInsert().
		Model(&companies).
		On("CONFLICT(uuid) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("normalized_domain = EXCLUDED.normalized_domain").
		Set("employees_count = EXCLUDED.employees_count").
		Set("industries = EXCLUDED.industries").
		Set("keywords = EXCLUDED.keywords").
		Set("address = EXCLUDED.address").
		Set("annual_revenue = EXCLUDED.annual_revenue").
		Set("total_funding = EXCLUDED.total_funding").
		Set("technologies = EXCLUDED.technologies").
		Set("city = EXCLUDED.city").
		Set("state = EXCLUDED.state").
		Set("country = EXCLUDED.country").
		Set("linkedin_url = EXCLUDED.linkedin_url").
		Set("website = EXCLUDED.website").
		Set("facebook_url = EXCLUDED.facebook_url").
		Set("twitter_url = EXCLUDED.twitter_url").
		Set("company_name_for_emails = EXCLUDED.company_name_for_emails").
		Set("phone_number = EXCLUDED.phone_number").
		Set("latest_funding = EXCLUDED.latest_funding").
		Set("latest_funding_amount = EXCLUDED.latest_funding_amount").
		Set("last_raised_at = EXCLUDED.last_raised_at").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(context.Background())

	return int64(len(companies)), err
}

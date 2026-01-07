package models

import (
	"context"
	"fmt"
	"time"
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
	SelectColumns     []string

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
		query.Where("uuid IN (?)", bun.In(utilities.UniqueStringSlice(f.Uuids)))
	}
	if len(f.Names) > 0 {
		query.Where("name IN (?)", bun.In(utilities.UniqueStringSlice(f.Names)))
	}
	if len(f.NormalizedDomains) > 0 {
		query.Where("normalized_domain IN (?)", bun.In(utilities.UniqueStringSlice(f.NormalizedDomains)))
	}
	if len(f.SelectColumns) > 0 {
		query.Column(utilities.UniqueStringSlice(f.SelectColumns)...)
	}

	return query
}

type PgCompanySvcRepo interface {
	GetFiltersByQuery(query FiltersDataQuery) ([]*PgCompany, error)
	ListByFilters(filters PgCompanyFilters) ([]*PgCompany, error)
	BulkUpsert(companies []*PgCompany) (int64, error)
	Create(company *PgCompany) error
	GetByUUID(uuid string) (*PgCompany, error)
	GetByUUIDOrDomain(uuid, normalizedDomain string) (*PgCompany, error)
	Update(company *PgCompany) error
	Delete(uuid string) error
}

func (t *PgCompanyStruct) GetFiltersByQuery(query FiltersDataQuery) ([]*PgCompany, error) {
	var companies []*PgCompany

	// fetch only filter column
	queryBuilder := t.PgDbClient.NewSelect().Model(&companies)
	if query.SearchText != "" {
		queryBuilder = queryBuilder.Where("? ILIKE ?", bun.Ident(query.FilterKey), "%"+query.SearchText+"%")
	}

	query.Limit = utilities.InlineIf(query.Limit > 0, query.Limit, constants.DefaultPageSize).(int)
	if query.Page > 0 {
		queryBuilder = queryBuilder.Offset((query.Page - 1) * query.Limit)
	}
	err := queryBuilder.Limit(query.Limit).Column(query.FilterKey).Scan(context.Background())
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

func (t *PgCompanyStruct) Create(company *PgCompany) error {
	_, err := t.PgDbClient.NewInsert().Model(company).Exec(context.Background())
	return err
}

func (t *PgCompanyStruct) GetByUUID(uuid string) (*PgCompany, error) {
	company := new(PgCompany)
	err := t.PgDbClient.NewSelect().Model(company).Where("uuid = ?", uuid).Where("deleted_at IS NULL").Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (t *PgCompanyStruct) GetByUUIDOrDomain(uuid, normalizedDomain string) (*PgCompany, error) {
	company := new(PgCompany)
	query := t.PgDbClient.NewSelect().Model(company).Where("deleted_at IS NULL")
	
	if uuid != "" {
		query = query.Where("uuid = ?", uuid)
	} else if normalizedDomain != "" {
		query = query.Where("normalized_domain = ?", normalizedDomain)
	} else {
		return nil, fmt.Errorf("either UUID or normalized_domain must be provided")
	}
	
	err := query.Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (t *PgCompanyStruct) Update(company *PgCompany) error {
	_, err := t.PgDbClient.NewUpdate().
		Model(company).
		Where("uuid = ?", company.UUID).
		Where("deleted_at IS NULL").
		Column("name").
		Column("normalized_domain").
		Column("employees_count").
		Column("industries").
		Column("keywords").
		Column("address").
		Column("annual_revenue").
		Column("total_funding").
		Column("technologies").
		Column("city").
		Column("state").
		Column("country").
		Column("linkedin_url").
		Column("website").
		Column("facebook_url").
		Column("twitter_url").
		Column("company_name_for_emails").
		Column("phone_number").
		Column("latest_funding").
		Column("latest_funding_amount").
		Column("last_raised_at").
		Column("updated_at").
		Exec(context.Background())
	return err
}

func (t *PgCompanyStruct) Delete(uuid string) error {
	now := time.Now()
	_, err := t.PgDbClient.NewUpdate().
		Model((*PgCompany)(nil)).
		Where("uuid = ?", uuid).
		Where("deleted_at IS NULL").
		Set("deleted_at = ?", now).
		Set("updated_at = ?", now).
		Exec(context.Background())
	return err
}

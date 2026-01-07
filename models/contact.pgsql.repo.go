package models

import (
	"context"
	"fmt"
	"time"
	"vivek-ray/constants"
	"vivek-ray/utilities"

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
	Uuids         []string
	CompanyIds    []string
	Emails        []string
	MobilePhones  []string
	SelectColumns []string

	Page  int
	Limit int
}

func (f *PgContactFilters) IsEmpty() bool {
	if len(f.Uuids) > 0 {
		return false
	}
	if len(f.CompanyIds) > 0 {
		return false
	}
	if len(f.Emails) > 0 {
		return false
	}
	if len(f.MobilePhones) > 0 {
		return false
	}

	return true
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

	if len(f.SelectColumns) > 0 {
		query.Column(utilities.UniqueStringSlice(f.SelectColumns)...)
	}

	return query
}

type PgContactSvcRepo interface {
	GetFiltersByQuery(query FiltersDataQuery) ([]*PgContact, error)
	ListByFilters(filters PgContactFilters) ([]*PgContact, error)
	BulkUpsert(contacts []*PgContact) (int64, error)
	Create(contact *PgContact) error
	GetByUUID(uuid string) (*PgContact, error)
	GetByUUIDOrEmail(uuid, email string) (*PgContact, error)
	Update(contact *PgContact) error
	Delete(uuid string) error
}

func (t *PgContactStruct) GetFiltersByQuery(query FiltersDataQuery) ([]*PgContact, error) {
	var contacts []*PgContact

	queryBuilder := t.PgDbClient.NewSelect().Model(&contacts)

	if query.SearchText != "" {
		queryBuilder = queryBuilder.Where("? ILIKE ?", bun.Ident(query.FilterKey), "%"+query.SearchText+"%")
	}
	query.Limit = utilities.InlineIf(query.Limit > 0, query.Limit, constants.DefaultPageSize).(int)
	if query.Page > 0 {
		queryBuilder = queryBuilder.Offset((query.Page - 1) * query.Limit)
	}
	err := queryBuilder.Limit(query.Limit).Column(query.FilterKey).Scan(context.Background())
	return contacts, err
}

func (t *PgContactStruct) ListByFilters(filters PgContactFilters) ([]*PgContact, error) {
	contacts := make([]*PgContact, 0)
	if filters.IsEmpty() {
		return contacts, nil
	}

	queryBuilder := t.PgDbClient.NewSelect().Model(&contacts)
	filters.ToQuery(queryBuilder)

	if filters.Page > 0 && filters.Limit > 0 {
		queryBuilder = queryBuilder.Offset((filters.Page - 1) * filters.Limit).Limit(filters.Limit)
	}
	err := queryBuilder.Scan(context.Background())
	return contacts, err
}

func (t *PgContactStruct) BulkUpsert(contacts []*PgContact) (int64, error) {
	_, err := t.PgDbClient.NewInsert().
		Model(&contacts).
		On("CONFLICT(uuid) DO UPDATE").
		Set("first_name = EXCLUDED.first_name").
		Set("last_name = EXCLUDED.last_name").
		Set("company_id = EXCLUDED.company_id").
		Set("email = EXCLUDED.email").
		Set("title = EXCLUDED.title").
		Set("departments = EXCLUDED.departments").
		Set("mobile_phone = EXCLUDED.mobile_phone").
		Set("email_status = EXCLUDED.email_status").
		Set("seniority = EXCLUDED.seniority").
		Set("city = EXCLUDED.city").
		Set("state = EXCLUDED.state").
		Set("country = EXCLUDED.country").
		Set("linkedin_url = EXCLUDED.linkedin_url").
		Set("facebook_url = EXCLUDED.facebook_url").
		Set("twitter_url = EXCLUDED.twitter_url").
		Set("website = EXCLUDED.website").
		Set("work_direct_phone = EXCLUDED.work_direct_phone").
		Set("home_phone = EXCLUDED.home_phone").
		Set("other_phone = EXCLUDED.other_phone").
		Set("stage = EXCLUDED.stage").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(context.Background())

	return int64(len(contacts)), err
}

func (t *PgContactStruct) Create(contact *PgContact) error {
	_, err := t.PgDbClient.NewInsert().Model(contact).Exec(context.Background())
	return err
}

func (t *PgContactStruct) GetByUUID(uuid string) (*PgContact, error) {
	contact := new(PgContact)
	err := t.PgDbClient.NewSelect().Model(contact).Where("uuid = ?", uuid).Where("deleted_at IS NULL").Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (t *PgContactStruct) GetByUUIDOrEmail(uuid, email string) (*PgContact, error) {
	contact := new(PgContact)
	query := t.PgDbClient.NewSelect().Model(contact).Where("deleted_at IS NULL")
	
	if uuid != "" {
		query = query.Where("uuid = ?", uuid)
	} else if email != "" {
		query = query.Where("email = ?", email)
	} else {
		return nil, fmt.Errorf("either UUID or email must be provided")
	}
	
	err := query.Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (t *PgContactStruct) Update(contact *PgContact) error {
	_, err := t.PgDbClient.NewUpdate().
		Model(contact).
		Where("uuid = ?", contact.UUID).
		Where("deleted_at IS NULL").
		Column("first_name").
		Column("last_name").
		Column("company_id").
		Column("email").
		Column("title").
		Column("departments").
		Column("mobile_phone").
		Column("email_status").
		Column("seniority").
		Column("city").
		Column("state").
		Column("country").
		Column("linkedin_url").
		Column("facebook_url").
		Column("twitter_url").
		Column("website").
		Column("work_direct_phone").
		Column("home_phone").
		Column("other_phone").
		Column("stage").
		Column("updated_at").
		Exec(context.Background())
	return err
}

func (t *PgContactStruct) Delete(uuid string) error {
	now := time.Now()
	_, err := t.PgDbClient.NewUpdate().
		Model((*PgContact)(nil)).
		Where("uuid = ?", uuid).
		Where("deleted_at IS NULL").
		Set("deleted_at = ?", now).
		Set("updated_at = ?", now).
		Exec(context.Background())
	return err
}

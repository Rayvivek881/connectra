package service

import (
	"errors"
	"sync"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/modules/contacts/helper"
	"vivek-ray/utilities"
)

type ContactService struct {
	contactElasticRepository models.ElasticContactSvcRepo
	contactPgRepository      models.PgContactSvcRepo
	companyPgRepository      models.PgCompanySvcRepo
	filtersDataRepository    models.FiltersDataSvcRepo
	tempFilters              []*models.ModelFilter
}

func NewContactService(tempFilters []*models.ModelFilter) ContactSvcRepo {
	return &ContactService{
		contactElasticRepository: models.ElasticContactRepository(connections.ElasticsearchConnection.Client),
		contactPgRepository:      models.PgContactRepository(connections.PgDBConnection.Client),
		companyPgRepository:      models.PgCompanyRepository(connections.PgDBConnection.Client),
		filtersDataRepository:    models.FiltersDataRepository(connections.PgDBConnection.Client),
		tempFilters:              tempFilters,
	}
}

type ContactSvcRepo interface {
	ListByFilters(query utilities.VQLQuery) ([]helper.ContactResponse, error)
	CountByFilters(query utilities.VQLQuery) (int64, error)
	BulkUpsert(pgContacts []*models.PgContact, esContacts []*models.ElasticContact) ([]*models.PgContact, error)
	BulkUpsertToDb(pgContacts []*models.PgContact, esContacts []*models.ElasticContact, filtersData []*models.ModelFilterData) error
}

func (s *ContactService) ListByFilters(query utilities.VQLQuery) ([]helper.ContactResponse, error) {
	sourceFields := []string{"uuid", "company_id"}
	elasticQuery := query.ToElasticsearchQuery(false, sourceFields)
	esHits, err := s.contactElasticRepository.ListByQueryMap(elasticQuery)
	if err != nil {
		return nil, err
	}
	contactResponses, contactUuids, companyIds := make([]helper.ContactResponse, 0), make([]string, 0), make([]string, 0)
	cursors := make(map[string][]string)
	for _, esHit := range esHits {
		contactUuids = append(contactUuids, esHit.Contact.UUID)
		cursors[esHit.Contact.UUID] = esHit.Cursor
		companyIds = append(companyIds, esHit.Contact.CompanyID)
	}
	if len(query.SelectColumns) != 0 {
		query.SelectColumns = append(query.SelectColumns, "company_id")
	}

	var (
		pgContacts []*models.PgContact
		companies  []*models.PgCompany
		contactErr error
		companyErr error
	)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		pgContacts, contactErr = s.contactPgRepository.ListByFilters(models.PgContactFilters{
			Uuids:         utilities.UniqueStringSlice(contactUuids),
			SelectColumns: query.SelectColumns,
		})
	}()

	shouldPopulateCompanies := query.CompanyConfig != nil && query.CompanyConfig.Populate
	if shouldPopulateCompanies {
		wg.Add(1)
		go func() {
			defer wg.Done()
			companies, companyErr = s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
				Uuids:         utilities.UniqueStringSlice(companyIds),
				SelectColumns: query.CompanyConfig.SelectColumns,
			})
		}()
	}
	wg.Wait()

	if contactErr != nil || companyErr != nil {
		return nil, constants.FailedToFetchDataError
	}

	pgContactsMap := make(map[string]*models.PgContact)
	for _, contact := range pgContacts {
		pgContactsMap[contact.UUID] = contact
	}
	for _, uuid := range contactUuids {
		if contact, ok := pgContactsMap[uuid]; ok {
			contactResponses = append(contactResponses, helper.ContactResponse{
				PgContact: contact,
				Company:   nil,
				Cursor:    cursors[uuid],
			})
		}
	}
	if shouldPopulateCompanies {
		companiesMap := make(map[string]*models.PgCompany)
		for _, company := range companies {
			companiesMap[company.UUID] = company
		}
		for i := range contactResponses {
			contactResponses[i].Company = companiesMap[contactResponses[i].PgContact.CompanyID]
		}
	}

	return contactResponses, nil
}

func (s *ContactService) CountByFilters(query utilities.VQLQuery) (int64, error) {
	elasticQuery := query.ToElasticsearchQuery(true, []string{})
	return s.contactElasticRepository.CountByQueryMap(elasticQuery)
}

func (s *ContactService) BulkUpsertToDb(pgContacts []*models.PgContact,
	esContacts []*models.ElasticContact, filtersData []*models.ModelFilterData) error {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var insertionError error

	wg.Add(3)
	go func() {
		defer wg.Done()
		if _, err := s.contactPgRepository.BulkUpsert(pgContacts); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.contactElasticRepository.BulkUpsert(esContacts); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.filtersDataRepository.BulkUpsert(filtersData); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	wg.Wait()
	return insertionError
}

func (s *ContactService) BulkUpsert(pgContacts []*models.PgContact, esContacts []*models.ElasticContact) ([]*models.PgContact, error) {
	insertedFilters, filtersData := make(map[string]struct{}), make([]*models.ModelFilterData, 0)

	for _, contact := range pgContacts {
		for _, filter := range s.tempFilters {
			if filter.Service != constants.ContactsService {
				continue
			}

			fieldValue := utilities.GetFieldValue(contact, filter.Key)
			values := utilities.ToStringSlice(fieldValue)

			for _, value := range values {
				filterUUID := utilities.GenerateUUID5(filter.Key + filter.Service + value)
				if _, ok := insertedFilters[filterUUID]; ok || value == "" {
					continue
				}
				insertedFilters[filterUUID] = struct{}{}
				filtersData = append(filtersData, &models.ModelFilterData{
					UUID:         filterUUID,
					FilterKey:    filter.Key,
					Service:      filter.Service,
					DisplayValue: value,
					Value:        value,
				})
			}
		}
	}
	return pgContacts, s.BulkUpsertToDb(pgContacts, esContacts, filtersData)
}

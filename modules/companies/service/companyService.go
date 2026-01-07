package service

import (
	"errors"
	"sync"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/modules/companies/helper"
	"vivek-ray/utilities"
)

type CompanyService struct {
	companyElasticRepository models.ElasticCompanySvcRepo
	companyPgRepository      models.PgCompanySvcRepo
	filtersDataRepository    models.FiltersDataSvcRepo
	tempFilters              []*models.ModelFilter
}

func NewCompanyService(tempFilters []*models.ModelFilter) CompanySvcRepo {
	return &CompanyService{
		companyElasticRepository: models.ElasticCompanyRepository(connections.ElasticsearchConnection.Client),
		companyPgRepository:      models.PgCompanyRepository(connections.PgDBConnection.Client),
		filtersDataRepository:    models.FiltersDataRepository(connections.PgDBConnection.Client),
		tempFilters:              tempFilters,
	}
}

type CompanySvcRepo interface {
	ListByFilters(query utilities.VQLQuery) ([]helper.CompanyResponse, error)
	CountByFilters(query utilities.VQLQuery) (int64, error)
	BulkUpsert(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany) error
	BulkUpsertWithDetails(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany) (*helper.BulkOperationResponse, error)
	BulkUpsertToDb(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany, filtersData []*models.ModelFilterData) error
	GetCompanyByUuids(uuids []string, selectColumns []string) ([]*models.PgCompany, error)
	GetCompanyByUUID(uuid string, selectColumns []string) (*models.PgCompany, error)
	Create(company *models.PgCompany) (*models.PgCompany, error)
	Update(company *models.PgCompany) (*models.PgCompany, error)
	Delete(uuid string) error
	Upsert(company *models.PgCompany) (*models.PgCompany, bool, error)
}

func (s *CompanyService) GetCompanyByUuids(uuids []string, selectColumns []string) ([]*models.PgCompany, error) {
	return s.companyPgRepository.ListByFilters(models.PgCompanyFilters{Uuids: uuids, SelectColumns: selectColumns})
}

func (s *CompanyService) GetCompanyByUUID(uuid string, selectColumns []string) (*models.PgCompany, error) {
	companies, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
		Uuids:         []string{uuid},
		SelectColumns: selectColumns,
	})
	if err != nil {
		return nil, err
	}
	if len(companies) == 0 {
		return nil, constants.CompanyNotFoundError
	}
	return companies[0], nil
}

func (s *CompanyService) ListByFilters(query utilities.VQLQuery) ([]helper.CompanyResponse, error) {
	sourceFields := []string{"uuid"}
	elasticQuery := query.ToElasticsearchQuery(false, sourceFields)
	esHits, err := s.companyElasticRepository.ListByQueryMap(elasticQuery)
	if err != nil {
		return nil, err
	}
	companyUuids, cursors := make([]string, 0), make(map[string][]string)
	for _, esHit := range esHits {
		companyUuids = append(companyUuids, esHit.Company.UUID)
		cursors[esHit.Company.UUID] = esHit.Cursor
	}
	companies, err := s.GetCompanyByUuids(companyUuids, query.SelectColumns)
	if err != nil {
		return nil, err
	}
	return helper.ToCompanyResponses(companies, companyUuids, cursors), nil
}

func (s *CompanyService) CountByFilters(query utilities.VQLQuery) (int64, error) {
	elasticQuery := query.ToElasticsearchQuery(true, []string{})
	return s.companyElasticRepository.CountByQueryMap(elasticQuery)
}

func (s *CompanyService) BulkUpsertToDb(pgCompanies []*models.PgCompany,
	esCompanies []*models.ElasticCompany, filtersData []*models.ModelFilterData) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var insertionError error

	wg.Add(3)
	go func() {
		defer wg.Done()
		if _, err := s.companyPgRepository.BulkUpsert(pgCompanies); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.companyElasticRepository.BulkUpsert(esCompanies); err != nil {
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

func (s *CompanyService) BulkUpsert(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany) error {
	insertedFilters, filtersData := make(map[string]struct{}), make([]*models.ModelFilterData, 0)

	for _, company := range pgCompanies {
		for _, filter := range s.tempFilters {
			if filter.Service != constants.CompaniesService {
				continue
			}
			value, _ := utilities.GetFieldValue(company, filter.Key).(string)
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
	return s.BulkUpsertToDb(pgCompanies, esCompanies, filtersData)
}

func (s *CompanyService) BulkUpsertWithDetails(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany) (*helper.BulkOperationResponse, error) {
	totalCount := int64(len(pgCompanies))
	response := &helper.BulkOperationResponse{
		TotalCount:   totalCount,
		SuccessCount: 0,
		ErrorCount:   0,
		Errors:       make([]helper.BulkOperationError, 0),
	}

	// Check which companies already exist to determine created vs updated
	existingUUIDs := make(map[string]bool)
	if totalCount > 0 {
		uuids := make([]string, 0, len(pgCompanies))
		for _, company := range pgCompanies {
			if company.UUID != "" {
				uuids = append(uuids, company.UUID)
			}
		}
		if len(uuids) > 0 {
			existing, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{Uuids: uuids})
			if err == nil {
				for _, company := range existing {
					existingUUIDs[company.UUID] = true
				}
			}
		}
	}

	// Perform bulk upsert
	err := s.BulkUpsert(pgCompanies, esCompanies)
	if err != nil {
		// If there's an error, we can't determine success count accurately
		// Return error details
		response.ErrorCount = totalCount
		response.Errors = append(response.Errors, helper.BulkOperationError{
			Index: -1,
			Error: err.Error(),
		})
		return response, err
	}

	// If successful, all records were processed
	// Note: PostgreSQL ON CONFLICT doesn't tell us created vs updated count
	// We can only estimate based on existing UUIDs check
	response.SuccessCount = totalCount
	response.ErrorCount = 0

	return response, nil
}

func (s *CompanyService) Create(company *models.PgCompany) (*models.PgCompany, error) {
	// Create in PostgreSQL first (source of truth)
	if err := s.companyPgRepository.Create(company); err != nil {
		return nil, err
	}

	// Queue Elasticsearch indexing (async with retry logic)
	esCompany := models.ElasticCompanyFromRawData(company)
	utilities.EnqueueCompanyOperation("create", company.UUID, esCompany)

	// Update filter data if needed
	filtersData := make([]*models.ModelFilterData, 0)
	for _, filter := range s.tempFilters {
		if filter.Service != constants.CompaniesService {
			continue
		}
		value, _ := utilities.GetFieldValue(company, filter.Key).(string)
		if value == "" {
			continue
		}
		filterUUID := utilities.GenerateUUID5(filter.Key + filter.Service + value)
		filtersData = append(filtersData, &models.ModelFilterData{
			UUID:         filterUUID,
			FilterKey:    filter.Key,
			Service:      filter.Service,
			DisplayValue: value,
			Value:        value,
		})
	}
	if len(filtersData) > 0 {
		_ = s.filtersDataRepository.BulkUpsert(filtersData) // Fire-and-forget
	}

	return company, nil
}

func (s *CompanyService) Update(company *models.PgCompany) (*models.PgCompany, error) {
	// Check if company exists
	existingCompany, err := s.companyPgRepository.GetByUUID(company.UUID)
	if err != nil {
		return nil, constants.CompanyNotFoundError
	}

	// Merge updates: only update fields that are provided (non-zero/non-empty)
	// For partial updates, we need to merge with existing data
	if company.Name == "" {
		company.Name = existingCompany.Name
	}
	if company.EmployeesCount == 0 && existingCompany.EmployeesCount > 0 {
		company.EmployeesCount = existingCompany.EmployeesCount
	}
	if len(company.Industries) == 0 && len(existingCompany.Industries) > 0 {
		company.Industries = existingCompany.Industries
	}
	if len(company.Keywords) == 0 && len(existingCompany.Keywords) > 0 {
		company.Keywords = existingCompany.Keywords
	}
	if company.Address == "" {
		company.Address = existingCompany.Address
	}
	if company.AnnualRevenue == 0 && existingCompany.AnnualRevenue > 0 {
		company.AnnualRevenue = existingCompany.AnnualRevenue
	}
	if company.TotalFunding == 0 && existingCompany.TotalFunding > 0 {
		company.TotalFunding = existingCompany.TotalFunding
	}
	if len(company.Technologies) == 0 && len(existingCompany.Technologies) > 0 {
		company.Technologies = existingCompany.Technologies
	}
	if company.City == "" {
		company.City = existingCompany.City
	}
	if company.State == "" {
		company.State = existingCompany.State
	}
	if company.Country == "" {
		company.Country = existingCompany.Country
	}
	if company.LinkedinURL == "" {
		company.LinkedinURL = existingCompany.LinkedinURL
	}
	if company.Website == "" {
		company.Website = existingCompany.Website
	}
	if company.NormalizedDomain == "" {
		company.NormalizedDomain = existingCompany.NormalizedDomain
	}
	if company.FacebookURL == "" {
		company.FacebookURL = existingCompany.FacebookURL
	}
	if company.TwitterURL == "" {
		company.TwitterURL = existingCompany.TwitterURL
	}
	if company.CompanyNameForEmails == "" {
		company.CompanyNameForEmails = existingCompany.CompanyNameForEmails
	}
	if company.PhoneNumber == "" {
		company.PhoneNumber = existingCompany.PhoneNumber
	}
	if company.LatestFunding == "" {
		company.LatestFunding = existingCompany.LatestFunding
	}
	if company.LatestFundingAmount == 0 && existingCompany.LatestFundingAmount > 0 {
		company.LatestFundingAmount = existingCompany.LatestFundingAmount
	}
	if company.LastRaisedAt == "" {
		company.LastRaisedAt = existingCompany.LastRaisedAt
	}
	company.CreatedAt = existingCompany.CreatedAt

	// Update in PostgreSQL
	if err := s.companyPgRepository.Update(company); err != nil {
		return nil, err
	}

	// Queue Elasticsearch update (async with retry logic)
	esCompany := models.ElasticCompanyFromRawData(company)
	utilities.EnqueueCompanyOperation("update", company.UUID, esCompany)

	// Update filter data if needed
	filtersData := make([]*models.ModelFilterData, 0)
	for _, filter := range s.tempFilters {
		if filter.Service != constants.CompaniesService {
			continue
		}
		value, _ := utilities.GetFieldValue(company, filter.Key).(string)
		if value == "" {
			continue
		}
		filterUUID := utilities.GenerateUUID5(filter.Key + filter.Service + value)
		filtersData = append(filtersData, &models.ModelFilterData{
			UUID:         filterUUID,
			FilterKey:    filter.Key,
			Service:      filter.Service,
			DisplayValue: value,
			Value:        value,
		})
	}
	if len(filtersData) > 0 {
		_ = s.filtersDataRepository.BulkUpsert(filtersData) // Fire-and-forget
	}

	return company, nil
}

func (s *CompanyService) Delete(uuid string) error {
	// Check if company exists
	_, err := s.companyPgRepository.GetByUUID(uuid)
	if err != nil {
		return constants.CompanyNotFoundError
	}

	// Soft delete in PostgreSQL
	if err := s.companyPgRepository.Delete(uuid); err != nil {
		return err
	}

	// Queue Elasticsearch delete (async with retry logic)
	utilities.EnqueueCompanyOperation("delete", uuid, nil)

	return nil
}

func (s *CompanyService) Upsert(company *models.PgCompany) (*models.PgCompany, bool, error) {
	// Try to find existing company by UUID or normalized_domain
	existingCompany, err := s.companyPgRepository.GetByUUIDOrDomain(company.UUID, company.NormalizedDomain)
	isNew := false

	if err != nil {
		// Company doesn't exist, create new one
		isNew = true
		createdCompany, createErr := s.Create(company)
		if createErr != nil {
			return nil, false, createErr
		}
		return createdCompany, isNew, nil
	}

	// Company exists, update it
	// Preserve the existing UUID and created_at
	company.UUID = existingCompany.UUID
	company.CreatedAt = existingCompany.CreatedAt

	updatedCompany, updateErr := s.Update(company)
	if updateErr != nil {
		return nil, false, updateErr
	}

	return updatedCompany, isNew, nil
}

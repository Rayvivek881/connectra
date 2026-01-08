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
	BulkUpsert(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany) ([]string, error)
	BulkUpsertToDb(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany, filtersData []*models.ModelFilterData) error
	GetCompanyByUuids(uuids []string, selectColumns []string) ([]*models.PgCompany, error)
}

func (s *CompanyService) GetCompanyByUuids(uuids []string, selectColumns []string) ([]*models.PgCompany, error) {
	return s.companyPgRepository.ListByFilters(models.PgCompanyFilters{Uuids: uuids, SelectColumns: selectColumns})
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

func (s *CompanyService) BulkUpsert(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany) ([]string, error) {
	insertedFilters, filtersData := make(map[string]struct{}), make([]*models.ModelFilterData, 0)
	uuids := make([]string, 0, len(pgCompanies))

	for _, company := range pgCompanies {
		uuids = append(uuids, company.UUID)
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
	return uuids, s.BulkUpsertToDb(pgCompanies, esCompanies, filtersData)
}

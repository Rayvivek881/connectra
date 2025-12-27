package service

import (
	"sync"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
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
	ListByFilters(query utilities.VQLQuery) ([]*models.PgCompany, error)
	CountByFilters(query utilities.VQLQuery) (int64, error)
	BulkUpsert(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany) error
	BulkUpsertToDb(pgCompanies []*models.PgCompany, esCompanies []*models.ElasticCompany, filtersData []*models.ModelFilterData) error
	GetCompanyByUuids(uuids []string, selectColumns []string) ([]*models.PgCompany, error)
}

func (s *CompanyService) GetCompanyByUuids(uuids []string, selectColumns []string) ([]*models.PgCompany, error) {
	return s.companyPgRepository.ListByFilters(models.PgCompanyFilters{Uuids: uuids, SelectColumns: selectColumns})
}

func (s *CompanyService) ListByFilters(query utilities.VQLQuery) ([]*models.PgCompany, error) {
	sourcefields := []string{"id"}
	elasticQuery := query.ToElasticsearchQuery(false, sourcefields)
	elasticCompanies, err := s.companyElasticRepository.ListByQueryMap(elasticQuery)
	if err != nil {
		return nil, err
	}
	companyUuids := make([]string, 0)
	for _, company := range elasticCompanies {
		companyUuids = append(companyUuids, company.Id)
	}
	return s.GetCompanyByUuids(companyUuids, query.SelectColumns)
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
			insertionError = err
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.companyElasticRepository.BulkUpsert(esCompanies); err != nil {
			mu.Lock()
			insertionError = err
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.filtersDataRepository.BulkUpsert(filtersData); err != nil {
			mu.Lock()
			insertionError = err
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

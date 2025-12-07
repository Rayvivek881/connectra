package service

import (
	"vivek-ray/connections"
	"vivek-ray/models"
	"vivek-ray/utilities"
)

type CompanyService struct {
	companyElasticRepository models.ElasticCompanySvcRepo
	companyPgRepository      models.PgCompanySvcRepo
}

func NewCompanyService() CompanySvcRepo {
	return &CompanyService{
		companyElasticRepository: models.ElasticCompanyRepository(connections.ElasticsearchConnection.Client),
		companyPgRepository:      models.PgCompanyRepository(connections.PgDBConnection.Client),
	}
}

type CompanySvcRepo interface {
	ListByFilters(query utilities.NQLQuery) ([]*models.PgCompany, error)
	CountByFilters(query utilities.NQLQuery) (int64, error)
}

func (s *CompanyService) ListByFilters(query utilities.NQLQuery) ([]*models.PgCompany, error) {
	elasticQuery := query.ToElasticsearchQuery(false)
	elasticCompanies, err := s.companyElasticRepository.ListByQueryMap(elasticQuery)
	if err != nil {
		return nil, err
	}
	companyUuids := make([]string, 0)
	for _, company := range elasticCompanies {
		companyUuids = append(companyUuids, company.Id)
	}
	return s.companyPgRepository.ListByFilters(models.PgCompanyFilters{Uuids: companyUuids})
}

func (s *CompanyService) CountByFilters(query utilities.NQLQuery) (int64, error) {
	elasticQuery := query.ToElasticsearchQuery(true)
	return s.companyElasticRepository.CountByQueryMap(elasticQuery)
}

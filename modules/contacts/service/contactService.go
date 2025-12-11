package service

import (
	"vivek-ray/connections"
	"vivek-ray/models"
	"vivek-ray/utilities"
)

type ContactService struct {
	contactElasticRepository models.ElasticContactSvcRepo
	contactPgRepository      models.PgContactSvcRepo
}

func NewContactService() ContactSvcRepo {
	return &ContactService{
		contactElasticRepository: models.ElasticContactRepository(connections.ElasticsearchConnection.Client),
		contactPgRepository:      models.PgContactRepository(connections.PgDBConnection.Client),
	}
}

type ContactSvcRepo interface {
	ListByFilters(query utilities.NQLQuery) ([]*models.PgContact, error)
	CountByFilters(query utilities.NQLQuery) (int64, error)
}

func (s *ContactService) ListByFilters(query utilities.NQLQuery) ([]*models.PgContact, error) {
	elasticQuery := query.ToElasticsearchQuery(false)
	elasticContacts, err := s.contactElasticRepository.ListByQueryMap(elasticQuery)
	if err != nil {
		return nil, err
	}
	contactUuids := make([]string, 0)
	for _, contact := range elasticContacts {
		contactUuids = append(contactUuids, contact.Id)
	}
	return s.contactPgRepository.ListByFilters(models.PgContactFilters{Uuids: contactUuids})
}

func (s *ContactService) CountByFilters(query utilities.NQLQuery) (int64, error) {
	elasticQuery := query.ToElasticsearchQuery(true)
	return s.contactElasticRepository.CountByQueryMap(elasticQuery)
}

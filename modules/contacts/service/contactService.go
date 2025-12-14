package service

import (
	"sync"
	"vivek-ray/connections"
	"vivek-ray/models"
	"vivek-ray/modules/contacts/helper"
	"vivek-ray/utilities"
)

type ContactService struct {
	contactElasticRepository models.ElasticContactSvcRepo
	contactPgRepository      models.PgContactSvcRepo
	companyPgRepository      models.PgCompanySvcRepo
}

func NewContactService() ContactSvcRepo {
	return &ContactService{
		contactElasticRepository: models.ElasticContactRepository(connections.ElasticsearchConnection.Client),
		contactPgRepository:      models.PgContactRepository(connections.PgDBConnection.Client),
		companyPgRepository:      models.PgCompanyRepository(connections.PgDBConnection.Client),
	}
}

type ContactSvcRepo interface {
	ListByFilters(query utilities.VQLQuery) ([]helper.ContactResponse, error)
	CountByFilters(query utilities.VQLQuery) (int64, error)
}

func (s *ContactService) ListByFilters(query utilities.VQLQuery) ([]helper.ContactResponse, error) {
	sourcefields := []string{"id", "company_id"}
	elasticQuery := query.ToElasticsearchQuery(false, sourcefields)
	elasticContacts, err := s.contactElasticRepository.ListByQueryMap(elasticQuery)
	if err != nil {
		return nil, err
	}
	contactResponses, contactUuids, companyIds := make([]helper.ContactResponse, 0), make([]string, 0), make([]string, 0)
	for _, contact := range elasticContacts {
		contactUuids = append(contactUuids, contact.Id)
		companyIds = append(companyIds, contact.CompanyID)
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
	shouldPopulateCompanies := query.CompanyConfig != nil && query.CompanyConfig.Populate
	if shouldPopulateCompanies {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			pgContacts, contactErr = s.contactPgRepository.ListByFilters(models.PgContactFilters{Uuids: contactUuids, SelectColumns: query.SelectColumns})
		}()

		go func() {
			defer wg.Done()
			companies, companyErr = s.companyPgRepository.ListByFilters(models.PgCompanyFilters{Uuids: companyIds, SelectColumns: query.CompanyConfig.SelectColumns})
		}()

		wg.Wait()

		if contactErr != nil {
			return nil, contactErr
		}
		if companyErr != nil {
			return nil, companyErr
		}
	} else {
		pgContacts, err = s.contactPgRepository.ListByFilters(models.PgContactFilters{Uuids: contactUuids, SelectColumns: query.SelectColumns})
		if err != nil {
			return nil, err
		}
	}

	for _, contact := range pgContacts {
		contactResponses = append(contactResponses, helper.ContactResponse{
			PgContact: contact,
			Company:   nil,
		})
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

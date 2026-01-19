package service

import (
	"errors"
	"sync"
	"vivek-ray/connections"
	"vivek-ray/models"
	companyService "vivek-ray/modules/companies/service"
	contactService "vivek-ray/modules/contacts/service"
	"vivek-ray/utilities"

	"github.com/rs/zerolog/log"
)

type BatchUpsertSvc interface {
	ProcessBatchUpsert(batch []map[string]string) error
}

type batchUpsertService struct {
	companyService companyService.CompanySvcRepo
	contactService contactService.ContactSvcRepo
}

func NewBatchUpsertService() BatchUpsertSvc {
	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get temp filters")
		return nil
	}
	return &batchUpsertService{
		companyService: companyService.NewCompanyService(tempFilters),
		contactService: contactService.NewContactService(tempFilters),
	}
}

func (s *batchUpsertService) UpsertBatch(pgCompanies []*models.PgCompany, pgContacts []*models.PgContact,
	esCompanies []*models.ElasticCompany, esContacts []*models.ElasticContact) error {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var insertionError error
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.companyService.BulkUpsert(pgCompanies, esCompanies); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.contactService.BulkUpsert(pgContacts, esContacts); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	wg.Wait()
	return insertionError
}

func (s *batchUpsertService) ProcessBatchUpsert(batch []map[string]string) error {
	batchLen := len(batch)

	pgCompanies := make([]*models.PgCompany, 0, batchLen)
	pgContacts := make([]*models.PgContact, 0, batchLen)
	esCompanies := make([]*models.ElasticCompany, 0, batchLen)
	esContacts := make([]*models.ElasticContact, 0, batchLen)

	insertedCompanies, insertedContacts := make(map[string]struct{}, batchLen), make(map[string]struct{}, batchLen)
	for _, row := range batch {
		cleanedRow := make(map[string]string, len(row))
		for key, value := range row {
			cleanedRow[utilities.GetCleanedString(key)] = utilities.GetCleanedString(value)
		}
		company := models.PgCompanyFromRawData(cleanedRow)
		contact := models.PgContactFromRowData(cleanedRow, company)
		elasticCompany := models.ElasticCompanyFromRawData(company)
		elasticContact := models.ElasticContactFromRawData(contact, company)

		if _, ok := insertedCompanies[company.UUID]; !ok {
			insertedCompanies[company.UUID] = struct{}{}
			pgCompanies = append(pgCompanies, company)
			esCompanies = append(esCompanies, elasticCompany)
		}

		if _, ok := insertedContacts[contact.UUID]; !ok {
			insertedContacts[contact.UUID] = struct{}{}
			pgContacts = append(pgContacts, contact)
			esContacts = append(esContacts, elasticContact)
		}
	}
	return s.UpsertBatch(pgCompanies, pgContacts, esCompanies, esContacts)
}

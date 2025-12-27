package service

import (
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
			insertionError = err
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.contactService.BulkUpsert(pgContacts, esContacts); err != nil {
			mu.Lock()
			insertionError = err
			mu.Unlock()
		}
	}()

	wg.Wait()
	return insertionError
}

func (s *batchUpsertService) ProcessBatchUpsert(batch []map[string]string) error {
	cleanedBatch := make([]map[string]string, 0, len(batch))
	for _, row := range batch {
		cleanedRow := make(map[string]string)
		for key, value := range row {
			cleanedRow[utilities.GetCleanedString(key)] = utilities.GetCleanedString(value)
		}
		cleanedBatch = append(cleanedBatch, cleanedRow)
	}

	pgCompanies := make([]*models.PgCompany, 0)
	pgContacts := make([]*models.PgContact, 0)
	esCompanies := make([]*models.ElasticCompany, 0)
	esContacts := make([]*models.ElasticContact, 0)

	insertedCompanies, insertedContacts := make(map[string]struct{}), make(map[string]struct{})
	for _, row := range cleanedBatch {
		company := models.PgCompanyFromRawData(row)
		contact := models.PgContactFromRowData(row, company)
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

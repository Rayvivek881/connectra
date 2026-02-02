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
	ProcessBatchUpsert(batch []map[string]string) ([]string, []string, error)
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
	osCompanies []*models.OpenSearchCompany, osContacts []*models.OpenSearchContact) error {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var insertionError error
	wg.Add(2)

	go func() {
		defer wg.Done()
		if _, err := s.companyService.BulkUpsert(pgCompanies, osCompanies); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.contactService.BulkUpsert(pgContacts, osContacts); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	wg.Wait()
	return insertionError
}

func (s *batchUpsertService) ProcessBatchUpsert(batch []map[string]string) ([]string, []string, error) {
	batchLen := len(batch)

	pgCompanies := make([]*models.PgCompany, 0, batchLen)
	pgContacts := make([]*models.PgContact, 0, batchLen)
	osCompanies := make([]*models.OpenSearchCompany, 0, batchLen)
	osContacts := make([]*models.OpenSearchContact, 0, batchLen)

	companyUuids, contactUuids := make([]string, 0), make([]string, 0)
	insertedCompanies, insertedContacts := make(map[string]struct{}), make(map[string]struct{})

	for _, row := range batch {
		cleanedRow := make(map[string]string, len(row))
		for key, value := range row {
			cleanedRow[utilities.GetCleanedString(key)] = utilities.GetCleanedString(value)
		}
		company := models.PgCompanyFromRawData(cleanedRow)
		contact := models.PgContactFromRowData(cleanedRow, company)
		osCompany := models.OpenSearchCompanyFromRawData(company)
		osContact := models.OpenSearchContactFromRawData(contact, company)

		if _, ok := insertedCompanies[company.UUID]; !ok {
			insertedCompanies[company.UUID] = struct{}{}
			pgCompanies = append(pgCompanies, company)
			osCompanies = append(osCompanies, osCompany)
			companyUuids = append(companyUuids, company.UUID)
		}

		if _, ok := insertedContacts[contact.UUID]; !ok {
			insertedContacts[contact.UUID] = struct{}{}
			pgContacts = append(pgContacts, contact)
			osContacts = append(osContacts, osContact)
			contactUuids = append(contactUuids, contact.UUID)
		}
	}
	return companyUuids, contactUuids, s.UpsertBatch(pgCompanies, pgContacts, osCompanies, osContacts)
}

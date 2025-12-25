package service

import (
	"sync"
	"vivek-ray/connections"
	"vivek-ray/models"
	"vivek-ray/utilities"
)

type BatchUpsertSvc interface {
	ProcessBatchUpsert(batch []map[string]string) error
}

type batchUpsertService struct {
	companyRepo        models.PgCompanySvcRepo
	contactRepo        models.PgContactSvcRepo
	companyElasticRepo models.ElasticCompanySvcRepo
	contactElasticRepo models.ElasticContactSvcRepo
}

func NewBatchUpsertService() BatchUpsertSvc {
	return &batchUpsertService{
		companyRepo:        models.PgCompanyRepository(connections.PgDBConnection.Client),
		contactRepo:        models.PgContactRepository(connections.PgDBConnection.Client),
		companyElasticRepo: models.ElasticCompanyRepository(connections.ElasticsearchConnection.Client),
		contactElasticRepo: models.ElasticContactRepository(connections.ElasticsearchConnection.Client),
	}
}

func (s *batchUpsertService) UpsertBatch(pgCompanies []*models.PgCompany, pgContacts []*models.PgContact,
	esCompanies []*models.ElasticCompany, esContacts []*models.ElasticContact) error {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var insertionError error
	wg.Add(4)

	go func() {
		defer wg.Done()
		if _, err := s.companyRepo.BulkUpsert(pgCompanies); err != nil {
			mu.Lock()
			insertionError = err
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.contactRepo.BulkUpsert(pgContacts); err != nil {
			mu.Lock()
			insertionError = err
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.companyElasticRepo.BulkUpsert(esCompanies); err != nil {
			mu.Lock()
			insertionError = err
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.contactElasticRepo.BulkUpsert(esContacts); err != nil {
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

	pgCompanies := make([]*models.PgCompany, 0, len(batch))
	pgContacts := make([]*models.PgContact, 0, len(cleanedBatch))
	esCompanies := make([]*models.ElasticCompany, 0, len(cleanedBatch))
	esContacts := make([]*models.ElasticContact, 0, len(cleanedBatch))

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

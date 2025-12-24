package service

import (
	"sync"
	"vivek-ray/connections"
	"vivek-ray/models"
	"vivek-ray/utilities"

	"github.com/rs/zerolog/log"
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

	for _, row := range cleanedBatch {
		company := models.PgCompanyFromRawData(row)
		contact := models.PgContactFromRowData(row, company)
		elasticCompany := models.ElasticCompanyFromRawData(company)
		elasticContact := models.ElasticContactFromRawData(contact, company)

		pgCompanies = append(pgCompanies, company)
		pgContacts = append(pgContacts, contact)
		esCompanies = append(esCompanies, elasticCompany)
		esContacts = append(esContacts, elasticContact)
	}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		if _, err := s.companyRepo.BulkUpsert(pgCompanies); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert companies")
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.contactRepo.BulkUpsert(pgContacts); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert contacts")
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.companyElasticRepo.BulkUpsert(esCompanies); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert companies to elasticsearch")
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.contactElasticRepo.BulkUpsert(esContacts); err != nil {
			log.Error().Err(err).Msg("Failed to bulk upsert contacts to elasticsearch")
		}
	}()

	wg.Wait()
	return nil
}

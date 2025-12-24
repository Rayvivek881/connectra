package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/utilities"
)

type CompanyWriteService struct {
	companyPgRepository models.PgCompanySvcRepo
	companyElasticRepo  models.ElasticCompanySvcRepo
}

func NewCompanyWriteService() CompanyWriteSvcRepo {
	return &CompanyWriteService{
		companyPgRepository: models.PgCompanyRepository(connections.PgDBConnection.Client),
		companyElasticRepo:  models.ElasticCompanyRepository(connections.ElasticsearchConnection.Client),
	}
}

type CompanyWriteSvcRepo interface {
	Create(company *models.PgCompany) (*models.PgCompany, error)
	Update(uuid string, updates map[string]interface{}) (*models.PgCompany, error)
	Delete(uuid string) error
	Upsert(company *models.PgCompany) (*models.PgCompany, bool, error)
	BulkUpsert(companies []models.PgCompany) (*BulkUpsertResult, error)
}

type BulkUpsertResult struct {
	Created              int64         `json:"created"`
	Updated              int64         `json:"updated"`
	Total                int64         `json:"total"`
	BatchesProcessed     int64         `json:"batches_processed"`
	ElasticsearchIndexed int64         `json:"elasticsearch_indexed"`
	ElasticsearchFailed  int64         `json:"elasticsearch_failed"`
	ProcessingTime       time.Duration `json:"processing_time"`
	Errors               []string      `json:"errors,omitempty"`
}

// Create creates a new company in PostgreSQL and indexes it in Elasticsearch
func (s *CompanyWriteService) Create(company *models.PgCompany) (*models.PgCompany, error) {
	ctx := context.Background()
	now := time.Now()
	company.CreatedAt = &now
	company.UpdatedAt = &now

	// Begin transaction
	tx, err := connections.PgDBConnection.Client.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert into PostgreSQL
	_, err = tx.NewInsert().
		Model(company).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to insert company: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Index in Elasticsearch (async - don't block on this)
	go func() {
		if err := s.indexCompanyInElasticsearch(company); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to index company in Elasticsearch: %v\n", err)
		}
	}()

	return company, nil
}

// Update updates an existing company in PostgreSQL and Elasticsearch
func (s *CompanyWriteService) Update(uuid string, updates map[string]interface{}) (*models.PgCompany, error) {
	ctx := context.Background()

	// Set updated_at
	now := time.Now()
	updates["updated_at"] = now

	// Begin transaction
	tx, err := connections.PgDBConnection.Client.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update in PostgreSQL
	_, err = tx.NewUpdate().
		Model((*models.PgCompany)(nil)).
		Where("uuid = ?", uuid).
		Set("updated_at = ?", now).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to set updated_at: %w", err)
	}

	// Apply all updates
	for key, value := range updates {
		if key != "updated_at" && key != "uuid" && key != "id" && key != "created_at" {
			_, err = tx.NewUpdate().
				Model((*models.PgCompany)(nil)).
				Where("uuid = ?", uuid).
				Set(fmt.Sprintf("%s = ?", key), value).
				Exec(ctx)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Fetch updated company
	companies, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
		Uuids: []string{uuid},
	})
	if err != nil || len(companies) == 0 {
		return nil, fmt.Errorf("failed to fetch updated company: %w", err)
	}

	company := companies[0]

	// Update Elasticsearch index (async)
	go func() {
		if err := s.indexCompanyInElasticsearch(company); err != nil {
			fmt.Printf("Warning: Failed to update company in Elasticsearch: %v\n", err)
		}
	}()

	return company, nil
}

// Delete soft deletes a company in PostgreSQL and removes from Elasticsearch
func (s *CompanyWriteService) Delete(uuid string) error {
	ctx := context.Background()
	now := time.Now()

	// Begin transaction
	tx, err := connections.PgDBConnection.Client.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Soft delete in PostgreSQL (set deleted_at)
	_, err = tx.NewUpdate().
		Model((*models.PgCompany)(nil)).
		Where("uuid = ?", uuid).
		Set("deleted_at = ?", now).
		Set("updated_at = ?", now).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Remove from Elasticsearch (async)
	go func() {
		if err := s.deleteCompanyFromElasticsearch(uuid); err != nil {
			fmt.Printf("Warning: Failed to delete company from Elasticsearch: %v\n", err)
		}
	}()

	return nil
}

// Upsert creates or updates a company
func (s *CompanyWriteService) Upsert(company *models.PgCompany) (*models.PgCompany, bool, error) {
	ctx := context.Background()
	now := time.Now()

	// Check if company exists
	existing, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
		Uuids: []string{company.UUID},
	})

	isNew := err != nil || len(existing) == 0

	if isNew {
		company.CreatedAt = &now
		company.UpdatedAt = &now
	} else {
		company.UpdatedAt = &now
	}

	// Begin transaction
	tx, err := connections.PgDBConnection.Client.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Use ON CONFLICT for upsert
	_, err = tx.NewInsert().
		Model(company).
		On("CONFLICT(uuid) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("normalized_domain = EXCLUDED.normalized_domain").
		Set("employees_count = EXCLUDED.employees_count").
		Set("industries = EXCLUDED.industries").
		Set("keywords = EXCLUDED.keywords").
		Set("address = EXCLUDED.address").
		Set("annual_revenue = EXCLUDED.annual_revenue").
		Set("total_funding = EXCLUDED.total_funding").
		Set("technologies = EXCLUDED.technologies").
		Set("city = EXCLUDED.city").
		Set("state = EXCLUDED.state").
		Set("country = EXCLUDED.country").
		Set("linkedin_url = EXCLUDED.linkedin_url").
		Set("website = EXCLUDED.website").
		Set("facebook_url = EXCLUDED.facebook_url").
		Set("twitter_url = EXCLUDED.twitter_url").
		Set("company_name_for_emails = EXCLUDED.company_name_for_emails").
		Set("phone_number = EXCLUDED.phone_number").
		Set("latest_funding = EXCLUDED.latest_funding").
		Set("latest_funding_amount = EXCLUDED.latest_funding_amount").
		Set("last_raised_at = EXCLUDED.last_raised_at").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("failed to upsert company: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Index in Elasticsearch (async)
	go func() {
		if err := s.indexCompanyInElasticsearch(company); err != nil {
			fmt.Printf("Warning: Failed to index company in Elasticsearch: %v\n", err)
		}
	}()

	return company, isNew, nil
}

// BulkUpsert performs bulk upsert of companies with batch processing
func (s *CompanyWriteService) BulkUpsert(companies []models.PgCompany) (*BulkUpsertResult, error) {
	if len(companies) == 0 {
		return &BulkUpsertResult{}, nil
	}

	startTime := time.Now()
	ctx := context.Background()
	now := time.Now()

	// Set timestamps for all companies
	for i := range companies {
		companies[i].UpdatedAt = &now
		if companies[i].CreatedAt == nil {
			companies[i].CreatedAt = &now
		}
	}

	// Get existing company UUIDs to count created vs updated
	uuids := make([]string, len(companies))
	for i, c := range companies {
		uuids[i] = c.UUID
	}

	existing, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
		Uuids: uuids,
	})
	existingMap := make(map[string]bool)
	if err == nil {
		for _, e := range existing {
			existingMap[e.UUID] = true
		}
	}

	// Count new vs updated
	created := int64(0)
	for _, c := range companies {
		if !existingMap[c.UUID] {
			created++
		}
	}
	updated := int64(len(companies)) - created

	// Statistics tracking
	var (
		batchesProcessed     int64
		elasticsearchIndexed int64
		elasticsearchFailed  int64
		errors               []string
		mu                   sync.Mutex
	)

	// Process in batches
	batchSize := constants.DefaultBulkBatchSize
	_, err = utilities.ProcessInBatches(batchSize, companies, func(batch []models.PgCompany) (int, int, error) {
		// Process PostgreSQL batch
		companyPtrs := make([]*models.PgCompany, len(batch))
		for i := range batch {
			companyPtrs[i] = &batch[i]
		}

		// Begin transaction for this batch
		tx, err := connections.PgDBConnection.Client.BeginTx(ctx, nil)
		if err != nil {
			return 0, len(batch), fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback()

		// Bulk upsert in PostgreSQL
		_, err = s.companyPgRepository.BulkUpsert(companyPtrs)
		if err != nil {
			return 0, len(batch), fmt.Errorf("failed to bulk upsert companies: %w", err)
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			return 0, len(batch), fmt.Errorf("failed to commit transaction: %w", err)
		}

		// Convert to ElasticCompany and bulk index in Elasticsearch
		elasticCompanies := make([]*models.ElasticCompany, 0, len(batch))
		for i := range batch {
			elasticCompany := s.convertToElasticCompany(&batch[i])
			elasticCompanies = append(elasticCompanies, elasticCompany)
		}

		// Bulk index in Elasticsearch (async but wait for completion)
		indexed, err := s.companyElasticRepo.BulkUpsert(elasticCompanies)
		if err != nil {
			mu.Lock()
			elasticsearchFailed += int64(len(elasticCompanies))
			errors = append(errors, fmt.Sprintf("batch failed to index in Elasticsearch: %v", err))
			mu.Unlock()
		} else {
			mu.Lock()
			elasticsearchIndexed += indexed
			mu.Unlock()
		}

		mu.Lock()
		batchesProcessed++
		mu.Unlock()

		return len(batch), 0, nil
	})

	if err != nil {
		mu.Lock()
		errors = append(errors, err.Error())
		mu.Unlock()
	}

	processingTime := time.Since(startTime)

	return &BulkUpsertResult{
		Created:              created,
		Updated:              updated,
		Total:                int64(len(companies)),
		BatchesProcessed:     batchesProcessed,
		ElasticsearchIndexed: elasticsearchIndexed,
		ElasticsearchFailed:  elasticsearchFailed,
		ProcessingTime:       processingTime,
		Errors:               errors,
	}, err
}

// convertToElasticCompany converts a PgCompany to ElasticCompany
func (s *CompanyWriteService) convertToElasticCompany(company *models.PgCompany) *models.ElasticCompany {
	var createdAt time.Time
	if company.CreatedAt != nil {
		createdAt = *company.CreatedAt
	} else {
		createdAt = time.Now()
	}

	return &models.ElasticCompany{
		Id:               company.UUID,
		Name:             company.Name,
		NormalizedDomain: company.NormalizedDomain,
		EmployeesCount:   company.EmployeesCount,
		Industries:       company.Industries,
		Keywords:         company.Keywords,
		Address:          company.Address,
		AnnualRevenue:    company.AnnualRevenue,
		TotalFunding:     company.TotalFunding,
		Technologies:     company.Technologies,
		City:             company.City,
		State:            company.State,
		Country:          company.Country,
		LinkedinURL:      company.LinkedinURL,
		Website:          company.Website,
		CreatedAt:        createdAt,
	}
}

// indexCompanyInElasticsearch indexes a company in Elasticsearch (used for single operations)
func (s *CompanyWriteService) indexCompanyInElasticsearch(company *models.PgCompany) error {
	elasticCompany := s.convertToElasticCompany(company)
	return s.companyElasticRepo.IndexCompany(elasticCompany)
}

// deleteCompanyFromElasticsearch removes a company from Elasticsearch
func (s *CompanyWriteService) deleteCompanyFromElasticsearch(uuid string) error {
	return s.companyElasticRepo.DeleteCompany(uuid)
}

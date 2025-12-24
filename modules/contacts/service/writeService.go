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

type ContactWriteService struct {
	contactPgRepository models.PgContactSvcRepo
	contactElasticRepo  models.ElasticContactSvcRepo
	companyPgRepository models.PgCompanySvcRepo
}

func NewContactWriteService() ContactWriteSvcRepo {
	return &ContactWriteService{
		contactPgRepository: models.PgContactRepository(connections.PgDBConnection.Client),
		contactElasticRepo:  models.ElasticContactRepository(connections.ElasticsearchConnection.Client),
		companyPgRepository: models.PgCompanyRepository(connections.PgDBConnection.Client),
	}
}

type ContactWriteSvcRepo interface {
	Create(contact *models.PgContact) (*models.PgContact, error)
	Update(uuid string, updates map[string]interface{}) (*models.PgContact, error)
	Delete(uuid string) error
	Upsert(contact *models.PgContact) (*models.PgContact, bool, error)
	BulkUpsert(contacts []models.PgContact) (*BulkUpsertResult, error)
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

// Create creates a new contact in PostgreSQL and indexes it in Elasticsearch
func (s *ContactWriteService) Create(contact *models.PgContact) (*models.PgContact, error) {
	ctx := context.Background()
	now := time.Now()
	contact.CreatedAt = &now
	contact.UpdatedAt = &now

	// Begin transaction
	tx, err := connections.PgDBConnection.Client.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert into PostgreSQL
	_, err = tx.NewInsert().
		Model(contact).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to insert contact: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Index in Elasticsearch (async - don't block on this)
	go func() {
		if err := s.indexContactInElasticsearch(contact); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to index contact in Elasticsearch: %v\n", err)
		}
	}()

	return contact, nil
}

// Update updates an existing contact in PostgreSQL and Elasticsearch
func (s *ContactWriteService) Update(uuid string, updates map[string]interface{}) (*models.PgContact, error) {
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
		Model((*models.PgContact)(nil)).
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
				Model((*models.PgContact)(nil)).
				Where("uuid = ?", uuid).
				Set(fmt.Sprintf("%s = ?", key), value).
				Exec(ctx)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Fetch updated contact
	contacts, err := s.contactPgRepository.ListByFilters(models.PgContactFilters{
		Uuids: []string{uuid},
	})
	if err != nil || len(contacts) == 0 {
		return nil, fmt.Errorf("failed to fetch updated contact: %w", err)
	}

	contact := contacts[0]

	// Update Elasticsearch index (async)
	go func() {
		if err := s.indexContactInElasticsearch(contact); err != nil {
			fmt.Printf("Warning: Failed to update contact in Elasticsearch: %v\n", err)
		}
	}()

	return contact, nil
}

// Delete soft deletes a contact in PostgreSQL and removes from Elasticsearch
func (s *ContactWriteService) Delete(uuid string) error {
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
		Model((*models.PgContact)(nil)).
		Where("uuid = ?", uuid).
		Set("deleted_at = ?", now).
		Set("updated_at = ?", now).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Remove from Elasticsearch (async)
	go func() {
		if err := s.deleteContactFromElasticsearch(uuid); err != nil {
			fmt.Printf("Warning: Failed to delete contact from Elasticsearch: %v\n", err)
		}
	}()

	return nil
}

// Upsert creates or updates a contact
func (s *ContactWriteService) Upsert(contact *models.PgContact) (*models.PgContact, bool, error) {
	ctx := context.Background()
	now := time.Now()

	// Check if contact exists
	existing, err := s.contactPgRepository.ListByFilters(models.PgContactFilters{
		Uuids: []string{contact.UUID},
	})

	isNew := err != nil || len(existing) == 0

	if isNew {
		contact.CreatedAt = &now
		contact.UpdatedAt = &now
	} else {
		contact.UpdatedAt = &now
	}

	// Begin transaction
	tx, err := connections.PgDBConnection.Client.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Use ON CONFLICT for upsert
	_, err = tx.NewInsert().
		Model(contact).
		On("CONFLICT(uuid) DO UPDATE").
		Set("first_name = EXCLUDED.first_name").
		Set("last_name = EXCLUDED.last_name").
		Set("company_id = EXCLUDED.company_id").
		Set("email = EXCLUDED.email").
		Set("title = EXCLUDED.title").
		Set("departments = EXCLUDED.departments").
		Set("mobile_phone = EXCLUDED.mobile_phone").
		Set("email_status = EXCLUDED.email_status").
		Set("seniority = EXCLUDED.seniority").
		Set("city = EXCLUDED.city").
		Set("state = EXCLUDED.state").
		Set("country = EXCLUDED.country").
		Set("linkedin_url = EXCLUDED.linkedin_url").
		Set("facebook_url = EXCLUDED.facebook_url").
		Set("twitter_url = EXCLUDED.twitter_url").
		Set("website = EXCLUDED.website").
		Set("work_direct_phone = EXCLUDED.work_direct_phone").
		Set("home_phone = EXCLUDED.home_phone").
		Set("other_phone = EXCLUDED.other_phone").
		Set("stage = EXCLUDED.stage").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("failed to upsert contact: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Index in Elasticsearch (async)
	go func() {
		if err := s.indexContactInElasticsearch(contact); err != nil {
			fmt.Printf("Warning: Failed to index contact in Elasticsearch: %v\n", err)
		}
	}()

	return contact, isNew, nil
}

// BulkUpsert performs bulk upsert of contacts with batch processing
func (s *ContactWriteService) BulkUpsert(contacts []models.PgContact) (*BulkUpsertResult, error) {
	if len(contacts) == 0 {
		return &BulkUpsertResult{}, nil
	}

	startTime := time.Now()
	ctx := context.Background()
	now := time.Now()

	// Set timestamps for all contacts
	for i := range contacts {
		contacts[i].UpdatedAt = &now
		if contacts[i].CreatedAt == nil {
			contacts[i].CreatedAt = &now
		}
	}

	// Get existing contact UUIDs to count created vs updated
	uuids := make([]string, len(contacts))
	for i, c := range contacts {
		uuids[i] = c.UUID
	}

	existing, err := s.contactPgRepository.ListByFilters(models.PgContactFilters{
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
	for _, c := range contacts {
		if !existingMap[c.UUID] {
			created++
		}
	}
	updated := int64(len(contacts)) - created

	// Optimize company data fetching: collect all unique company IDs
	companyIDSet := make(map[string]bool)
	for _, c := range contacts {
		if c.CompanyID != "" {
			companyIDSet[c.CompanyID] = true
		}
	}

	// Batch fetch all companies needed for denormalization
	companyIDs := make([]string, 0, len(companyIDSet))
	for id := range companyIDSet {
		companyIDs = append(companyIDs, id)
	}

	companyMap := make(map[string]*models.PgCompany)
	if len(companyIDs) > 0 {
		companies, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
			Uuids: companyIDs,
		})
		if err == nil {
			for i := range companies {
				companyMap[companies[i].UUID] = companies[i]
			}
		}
	}

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
	_, err = utilities.ProcessInBatches(batchSize, contacts, func(batch []models.PgContact) (int, int, error) {
		// Process PostgreSQL batch
		contactPtrs := make([]*models.PgContact, len(batch))
		for i := range batch {
			contactPtrs[i] = &batch[i]
		}

		// Begin transaction for this batch
		tx, err := connections.PgDBConnection.Client.BeginTx(ctx, nil)
		if err != nil {
			return 0, len(batch), fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback()

		// Bulk upsert in PostgreSQL
		_, err = s.contactPgRepository.BulkUpsert(contactPtrs)
		if err != nil {
			return 0, len(batch), fmt.Errorf("failed to bulk upsert contacts: %w", err)
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			return 0, len(batch), fmt.Errorf("failed to commit transaction: %w", err)
		}

		// Convert to ElasticContact with denormalized company data and bulk index
		elasticContacts := make([]*models.ElasticContact, 0, len(batch))
		for i := range batch {
			elasticContact := s.convertToElasticContact(&batch[i], companyMap)
			elasticContacts = append(elasticContacts, elasticContact)
		}

		// Bulk index in Elasticsearch
		indexed, err := s.contactElasticRepo.BulkUpsert(elasticContacts)
		if err != nil {
			mu.Lock()
			elasticsearchFailed += int64(len(elasticContacts))
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
		Total:                int64(len(contacts)),
		BatchesProcessed:     batchesProcessed,
		ElasticsearchIndexed: elasticsearchIndexed,
		ElasticsearchFailed:  elasticsearchFailed,
		ProcessingTime:       processingTime,
		Errors:               errors,
	}, err
}

// convertToElasticContact converts a PgContact to ElasticContact with company data
func (s *ContactWriteService) convertToElasticContact(contact *models.PgContact, companyMap map[string]*models.PgCompany) *models.ElasticContact {
	var createdAt time.Time
	if contact.CreatedAt != nil {
		createdAt = *contact.CreatedAt
	} else {
		createdAt = time.Now()
	}

	elasticContact := &models.ElasticContact{
		Id:          contact.UUID,
		FirstName:   contact.FirstName,
		LastName:    contact.LastName,
		CompanyID:   contact.CompanyID,
		Email:       contact.Email,
		Title:       contact.Title,
		Departments: contact.Departments,
		MobilePhone: contact.MobilePhone,
		EmailStatus: contact.EmailStatus,
		Seniority:   contact.Seniority,
		City:        contact.City,
		State:       contact.State,
		Country:     contact.Country,
		LinkedinURL: contact.LinkedinURL,
		CreatedAt:   createdAt,
	}

	// Add company data if available
	if contact.CompanyID != "" {
		if company, ok := companyMap[contact.CompanyID]; ok {
			elasticContact.CompanyName = company.Name
			elasticContact.CompanyEmployeesCount = company.EmployeesCount
			elasticContact.CompanyIndustries = company.Industries
			elasticContact.CompanyKeywords = company.Keywords
			elasticContact.CompanyAddress = company.Address
			elasticContact.CompanyAnnualRevenue = company.AnnualRevenue
			elasticContact.CompanyTotalFunding = company.TotalFunding
			elasticContact.CompanyTechnologies = company.Technologies
			elasticContact.CompanyCity = company.City
			elasticContact.CompanyState = company.State
			elasticContact.CompanyCountry = company.Country
			elasticContact.CompanyLinkedinURL = company.LinkedinURL
			elasticContact.CompanyWebsite = company.Website
			elasticContact.CompanyNormalizedDomain = company.NormalizedDomain
		}
	}

	return elasticContact
}

// indexContactInElasticsearch indexes a contact in Elasticsearch with company data (used for single operations)
func (s *ContactWriteService) indexContactInElasticsearch(contact *models.PgContact) error {
	// Fetch company data if company_id is set
	var company *models.PgCompany
	if contact.CompanyID != "" {
		companies, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
			Uuids: []string{contact.CompanyID},
		})
		if err == nil && len(companies) > 0 {
			company = companies[0]
		}
	}

	companyMap := make(map[string]*models.PgCompany)
	if company != nil {
		companyMap[company.UUID] = company
	}

	elasticContact := s.convertToElasticContact(contact, companyMap)
	return s.contactElasticRepo.IndexContact(elasticContact)
}

// deleteContactFromElasticsearch removes a contact from Elasticsearch
func (s *ContactWriteService) deleteContactFromElasticsearch(uuid string) error {
	return s.contactElasticRepo.DeleteContact(uuid)
}

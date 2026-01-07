package service

import (
	"errors"
	"fmt"
	"sync"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/modules/contacts/helper"
	"vivek-ray/utilities"
)

type ContactService struct {
	contactElasticRepository models.ElasticContactSvcRepo
	contactPgRepository      models.PgContactSvcRepo
	companyPgRepository      models.PgCompanySvcRepo
	filtersDataRepository    models.FiltersDataSvcRepo
	tempFilters              []*models.ModelFilter
}

func NewContactService(tempFilters []*models.ModelFilter) ContactSvcRepo {
	return &ContactService{
		contactElasticRepository: models.ElasticContactRepository(connections.ElasticsearchConnection.Client),
		contactPgRepository:      models.PgContactRepository(connections.PgDBConnection.Client),
		companyPgRepository:      models.PgCompanyRepository(connections.PgDBConnection.Client),
		filtersDataRepository:    models.FiltersDataRepository(connections.PgDBConnection.Client),
		tempFilters:              tempFilters,
	}
}

type ContactSvcRepo interface {
	ListByFilters(query utilities.VQLQuery) ([]helper.ContactResponse, error)
	CountByFilters(query utilities.VQLQuery) (int64, error)
	BulkUpsert(pgContacts []*models.PgContact, esContacts []*models.ElasticContact) error
	BulkUpsertWithDetails(pgContacts []*models.PgContact, esContacts []*models.ElasticContact) (*helper.BulkOperationResponse, error)
	BulkUpsertToDb(pgContacts []*models.PgContact, esContacts []*models.ElasticContact, filtersData []*models.ModelFilterData) error
	GetContactByUUID(uuid string, selectColumns []string) (*models.PgContact, error)
	Create(contact *models.PgContact) (*models.PgContact, error)
	Update(contact *models.PgContact) (*models.PgContact, error)
	Delete(uuid string) error
	Upsert(contact *models.PgContact) (*models.PgContact, bool, error)
}

func (s *ContactService) GetContactByUUID(uuid string, selectColumns []string) (*models.PgContact, error) {
	contacts, err := s.contactPgRepository.ListByFilters(models.PgContactFilters{
		Uuids:         []string{uuid},
		SelectColumns: selectColumns,
	})
	if err != nil {
		return nil, err
	}
	if len(contacts) == 0 {
		return nil, constants.ContactNotFoundError
	}
	return contacts[0], nil
}

func (s *ContactService) ListByFilters(query utilities.VQLQuery) ([]helper.ContactResponse, error) {
	sourceFields := []string{"uuid", "company_id"}
	elasticQuery := query.ToElasticsearchQuery(false, sourceFields)
	esHits, err := s.contactElasticRepository.ListByQueryMap(elasticQuery)
	if err != nil {
		return nil, err
	}
	contactResponses, contactUuids, companyIds := make([]helper.ContactResponse, 0), make([]string, 0), make([]string, 0)
	cursors := make(map[string][]string)
	for _, esHit := range esHits {
		contactUuids = append(contactUuids, esHit.Contact.UUID)
		cursors[esHit.Contact.UUID] = esHit.Cursor
		companyIds = append(companyIds, esHit.Contact.CompanyID)
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		pgContacts, contactErr = s.contactPgRepository.ListByFilters(models.PgContactFilters{
			Uuids:         utilities.UniqueStringSlice(contactUuids),
			SelectColumns: query.SelectColumns,
		})
	}()

	shouldPopulateCompanies := query.CompanyConfig != nil && query.CompanyConfig.Populate
	if shouldPopulateCompanies {
		wg.Add(1)
		go func() {
			defer wg.Done()
			companies, companyErr = s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
				Uuids:         utilities.UniqueStringSlice(companyIds),
				SelectColumns: query.CompanyConfig.SelectColumns,
			})
		}()
	}
	wg.Wait()

	if contactErr != nil || companyErr != nil {
		return nil, constants.FailedToFetchDataError
	}

	pgContactsMap := make(map[string]*models.PgContact)
	for _, contact := range pgContacts {
		pgContactsMap[contact.UUID] = contact
	}
	for _, uuid := range contactUuids {
		if contact, ok := pgContactsMap[uuid]; ok {
			contactResponses = append(contactResponses, helper.ContactResponse{
				PgContact: contact,
				Company:   nil,
				Cursor:    cursors[uuid],
			})
		}
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

func (s *ContactService) BulkUpsertToDb(pgContacts []*models.PgContact,
	esContacts []*models.ElasticContact, filtersData []*models.ModelFilterData) error {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var insertionError error

	wg.Add(3)
	go func() {
		defer wg.Done()
		if _, err := s.contactPgRepository.BulkUpsert(pgContacts); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := s.contactElasticRepository.BulkUpsert(esContacts); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.filtersDataRepository.BulkUpsert(filtersData); err != nil {
			mu.Lock()
			insertionError = errors.Join(insertionError, err)
			mu.Unlock()
		}
	}()

	wg.Wait()
	return insertionError
}

func (s *ContactService) BulkUpsert(pgContacts []*models.PgContact, esContacts []*models.ElasticContact) error {
	insertedFilters, filtersData := make(map[string]struct{}), make([]*models.ModelFilterData, 0)

	for _, contact := range pgContacts {
		for _, filter := range s.tempFilters {
			if filter.Service != constants.ContactsService {
				continue
			}

			fieldValue := utilities.GetFieldValue(contact, filter.Key)
			values := utilities.ToStringSlice(fieldValue)

			for _, value := range values {
				filterUUID := utilities.GenerateUUID5(filter.Key + filter.Service + value)
				if _, ok := insertedFilters[filterUUID]; ok || value == "" {
					continue
				}
				insertedFilters[filterUUID] = struct{}{}
				filtersData = append(filtersData, &models.ModelFilterData{
					UUID:         filterUUID,
					FilterKey:    filter.Key,
					Service:      filter.Service,
					DisplayValue: value,
					Value:        value,
				})
			}
		}
	}
	return s.BulkUpsertToDb(pgContacts, esContacts, filtersData)
}

func (s *ContactService) BulkUpsertWithDetails(pgContacts []*models.PgContact, esContacts []*models.ElasticContact) (*helper.BulkOperationResponse, error) {
	totalCount := int64(len(pgContacts))
	response := &helper.BulkOperationResponse{
		TotalCount:   totalCount,
		SuccessCount: 0,
		ErrorCount:   0,
		Errors:       make([]helper.BulkOperationError, 0),
	}

	// Check which contacts already exist to determine created vs updated
	existingUUIDs := make(map[string]bool)
	if totalCount > 0 {
		uuids := make([]string, 0, len(pgContacts))
		for _, contact := range pgContacts {
			if contact.UUID != "" {
				uuids = append(uuids, contact.UUID)
			}
		}
		if len(uuids) > 0 {
			existing, err := s.contactPgRepository.ListByFilters(models.PgContactFilters{Uuids: uuids})
			if err == nil {
				for _, contact := range existing {
					existingUUIDs[contact.UUID] = true
				}
			}
		}
	}

	// Perform bulk upsert
	err := s.BulkUpsert(pgContacts, esContacts)
	if err != nil {
		// If there's an error, we can't determine success count accurately
		// Return error details
		response.ErrorCount = totalCount
		response.Errors = append(response.Errors, helper.BulkOperationError{
			Index: -1,
			Error: err.Error(),
		})
		return response, err
	}

	// If successful, all records were processed
	// Note: PostgreSQL ON CONFLICT doesn't tell us created vs updated count
	// We can only estimate based on existing UUIDs check
	response.SuccessCount = totalCount
	response.ErrorCount = 0

	return response, nil
}

func (s *ContactService) Create(contact *models.PgContact) (*models.PgContact, error) {
	// Validate company_id if provided
	var company *models.PgCompany
	if contact.CompanyID != "" {
		companies, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
			Uuids: []string{contact.CompanyID},
		})
		if err != nil || len(companies) == 0 {
			return nil, fmt.Errorf("company with id %s not found", contact.CompanyID)
		}
		company = companies[0]
	} else {
		// Create empty company for Elasticsearch if no company_id
		company = &models.PgCompany{}
	}

	// Create in PostgreSQL first (source of truth)
	if err := s.contactPgRepository.Create(contact); err != nil {
		return nil, err
	}

	// Queue Elasticsearch indexing (async with retry logic)
	// Need company data for denormalized fields
	esContact := models.ElasticContactFromRawData(contact, company)
	utilities.EnqueueContactOperation("create", contact.UUID, esContact)

	// Update filter data if needed
	filtersData := make([]*models.ModelFilterData, 0)
	for _, filter := range s.tempFilters {
		if filter.Service != constants.ContactsService {
			continue
		}

		fieldValue := utilities.GetFieldValue(contact, filter.Key)
		values := utilities.ToStringSlice(fieldValue)

		for _, value := range values {
			if value == "" {
				continue
			}
			filterUUID := utilities.GenerateUUID5(filter.Key + filter.Service + value)
			filtersData = append(filtersData, &models.ModelFilterData{
				UUID:         filterUUID,
				FilterKey:    filter.Key,
				Service:      filter.Service,
				DisplayValue: value,
				Value:        value,
			})
		}
	}
	if len(filtersData) > 0 {
		_ = s.filtersDataRepository.BulkUpsert(filtersData) // Fire-and-forget
	}

	return contact, nil
}

func (s *ContactService) Delete(uuid string) error {
	// Check if contact exists
	_, err := s.contactPgRepository.GetByUUID(uuid)
	if err != nil {
		return constants.ContactNotFoundError
	}

	// Soft delete in PostgreSQL
	if err := s.contactPgRepository.Delete(uuid); err != nil {
		return err
	}

	// Queue Elasticsearch delete (async with retry logic)
	utilities.EnqueueContactOperation("delete", uuid, nil)

	return nil
}

func (s *ContactService) Update(contact *models.PgContact) (*models.PgContact, error) {
	// Check if contact exists
	existingContact, err := s.contactPgRepository.GetByUUID(contact.UUID)
	if err != nil {
		return nil, constants.ContactNotFoundError
	}

	// Merge updates: only update fields that are provided (non-zero/non-empty)
	// For partial updates, we need to merge with existing data
	if contact.FirstName == "" {
		contact.FirstName = existingContact.FirstName
	}
	if contact.LastName == "" {
		contact.LastName = existingContact.LastName
	}
	if contact.Email == "" {
		contact.Email = existingContact.Email
	}
	if contact.CompanyID == "" {
		contact.CompanyID = existingContact.CompanyID
	}
	if contact.Title == "" {
		contact.Title = existingContact.Title
	}
	if len(contact.Departments) == 0 && len(existingContact.Departments) > 0 {
		contact.Departments = existingContact.Departments
	}
	if contact.MobilePhone == "" {
		contact.MobilePhone = existingContact.MobilePhone
	}
	if contact.EmailStatus == "" {
		contact.EmailStatus = existingContact.EmailStatus
	}
	if contact.Seniority == "" {
		contact.Seniority = existingContact.Seniority
	}
	if contact.City == "" {
		contact.City = existingContact.City
	}
	if contact.State == "" {
		contact.State = existingContact.State
	}
	if contact.Country == "" {
		contact.Country = existingContact.Country
	}
	if contact.LinkedinURL == "" {
		contact.LinkedinURL = existingContact.LinkedinURL
	}
	if contact.FacebookURL == "" {
		contact.FacebookURL = existingContact.FacebookURL
	}
	if contact.TwitterURL == "" {
		contact.TwitterURL = existingContact.TwitterURL
	}
	if contact.Website == "" {
		contact.Website = existingContact.Website
	}
	if contact.WorkDirectPhone == "" {
		contact.WorkDirectPhone = existingContact.WorkDirectPhone
	}
	if contact.HomePhone == "" {
		contact.HomePhone = existingContact.HomePhone
	}
	if contact.OtherPhone == "" {
		contact.OtherPhone = existingContact.OtherPhone
	}
	if contact.Stage == "" {
		contact.Stage = existingContact.Stage
	}
	contact.CreatedAt = existingContact.CreatedAt

	// Validate company_id if it changed
	var company *models.PgCompany
	if contact.CompanyID != "" {
		companies, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
			Uuids: []string{contact.CompanyID},
		})
		if err != nil || len(companies) == 0 {
			return nil, fmt.Errorf("company with id %s not found", contact.CompanyID)
		}
		company = companies[0]
	} else {
		// Use existing company or empty company
		if existingContact.CompanyID != "" {
			companies, err := s.companyPgRepository.ListByFilters(models.PgCompanyFilters{
				Uuids: []string{existingContact.CompanyID},
			})
			if err == nil && len(companies) > 0 {
				company = companies[0]
			} else {
				company = &models.PgCompany{}
			}
		} else {
			company = &models.PgCompany{}
		}
	}

	// Update in PostgreSQL
	if err := s.contactPgRepository.Update(contact); err != nil {
		return nil, err
	}

	// Queue Elasticsearch update (async with retry logic)
	// Need company data for denormalized fields
	esContact := models.ElasticContactFromRawData(contact, company)
	utilities.EnqueueContactOperation("update", contact.UUID, esContact)

	// Update filter data if needed
	filtersData := make([]*models.ModelFilterData, 0)
	for _, filter := range s.tempFilters {
		if filter.Service != constants.ContactsService {
			continue
		}

		fieldValue := utilities.GetFieldValue(contact, filter.Key)
		values := utilities.ToStringSlice(fieldValue)

		for _, value := range values {
			if value == "" {
				continue
			}
			filterUUID := utilities.GenerateUUID5(filter.Key + filter.Service + value)
			filtersData = append(filtersData, &models.ModelFilterData{
				UUID:         filterUUID,
				FilterKey:    filter.Key,
				Service:      filter.Service,
				DisplayValue: value,
				Value:        value,
			})
		}
	}
	if len(filtersData) > 0 {
		_ = s.filtersDataRepository.BulkUpsert(filtersData) // Fire-and-forget
	}

	return contact, nil
}

func (s *ContactService) Upsert(contact *models.PgContact) (*models.PgContact, bool, error) {
	// Try to find existing contact by UUID or email
	existingContact, err := s.contactPgRepository.GetByUUIDOrEmail(contact.UUID, contact.Email)
	isNew := false

	if err != nil {
		// Contact doesn't exist, create new one
		isNew = true
		createdContact, createErr := s.Create(contact)
		if createErr != nil {
			return nil, false, createErr
		}
		return createdContact, isNew, nil
	}

	// Contact exists, update it
	// Preserve the existing UUID and created_at
	contact.UUID = existingContact.UUID
	contact.CreatedAt = existingContact.CreatedAt

	updatedContact, updateErr := s.Update(contact)
	if updateErr != nil {
		return nil, false, updateErr
	}

	return updatedContact, isNew, nil
}

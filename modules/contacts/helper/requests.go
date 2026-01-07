package helper

import (
	"fmt"
	"strings"
	"time"
	"vivek-ray/constants"
	"vivek-ray/models"
	companyService "vivek-ray/modules/companies/service"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FilterStatusUpdate struct {
	Active    bool   `json:"active"`
	FilterKey string `json:"filter_key"`
}

func BindAndValidateVQLQuery(c *gin.Context) (utilities.VQLQuery, error) {
	var query utilities.VQLQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		return query, err
	}
	if err := utilities.ValidateElasticPagination(query.Page, query.Limit); err != nil {
		return query, err
	}
	return query, nil
}

func BindAndValidateFiltersDataQuery(c *gin.Context) (models.FiltersDataQuery, error) {
	var query models.FiltersDataQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		return query, err
	}
	if err := utilities.ValidatePageSize(query.Limit); err != nil {
		return query, err
	}
	query.Service = constants.ContactsService
	return query, nil
}

func BindFilterUpdateStatus(c *gin.Context) (FilterStatusUpdate, error) {
	var statusUpdate FilterStatusUpdate
	err := c.ShouldBindJSON(&statusUpdate)

	return statusUpdate, err
}

func BindBatchUpsertRequest(c *gin.Context) ([]*models.PgContact, []*models.ElasticContact, error) {
	pgContacts := make([]*models.PgContact, 0)
	esContacts, companyUuids := make([]*models.ElasticContact, 0), make([]string, 0)
	err := c.ShouldBindJSON(&pgContacts)
	if err != nil {
		return nil, nil, err
	}
	if len(pgContacts) > constants.MaxPageSize {
		return nil, nil, constants.PageSizeExceededError
	}
	for _, contact := range pgContacts {
		if _, err := uuid.Parse(contact.CompanyID); err == nil {
			companyUuids = append(companyUuids, contact.CompanyID)
		}
	}
	companies, err := companyService.NewCompanyService([]*models.ModelFilter{}).GetCompanyByUuids(companyUuids, []string{})
	if err != nil {
		return nil, nil, err
	}
	companyMap := make(map[string]*models.PgCompany)
	for _, company := range companies {
		companyMap[company.UUID] = company
	}
	for _, contact := range pgContacts {
		company := &models.PgCompany{}
		if companyData, ok := companyMap[contact.CompanyID]; ok {
			company = companyData
		}
		esContacts = append(esContacts, models.ElasticContactFromRawData(contact, company))
	}
	return pgContacts, esContacts, nil
}

func BindAndValidateCreateContactRequest(c *gin.Context) (*models.PgContact, error) {
	var contact models.PgContact
	if err := c.ShouldBindJSON(&contact); err != nil {
		return nil, err
	}

	// Validate required fields
	if err := utilities.ValidateRequiredString(contact.FirstName, "first_name"); err != nil {
		return nil, err
	}
	if err := utilities.ValidateRequiredString(contact.LastName, "last_name"); err != nil {
		return nil, err
	}
	if err := utilities.ValidateEmail(contact.Email); err != nil {
		return nil, err
	}

	// Validate company_id if provided
	if contact.CompanyID != "" {
		if err := utilities.ValidateUUID(contact.CompanyID); err != nil {
			return nil, fmt.Errorf("invalid company_id: %v", err)
		}
	}

	// Validate optional URLs
	if contact.LinkedinURL != "" {
		if err := utilities.ValidateLinkedInURL(contact.LinkedinURL); err != nil {
			return nil, err
		}
	}
	if contact.Website != "" {
		if err := utilities.ValidateURL(contact.Website, "website"); err != nil {
			return nil, err
		}
	}
	if contact.FacebookURL != "" {
		if err := utilities.ValidateURL(contact.FacebookURL, "facebook_url"); err != nil {
			return nil, err
		}
	}
	if contact.TwitterURL != "" {
		if err := utilities.ValidateURL(contact.TwitterURL, "twitter_url"); err != nil {
			return nil, err
		}
	}

	// Normalize email and names
	contact.Email = strings.ToLower(strings.TrimSpace(contact.Email))
	contact.FirstName = strings.ToLower(strings.TrimSpace(contact.FirstName))
	contact.LastName = strings.ToLower(strings.TrimSpace(contact.LastName))

	// Normalize LinkedIn URL
	if contact.LinkedinURL != "" {
		contact.LinkedinURL = strings.ToLower(strings.TrimSpace(contact.LinkedinURL))
	}

	// Normalize other string fields
	if contact.City != "" {
		contact.City = strings.ToLower(strings.TrimSpace(contact.City))
	}
	if contact.State != "" {
		contact.State = strings.ToLower(strings.TrimSpace(contact.State))
	}
	if contact.Country != "" {
		contact.Country = strings.ToLower(strings.TrimSpace(contact.Country))
	}
	if contact.EmailStatus != "" {
		contact.EmailStatus = strings.ToLower(strings.TrimSpace(contact.EmailStatus))
	}
	if contact.Seniority != "" {
		contact.Seniority = strings.ToLower(strings.TrimSpace(contact.Seniority))
	}
	if contact.Stage != "" {
		contact.Stage = strings.ToLower(strings.TrimSpace(contact.Stage))
	}

	// Clean phone numbers
	if contact.MobilePhone != "" {
		contact.MobilePhone = utilities.GetCleanedPhoneNumber(contact.MobilePhone)
	}
	if contact.WorkDirectPhone != "" {
		contact.WorkDirectPhone = utilities.GetCleanedPhoneNumber(contact.WorkDirectPhone)
	}
	if contact.HomePhone != "" {
		contact.HomePhone = utilities.GetCleanedPhoneNumber(contact.HomePhone)
	}
	if contact.OtherPhone != "" {
		contact.OtherPhone = utilities.GetCleanedPhoneNumber(contact.OtherPhone)
	}

	// Generate UUID if not provided
	if contact.UUID == "" {
		contact.UUID = utilities.GenerateUUID5(fmt.Sprintf("%s%s%s", contact.FirstName, contact.LastName, contact.LinkedinURL))
	} else {
		// Validate UUID format if provided
		if err := utilities.ValidateUUID(contact.UUID); err != nil {
			return nil, err
		}
	}

	// Set timestamps
	now := time.Now()
	contact.CreatedAt = &now
	contact.UpdatedAt = &now

	return &contact, nil
}

func BindAndValidateUpdateContactRequest(c *gin.Context) (*models.PgContact, error) {
	var contact models.PgContact
	if err := c.ShouldBindJSON(&contact); err != nil {
		return nil, err
	}

	// Get UUID from path parameter
	uuid := c.Param("uuid")
	if uuid == "" {
		return nil, fmt.Errorf("UUID is required in path parameter")
	}
	if err := utilities.ValidateUUID(uuid); err != nil {
		return nil, err
	}
	contact.UUID = uuid

	// Validate email if provided
	if contact.Email != "" {
		if err := utilities.ValidateEmail(contact.Email); err != nil {
			return nil, err
		}
		contact.Email = strings.ToLower(strings.TrimSpace(contact.Email))
	}

	// Validate company_id if provided
	if contact.CompanyID != "" {
		if err := utilities.ValidateUUID(contact.CompanyID); err != nil {
			return nil, fmt.Errorf("invalid company_id: %v", err)
		}
	}

	// Validate optional URLs if provided
	if contact.LinkedinURL != "" {
		if err := utilities.ValidateLinkedInURL(contact.LinkedinURL); err != nil {
			return nil, err
		}
		contact.LinkedinURL = strings.ToLower(strings.TrimSpace(contact.LinkedinURL))
	}
	if contact.Website != "" {
		if err := utilities.ValidateURL(contact.Website, "website"); err != nil {
			return nil, err
		}
	}
	if contact.FacebookURL != "" {
		if err := utilities.ValidateURL(contact.FacebookURL, "facebook_url"); err != nil {
			return nil, err
		}
	}
	if contact.TwitterURL != "" {
		if err := utilities.ValidateURL(contact.TwitterURL, "twitter_url"); err != nil {
			return nil, err
		}
	}

	// Normalize string fields if provided
	if contact.FirstName != "" {
		contact.FirstName = strings.ToLower(strings.TrimSpace(contact.FirstName))
	}
	if contact.LastName != "" {
		contact.LastName = strings.ToLower(strings.TrimSpace(contact.LastName))
	}
	if contact.City != "" {
		contact.City = strings.ToLower(strings.TrimSpace(contact.City))
	}
	if contact.State != "" {
		contact.State = strings.ToLower(strings.TrimSpace(contact.State))
	}
	if contact.Country != "" {
		contact.Country = strings.ToLower(strings.TrimSpace(contact.Country))
	}
	if contact.EmailStatus != "" {
		contact.EmailStatus = strings.ToLower(strings.TrimSpace(contact.EmailStatus))
	}
	if contact.Seniority != "" {
		contact.Seniority = strings.ToLower(strings.TrimSpace(contact.Seniority))
	}
	if contact.Stage != "" {
		contact.Stage = strings.ToLower(strings.TrimSpace(contact.Stage))
	}

	// Clean phone numbers if provided
	if contact.MobilePhone != "" {
		contact.MobilePhone = utilities.GetCleanedPhoneNumber(contact.MobilePhone)
	}
	if contact.WorkDirectPhone != "" {
		contact.WorkDirectPhone = utilities.GetCleanedPhoneNumber(contact.WorkDirectPhone)
	}
	if contact.HomePhone != "" {
		contact.HomePhone = utilities.GetCleanedPhoneNumber(contact.HomePhone)
	}
	if contact.OtherPhone != "" {
		contact.OtherPhone = utilities.GetCleanedPhoneNumber(contact.OtherPhone)
	}

	// Set updated timestamp
	now := time.Now()
	contact.UpdatedAt = &now

	return &contact, nil
}

func BindAndValidateUpsertContactRequest(c *gin.Context) (*models.PgContact, error) {
	var contact models.PgContact
	if err := c.ShouldBindJSON(&contact); err != nil {
		return nil, err
	}

	// Validate email if provided (required for upsert by email)
	if contact.Email == "" && contact.UUID == "" {
		return nil, fmt.Errorf("either email or UUID is required for upsert")
	}

	// Validate email format if provided
	if contact.Email != "" {
		if err := utilities.ValidateEmail(contact.Email); err != nil {
			return nil, err
		}
		contact.Email = strings.ToLower(strings.TrimSpace(contact.Email))
	}

	// Validate UUID format if provided
	if contact.UUID != "" {
		if err := utilities.ValidateUUID(contact.UUID); err != nil {
			return nil, err
		}
	}

	// Validate company_id if provided
	if contact.CompanyID != "" {
		if err := utilities.ValidateUUID(contact.CompanyID); err != nil {
			return nil, fmt.Errorf("invalid company_id: %v", err)
		}
	}

	// Validate optional URLs if provided
	if contact.LinkedinURL != "" {
		if err := utilities.ValidateLinkedInURL(contact.LinkedinURL); err != nil {
			return nil, err
		}
		contact.LinkedinURL = strings.ToLower(strings.TrimSpace(contact.LinkedinURL))
	}
	if contact.Website != "" {
		if err := utilities.ValidateURL(contact.Website, "website"); err != nil {
			return nil, err
		}
	}
	if contact.FacebookURL != "" {
		if err := utilities.ValidateURL(contact.FacebookURL, "facebook_url"); err != nil {
			return nil, err
		}
	}
	if contact.TwitterURL != "" {
		if err := utilities.ValidateURL(contact.TwitterURL, "twitter_url"); err != nil {
			return nil, err
		}
	}

	// Normalize string fields if provided
	if contact.FirstName != "" {
		contact.FirstName = strings.ToLower(strings.TrimSpace(contact.FirstName))
	}
	if contact.LastName != "" {
		contact.LastName = strings.ToLower(strings.TrimSpace(contact.LastName))
	}
	if contact.City != "" {
		contact.City = strings.ToLower(strings.TrimSpace(contact.City))
	}
	if contact.State != "" {
		contact.State = strings.ToLower(strings.TrimSpace(contact.State))
	}
	if contact.Country != "" {
		contact.Country = strings.ToLower(strings.TrimSpace(contact.Country))
	}
	if contact.EmailStatus != "" {
		contact.EmailStatus = strings.ToLower(strings.TrimSpace(contact.EmailStatus))
	}
	if contact.Seniority != "" {
		contact.Seniority = strings.ToLower(strings.TrimSpace(contact.Seniority))
	}
	if contact.Stage != "" {
		contact.Stage = strings.ToLower(strings.TrimSpace(contact.Stage))
	}

	// Clean phone numbers if provided
	if contact.MobilePhone != "" {
		contact.MobilePhone = utilities.GetCleanedPhoneNumber(contact.MobilePhone)
	}
	if contact.WorkDirectPhone != "" {
		contact.WorkDirectPhone = utilities.GetCleanedPhoneNumber(contact.WorkDirectPhone)
	}
	if contact.HomePhone != "" {
		contact.HomePhone = utilities.GetCleanedPhoneNumber(contact.HomePhone)
	}
	if contact.OtherPhone != "" {
		contact.OtherPhone = utilities.GetCleanedPhoneNumber(contact.OtherPhone)
	}

	// Generate UUID if not provided (for new contacts)
	if contact.UUID == "" {
		contact.UUID = utilities.GenerateUUID5(fmt.Sprintf("%s%s%s", contact.FirstName, contact.LastName, contact.LinkedinURL))
	}

	// Set timestamps
	now := time.Now()
	if contact.CreatedAt == nil {
		contact.CreatedAt = &now
	}
	contact.UpdatedAt = &now

	return &contact, nil
}

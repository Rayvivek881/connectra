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

func CleanContactData(c *models.PgContact) *models.PgContact {
	serverTime := time.Now()
	firstName := strings.ToLower(utilities.GetCleanedString(c.FirstName))
	lastName := strings.ToLower(utilities.GetCleanedString(c.LastName))
	linkedinURL := strings.ToLower(utilities.GetCleanedString(c.LinkedinURL))

	if c.UUID == "" {
		c.UUID = utilities.GenerateUUID5(fmt.Sprintf("%s%s%s", firstName, lastName, linkedinURL))
	}

	createdAt := c.CreatedAt
	if createdAt == nil {
		createdAt = &serverTime
	}

	return &models.PgContact{
		UUID: c.UUID,

		FirstName:   firstName,
		LastName:    lastName,
		CompanyID:   utilities.GetCleanedString(c.CompanyID),
		Email:       strings.ToLower(utilities.GetCleanedString(c.Email)),
		Title:       utilities.GetCleanedString(c.Title),
		Departments: cleanStringSlice(c.Departments),

		MobilePhone: utilities.GetCleanedPhoneNumber(c.MobilePhone),
		EmailStatus: strings.ToLower(utilities.GetCleanedString(c.EmailStatus)),
		Seniority:   strings.ToLower(utilities.GetCleanedString(c.Seniority)),
		City:        strings.ToLower(utilities.GetCleanedString(c.City)),
		State:       strings.ToLower(utilities.GetCleanedString(c.State)),
		Country:     strings.ToLower(utilities.GetCleanedString(c.Country)),
		LinkedinURL: linkedinURL,

		FacebookURL:     utilities.GetCleanedString(c.FacebookURL),
		TwitterURL:      utilities.GetCleanedString(c.TwitterURL),
		Website:         utilities.GetCleanedString(c.Website),
		WorkDirectPhone: utilities.GetCleanedPhoneNumber(c.WorkDirectPhone),
		HomePhone:       utilities.GetCleanedPhoneNumber(c.HomePhone),
		OtherPhone:      utilities.GetCleanedPhoneNumber(c.OtherPhone),
		Stage:           strings.ToLower(utilities.GetCleanedString(c.Stage)),

		CreatedAt: createdAt,
		UpdatedAt: &serverTime,
	}
}

func cleanStringSlice(slice []string) []string {
	cleaned := make([]string, 0, len(slice))
	for _, s := range slice {
		if trimmed := utilities.GetCleanedString(s); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func BindBatchUpsertRequest(c *gin.Context) ([]*models.PgContact, []*models.ElasticContact, error) {
	var rawContacts []*models.PgContact
	if err := c.ShouldBindJSON(&rawContacts); err != nil {
		return nil, nil, err
	}
	if len(rawContacts) > constants.MaxPageSize {
		return nil, nil, constants.PageSizeExceededError
	}

	companyUuids := make([]string, 0, len(rawContacts))
	for _, contact := range rawContacts {
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

	pgContacts := make([]*models.PgContact, 0, len(rawContacts))
	esContacts := make([]*models.ElasticContact, 0, len(rawContacts))

	for _, contact := range rawContacts {
		cleanedContact := CleanContactData(contact)
		pgContacts = append(pgContacts, cleanedContact)

		company := &models.PgCompany{}
		if companyData, ok := companyMap[cleanedContact.CompanyID]; ok {
			company = companyData
		}
		esContacts = append(esContacts, models.ElasticContactFromRawData(cleanedContact, company))
	}

	return pgContacts, esContacts, nil
}

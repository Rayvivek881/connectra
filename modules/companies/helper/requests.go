package helper

import (
	"fmt"
	"strings"
	"time"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
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
	if err := utilities.ValidateOpenSearchPagination(query.Page, query.Limit); err != nil {
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
	query.Service = constants.CompaniesService
	return query, nil
}

func BindFilterUpdateStatus(c *gin.Context) (FilterStatusUpdate, error) {
	var statusUpdate FilterStatusUpdate
	err := c.ShouldBindJSON(&statusUpdate)

	return statusUpdate, err
}

func CleanCompanyData(c *models.PgCompany) *models.PgCompany {
	serverTime := time.Now()
	name := utilities.GetCleanedString(c.Name)
	linkedinURL := strings.ToLower(utilities.GetCleanedString(c.LinkedinURL))

	if c.UUID == "" {
		c.UUID = utilities.GenerateUUID5(fmt.Sprintf("%s%s", strings.ToLower(name), linkedinURL))
	}

	createdAt := c.CreatedAt
	if createdAt == nil {
		createdAt = &serverTime
	}

	return &models.PgCompany{
		UUID: c.UUID,

		Name:             name,
		EmployeesCount:   c.EmployeesCount,
		Industries:       cleanStringSlice(c.Industries),
		Keywords:         cleanStringSlice(c.Keywords),
		Address:          utilities.GetCleanedString(c.Address),
		AnnualRevenue:    c.AnnualRevenue,
		TotalFunding:     c.TotalFunding,
		Technologies:     cleanStringSlice(c.Technologies),
		Website:          utilities.GetCleanedString(c.Website),
		LinkedinURL:      linkedinURL,
		City:             strings.ToLower(utilities.GetCleanedString(c.City)),
		State:            strings.ToLower(utilities.GetCleanedString(c.State)),
		Country:          strings.ToLower(utilities.GetCleanedString(c.Country)),
		NormalizedDomain: strings.ToLower(utilities.GetCleanedString(c.NormalizedDomain)),

		FacebookURL:          utilities.GetCleanedString(c.FacebookURL),
		TwitterURL:           utilities.GetCleanedString(c.TwitterURL),
		CompanyNameForEmails: utilities.GetCleanedString(c.CompanyNameForEmails),
		PhoneNumber:          utilities.GetCleanedPhoneNumber(c.PhoneNumber),
		LatestFunding:        utilities.GetCleanedString(c.LatestFunding),
		LatestFundingAmount:  c.LatestFundingAmount,
		LastRaisedAt:         utilities.GetCleanedString(c.LastRaisedAt),

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

func BindBatchUpsertRequest(c *gin.Context) ([]*models.PgCompany, []*models.OpenSearchCompany, error) {
	var rawCompanies []*models.PgCompany
	if err := c.ShouldBindJSON(&rawCompanies); err != nil {
		return nil, nil, err
	}
	if len(rawCompanies) > constants.MaxPageSize {
		return nil, nil, constants.PageSizeExceededError
	}

	pgCompanies := make([]*models.PgCompany, 0, len(rawCompanies))
	osCompanies := make([]*models.OpenSearchCompany, 0, len(rawCompanies))

	for _, company := range rawCompanies {
		cleanedCompany := CleanCompanyData(company)
		pgCompanies = append(pgCompanies, cleanedCompany)
		if !utilities.IsUuidValid(cleanedCompany.UUID) {
			return nil, nil, constants.InvalidUUIDError(cleanedCompany.UUID)
		}
		osCompanies = append(osCompanies, models.OpenSearchCompanyFromRawData(cleanedCompany))
	}

	return pgCompanies, osCompanies, nil
}

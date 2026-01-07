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
	query.Service = constants.CompaniesService
	return query, nil
}

func BindFilterUpdateStatus(c *gin.Context) (FilterStatusUpdate, error) {
	var statusUpdate FilterStatusUpdate
	err := c.ShouldBindJSON(&statusUpdate)

	return statusUpdate, err
}

func BindBatchUpsertRequest(c *gin.Context) ([]*models.PgCompany, []*models.ElasticCompany, error) {
	pgCompanies := make([]*models.PgCompany, 0)
	esCompanies := make([]*models.ElasticCompany, 0)
	err := c.ShouldBindJSON(&pgCompanies)
	if err != nil {
		return nil, nil, err
	}
	if len(pgCompanies) > constants.MaxPageSize {
		return nil, nil, constants.PageSizeExceededError
	}
	for _, pgCompany := range pgCompanies {
		esCompanies = append(esCompanies, models.ElasticCompanyFromRawData(pgCompany))
	}
	return pgCompanies, esCompanies, nil
}

func BindAndValidateCreateCompanyRequest(c *gin.Context) (*models.PgCompany, error) {
	var company models.PgCompany
	if err := c.ShouldBindJSON(&company); err != nil {
		return nil, err
	}

	// Validate required fields
	if err := utilities.ValidateRequiredString(company.Name, "name"); err != nil {
		return nil, err
	}

	// Validate optional URLs
	if err := utilities.ValidateURL(company.Website, "website"); err != nil {
		return nil, err
	}
	if err := utilities.ValidateLinkedInURL(company.LinkedinURL); err != nil {
		return nil, err
	}
	if err := utilities.ValidateURL(company.FacebookURL, "facebook_url"); err != nil {
		return nil, err
	}
	if err := utilities.ValidateURL(company.TwitterURL, "twitter_url"); err != nil {
		return nil, err
	}

	// Validate numeric fields
	if err := utilities.ValidateNonNegativeInt64(company.EmployeesCount, "employees_count"); err != nil {
		return nil, err
	}
	if err := utilities.ValidateNonNegativeInt64(company.AnnualRevenue, "annual_revenue"); err != nil {
		return nil, err
	}
	if err := utilities.ValidateNonNegativeInt64(company.TotalFunding, "total_funding"); err != nil {
		return nil, err
	}
	if err := utilities.ValidateNonNegativeInt64(company.LatestFundingAmount, "latest_funding_amount"); err != nil {
		return nil, err
	}

	// Generate UUID if not provided
	if company.UUID == "" {
		company.UUID = utilities.GenerateUUID5(company.Name + company.LinkedinURL)
	} else {
		// Validate UUID format if provided
		if err := utilities.ValidateUUID(company.UUID); err != nil {
			return nil, err
		}
	}

	// Normalize domain from website if available
	if company.Website != "" && company.NormalizedDomain == "" {
		// Extract domain from website URL
		website := strings.ToLower(company.Website)
		if strings.HasPrefix(website, "http://") || strings.HasPrefix(website, "https://") {
			parts := strings.Split(website, "/")
			if len(parts) > 2 {
				company.NormalizedDomain = strings.TrimPrefix(parts[2], "www.")
			}
		} else {
			company.NormalizedDomain = strings.TrimPrefix(website, "www.")
		}
	}

	// Set timestamps
	now := time.Now()
	company.CreatedAt = &now
	company.UpdatedAt = &now

	return &company, nil
}

func BindAndValidateUpdateCompanyRequest(c *gin.Context) (*models.PgCompany, error) {
	var company models.PgCompany
	if err := c.ShouldBindJSON(&company); err != nil {
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
	company.UUID = uuid

	// Validate optional URLs if provided
	if company.Website != "" {
		if err := utilities.ValidateURL(company.Website, "website"); err != nil {
			return nil, err
		}
	}
	if company.LinkedinURL != "" {
		if err := utilities.ValidateLinkedInURL(company.LinkedinURL); err != nil {
			return nil, err
		}
	}
	if company.FacebookURL != "" {
		if err := utilities.ValidateURL(company.FacebookURL, "facebook_url"); err != nil {
			return nil, err
		}
	}
	if company.TwitterURL != "" {
		if err := utilities.ValidateURL(company.TwitterURL, "twitter_url"); err != nil {
			return nil, err
		}
	}

	// Validate numeric fields if provided
	if company.EmployeesCount < 0 {
		return nil, fmt.Errorf("employees_count must be non-negative")
	}
	if company.AnnualRevenue < 0 {
		return nil, fmt.Errorf("annual_revenue must be non-negative")
	}
	if company.TotalFunding < 0 {
		return nil, fmt.Errorf("total_funding must be non-negative")
	}
	if company.LatestFundingAmount < 0 {
		return nil, fmt.Errorf("latest_funding_amount must be non-negative")
	}

	// Normalize domain from website if provided and normalized_domain is empty
	if company.Website != "" && company.NormalizedDomain == "" {
		website := strings.ToLower(company.Website)
		if strings.HasPrefix(website, "http://") || strings.HasPrefix(website, "https://") {
			parts := strings.Split(website, "/")
			if len(parts) > 2 {
				company.NormalizedDomain = strings.TrimPrefix(parts[2], "www.")
			}
		} else {
			company.NormalizedDomain = strings.TrimPrefix(website, "www.")
		}
	}

	// Set updated timestamp
	now := time.Now()
	company.UpdatedAt = &now

	return &company, nil
}

func BindAndValidateUpsertCompanyRequest(c *gin.Context) (*models.PgCompany, error) {
	var company models.PgCompany
	if err := c.ShouldBindJSON(&company); err != nil {
		return nil, err
	}

	// Validate required fields
	if err := utilities.ValidateRequiredString(company.Name, "name"); err != nil {
		return nil, err
	}

	// Validate optional URLs
	if company.Website != "" {
		if err := utilities.ValidateURL(company.Website, "website"); err != nil {
			return nil, err
		}
	}
	if company.LinkedinURL != "" {
		if err := utilities.ValidateLinkedInURL(company.LinkedinURL); err != nil {
			return nil, err
		}
	}
	if company.FacebookURL != "" {
		if err := utilities.ValidateURL(company.FacebookURL, "facebook_url"); err != nil {
			return nil, err
		}
	}
	if company.TwitterURL != "" {
		if err := utilities.ValidateURL(company.TwitterURL, "twitter_url"); err != nil {
			return nil, err
		}
	}

	// Validate numeric fields
	if company.EmployeesCount < 0 {
		return nil, fmt.Errorf("employees_count must be non-negative")
	}
	if company.AnnualRevenue < 0 {
		return nil, fmt.Errorf("annual_revenue must be non-negative")
	}
	if company.TotalFunding < 0 {
		return nil, fmt.Errorf("total_funding must be non-negative")
	}
	if company.LatestFundingAmount < 0 {
		return nil, fmt.Errorf("latest_funding_amount must be non-negative")
	}

	// Generate UUID if not provided
	if company.UUID == "" {
		company.UUID = utilities.GenerateUUID5(company.Name + company.LinkedinURL)
	} else {
		// Validate UUID format if provided
		if err := utilities.ValidateUUID(company.UUID); err != nil {
			return nil, err
		}
	}

	// Normalize domain from website if available
	if company.Website != "" && company.NormalizedDomain == "" {
		website := strings.ToLower(company.Website)
		if strings.HasPrefix(website, "http://") || strings.HasPrefix(website, "https://") {
			parts := strings.Split(website, "/")
			if len(parts) > 2 {
				company.NormalizedDomain = strings.TrimPrefix(parts[2], "www.")
			}
		} else {
			company.NormalizedDomain = strings.TrimPrefix(website, "www.")
		}
	}

	// Set timestamps (will be updated if existing record found)
	now := time.Now()
	if company.CreatedAt == nil {
		company.CreatedAt = &now
	}
	company.UpdatedAt = &now

	return &company, nil
}

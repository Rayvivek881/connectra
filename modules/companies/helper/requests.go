package helper

import (
	"errors"
	"strings"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
)

type FilterStatusUpdate struct {
	Active    bool   `json:"active"`
	FilterKey string `json:"filter_key"`
}

type BulkUpsertRequest struct {
	Companies []models.PgCompany `json:"companies" binding:"required"`
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

// ValidateCompany validates company data
func ValidateCompany(company *models.PgCompany) error {
	if company.UUID == "" {
		return errors.New("UUID is required")
	}

	// Validate required fields based on business logic
	if company.Name == "" {
		return errors.New("company name is required")
	}

	// Validate LinkedIn URL format if provided
	if company.LinkedinURL != "" {
		if !strings.Contains(company.LinkedinURL, "linkedin.com") {
			return errors.New("invalid LinkedIn URL format")
		}
	}

	// Validate website URL format if provided
	if company.Website != "" {
		if !strings.Contains(company.Website, ".") {
			return errors.New("invalid website URL format")
		}
	}

	// Validate numeric fields are non-negative
	if company.EmployeesCount < 0 {
		return errors.New("employees_count cannot be negative")
	}
	if company.AnnualRevenue < 0 {
		return errors.New("annual_revenue cannot be negative")
	}
	if company.TotalFunding < 0 {
		return errors.New("total_funding cannot be negative")
	}

	return nil
}

// ValidateCompanyUpdates validates company update data
func ValidateCompanyUpdates(uuid string, updates map[string]interface{}) error {
	if uuid == "" {
		return errors.New("UUID is required")
	}

	// Don't allow updating UUID, ID, or created_at
	disallowedFields := []string{"uuid", "id", "created_at"}
	for _, field := range disallowedFields {
		if _, exists := updates[field]; exists {
			return errors.New("cannot update " + field)
		}
	}

	// Validate LinkedIn URL if being updated
	if linkedinURL, exists := updates["linkedin_url"]; exists {
		if urlStr, ok := linkedinURL.(string); ok && urlStr != "" {
			if !strings.Contains(urlStr, "linkedin.com") {
				return errors.New("invalid LinkedIn URL format")
			}
		}
	}

	// Validate website URL if being updated
	if website, exists := updates["website"]; exists {
		if urlStr, ok := website.(string); ok && urlStr != "" {
			if !strings.Contains(urlStr, ".") {
				return errors.New("invalid website URL format")
			}
		}
	}

	// Validate numeric fields
	if employeesCount, exists := updates["employees_count"]; exists {
		if count, ok := employeesCount.(float64); ok && count < 0 {
			return errors.New("employees_count cannot be negative")
		}
	}
	if annualRevenue, exists := updates["annual_revenue"]; exists {
		if revenue, ok := annualRevenue.(float64); ok && revenue < 0 {
			return errors.New("annual_revenue cannot be negative")
		}
	}
	if totalFunding, exists := updates["total_funding"]; exists {
		if funding, ok := totalFunding.(float64); ok && funding < 0 {
			return errors.New("total_funding cannot be negative")
		}
	}

	return nil
}

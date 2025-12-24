package helper

import (
	"errors"
	"net/mail"
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
	Contacts []models.PgContact `json:"contacts" binding:"required"`
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

// ValidateContact validates contact data
func ValidateContact(contact *models.PgContact) error {
	if contact.UUID == "" {
		return errors.New("UUID is required")
	}

	// Validate email format if provided
	if contact.Email != "" {
		if _, err := mail.ParseAddress(contact.Email); err != nil {
			return errors.New("invalid email format")
		}
	}

	// Validate required fields based on business logic
	if contact.FirstName == "" && contact.LastName == "" {
		return errors.New("at least one of first_name or last_name is required")
	}

	// Validate LinkedIn URL format if provided
	if contact.LinkedinURL != "" {
		if !strings.Contains(contact.LinkedinURL, "linkedin.com") {
			return errors.New("invalid LinkedIn URL format")
		}
	}

	return nil
}

// ValidateContactUpdates validates contact update data
func ValidateContactUpdates(uuid string, updates map[string]interface{}) error {
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

	// Validate email format if being updated
	if email, exists := updates["email"]; exists {
		if emailStr, ok := email.(string); ok && emailStr != "" {
			if _, err := mail.ParseAddress(emailStr); err != nil {
				return errors.New("invalid email format")
			}
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

	return nil
}

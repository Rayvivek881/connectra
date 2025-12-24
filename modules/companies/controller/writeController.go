package controller

import (
	"net/http"
	"vivek-ray/models"
	"vivek-ray/modules/companies/helper"
	"vivek-ray/modules/companies/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateCompany handles POST /companies/create
func CreateCompany(c *gin.Context) {
	var company models.PgCompany
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Generate UUID if not provided
	if company.UUID == "" {
		company.UUID = uuid.New().String()
	}

	// Validate company data
	if err := helper.ValidateCompany(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Create company using write service
	result, err := service.NewCompanyWriteService().Create(&company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": result, "success": true})
}

// UpdateCompany handles PUT /companies/:uuid
func UpdateCompany(c *gin.Context) {
	companyUUID := c.Param("uuid")
	if companyUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required", "success": false})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Validate updates
	if err := helper.ValidateCompanyUpdates(companyUUID, updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Update company using write service
	result, err := service.NewCompanyWriteService().Update(companyUUID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

// DeleteCompany handles DELETE /companies/:uuid
func DeleteCompany(c *gin.Context) {
	companyUUID := c.Param("uuid")
	if companyUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required", "success": false})
		return
	}

	// Delete company using write service
	err := service.NewCompanyWriteService().Delete(companyUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully", "success": true})
}

// UpsertCompany handles POST /companies/upsert
func UpsertCompany(c *gin.Context) {
	var company models.PgCompany
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Generate UUID if not provided
	if company.UUID == "" {
		company.UUID = uuid.New().String()
	}

	// Validate company data
	if err := helper.ValidateCompany(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Upsert company using write service
	result, isNew, err := service.NewCompanyWriteService().Upsert(&company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	statusCode := http.StatusOK
	if isNew {
		statusCode = http.StatusCreated
	}

	c.JSON(statusCode, gin.H{"data": result, "is_new": isNew, "success": true})
}

// BulkUpsertCompanies handles POST /companies/bulk
func BulkUpsertCompanies(c *gin.Context) {
	var request helper.BulkUpsertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	if len(request.Companies) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No companies provided", "success": false})
		return
	}

	// Validate all companies
	for i, company := range request.Companies {
		// Generate UUID if not provided
		if company.UUID == "" {
			request.Companies[i].UUID = uuid.New().String()
		}

		if err := helper.ValidateCompany(&request.Companies[i]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         err.Error(),
				"company_index": i,
				"success":       false,
			})
			return
		}
	}

	// Bulk upsert using write service
	result, err := service.NewCompanyWriteService().BulkUpsert(request.Companies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}


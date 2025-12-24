package controller

import (
	"net/http"
	"vivek-ray/models"
	"vivek-ray/modules/contacts/helper"
	"vivek-ray/modules/contacts/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateContact handles POST /contacts/create
func CreateContact(c *gin.Context) {
	var contact models.PgContact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Generate UUID if not provided
	if contact.UUID == "" {
		contact.UUID = uuid.New().String()
	}

	// Validate contact data
	if err := helper.ValidateContact(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Create contact using write service
	result, err := service.NewContactWriteService().Create(&contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": result, "success": true})
}

// UpdateContact handles PUT /contacts/:uuid
func UpdateContact(c *gin.Context) {
	contactUUID := c.Param("uuid")
	if contactUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required", "success": false})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Validate updates
	if err := helper.ValidateContactUpdates(contactUUID, updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Update contact using write service
	result, err := service.NewContactWriteService().Update(contactUUID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

// DeleteContact handles DELETE /contacts/:uuid
func DeleteContact(c *gin.Context) {
	contactUUID := c.Param("uuid")
	if contactUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required", "success": false})
		return
	}

	// Delete contact using write service
	err := service.NewContactWriteService().Delete(contactUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully", "success": true})
}

// UpsertContact handles POST /contacts/upsert
func UpsertContact(c *gin.Context) {
	var contact models.PgContact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Generate UUID if not provided
	if contact.UUID == "" {
		contact.UUID = uuid.New().String()
	}

	// Validate contact data
	if err := helper.ValidateContact(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// Upsert contact using write service
	result, isNew, err := service.NewContactWriteService().Upsert(&contact)
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

// BulkUpsertContacts handles POST /contacts/bulk
func BulkUpsertContacts(c *gin.Context) {
	var request helper.BulkUpsertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	if len(request.Contacts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No contacts provided", "success": false})
		return
	}

	// Validate all contacts
	for i, contact := range request.Contacts {
		// Generate UUID if not provided
		if contact.UUID == "" {
			request.Contacts[i].UUID = uuid.New().String()
		}

		if err := helper.ValidateContact(&request.Contacts[i]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         err.Error(),
				"contact_index": i,
				"success":       false,
			})
			return
		}
	}

	// Bulk upsert using write service
	result, err := service.NewContactWriteService().BulkUpsert(request.Contacts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

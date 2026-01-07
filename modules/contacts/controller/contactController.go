package controller

import (
	"net/http"
	"vivek-ray/connections"
	"vivek-ray/models"
	"vivek-ray/modules/contacts/helper"
	"vivek-ray/modules/contacts/service"

	"github.com/gin-gonic/gin"
)

func GetContactsByFilter(c *gin.Context) {
	query, err := helper.BindAndValidateVQLQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}
	tempFilters := make([]*models.ModelFilter, 0)
	result, err := service.NewContactService(tempFilters).ListByFilters(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

func GetContactsCountByFilter(c *gin.Context) {
	query, err := helper.BindAndValidateVQLQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}
	tempFilters := make([]*models.ModelFilter, 0)
	count, err := service.NewContactService(tempFilters).CountByFilters(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count, "success": true})
}

func BatchUpsert(c *gin.Context) {
	pgContacts, esContacts, err := helper.BindBatchUpsertRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}
	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	pgContacts, _ = service.NewContactService(tempFilters).BulkUpsert(pgContacts, esContacts)
	c.JSON(http.StatusOK, gin.H{"success": true, "contacts": pgContacts})
}

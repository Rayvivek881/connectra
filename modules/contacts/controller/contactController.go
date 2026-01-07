package controller

import (
	"net/http"
	"strings"
	"vivek-ray/connections"
	"vivek-ray/models"
	commonHelper "vivek-ray/modules/common/helper"
	"vivek-ray/modules/contacts/helper"
	"vivek-ray/modules/contacts/service"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
)

func GetContactsByFilter(c *gin.Context) {
	query, err := helper.BindAndValidateVQLQuery(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}
	tempFilters := make([]*models.ModelFilter, 0)
	result, err := service.NewContactService(tempFilters).ListByFilters(query)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

func GetContactsCountByFilter(c *gin.Context) {
	query, err := helper.BindAndValidateVQLQuery(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}
	tempFilters := make([]*models.ModelFilter, 0)
	count, err := service.NewContactService(tempFilters).CountByFilters(query)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count, "success": true})
}

func BatchUpsert(c *gin.Context) {
	pgContacts, esContacts, err := helper.BindBatchUpsertRequest(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}
	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}
	
	result, err := service.NewContactService(tempFilters).BulkUpsertWithDetails(pgContacts, esContacts)
	if err != nil {
		// Return partial success if some records succeeded
		if result != nil {
			c.JSON(http.StatusPartialContent, gin.H{"data": result, "success": false})
			return
		}
		commonHelper.HandleServiceError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

func CreateContact(c *gin.Context) {
	contact, err := helper.BindAndValidateCreateContactRequest(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}

	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	createdContact, err := service.NewContactService(tempFilters).Create(contact)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdContact, "success": true})
}

func UpdateContact(c *gin.Context) {
	contact, err := helper.BindAndValidateUpdateContactRequest(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}

	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	updatedContact, err := service.NewContactService(tempFilters).Update(contact)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedContact, "success": true})
}

func DeleteContact(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		commonHelper.SendBadRequestError(c, "UUID is required in path parameter")
		return
	}

	if err := utilities.ValidateUUID(uuid); err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}

	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	if err := service.NewContactService(tempFilters).Delete(uuid); err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Contact deleted successfully"})
}

func UpsertContact(c *gin.Context) {
	contact, err := helper.BindAndValidateUpsertContactRequest(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}

	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	upsertedContact, isNew, err := service.NewContactService(tempFilters).Upsert(contact)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	statusCode := http.StatusOK
	if isNew {
		statusCode = http.StatusCreated
	}

	c.JSON(statusCode, gin.H{"data": upsertedContact, "is_new": isNew, "success": true})
}

func GetContactByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		commonHelper.SendBadRequestError(c, "UUID is required in path parameter")
		return
	}

	if err := utilities.ValidateUUID(uuid); err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}

	// Parse select_columns from query parameter
	selectColumns := []string{}
	if selectCols := c.Query("select_columns"); selectCols != "" {
		// Parse comma-separated list
		cols := strings.Split(selectCols, ",")
		for _, col := range cols {
			col = strings.TrimSpace(col)
			if col != "" {
				selectColumns = append(selectColumns, col)
			}
		}
	}

	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	contact, err := service.NewContactService(tempFilters).GetContactByUUID(uuid, selectColumns)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contact, "success": true})
}

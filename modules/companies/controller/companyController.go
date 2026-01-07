package controller

import (
	"net/http"
	"strings"
	"vivek-ray/connections"
	"vivek-ray/models"
	commonHelper "vivek-ray/modules/common/helper"
	"vivek-ray/modules/companies/helper"
	"vivek-ray/modules/companies/service"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
)

func GetCompaniesByFilter(c *gin.Context) {
	query, err := helper.BindAndValidateVQLQuery(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}
	tempFilters := make([]*models.ModelFilter, 0)
	result, err := service.NewCompanyService(tempFilters).ListByFilters(query)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

func GetCompaniesCountByFilter(c *gin.Context) {
	query, err := helper.BindAndValidateVQLQuery(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}
	tempFilters := make([]*models.ModelFilter, 0)
	count, err := service.NewCompanyService(tempFilters).CountByFilters(query)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count, "success": true})
}

func BatchUpsert(c *gin.Context) {
	pgCompanies, esCompanies, err := helper.BindBatchUpsertRequest(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}
	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	result, err := service.NewCompanyService(tempFilters).BulkUpsertWithDetails(pgCompanies, esCompanies)
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

func CreateCompany(c *gin.Context) {
	company, err := helper.BindAndValidateCreateCompanyRequest(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}

	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	createdCompany, err := service.NewCompanyService(tempFilters).Create(company)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdCompany, "success": true})
}

func UpdateCompany(c *gin.Context) {
	company, err := helper.BindAndValidateUpdateCompanyRequest(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}

	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	updatedCompany, err := service.NewCompanyService(tempFilters).Update(company)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedCompany, "success": true})
}

func DeleteCompany(c *gin.Context) {
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

	if err := service.NewCompanyService(tempFilters).Delete(uuid); err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Company deleted successfully"})
}

func UpsertCompany(c *gin.Context) {
	company, err := helper.BindAndValidateUpsertCompanyRequest(c)
	if err != nil {
		commonHelper.SendValidationError(c, err.Error(), nil)
		return
	}

	tempFilters, err := models.FiltersRepository(connections.PgDBConnection.Client).GetTempFilters()
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	upsertedCompany, isNew, err := service.NewCompanyService(tempFilters).Upsert(company)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	statusCode := http.StatusOK
	if isNew {
		statusCode = http.StatusCreated
	}

	c.JSON(statusCode, gin.H{"data": upsertedCompany, "is_new": isNew, "success": true})
}

func GetCompanyByUUID(c *gin.Context) {
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

	company, err := service.NewCompanyService(tempFilters).GetCompanyByUUID(uuid, selectColumns)
	if err != nil {
		commonHelper.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": company, "success": true})
}

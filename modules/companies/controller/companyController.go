package controller

import (
	"net/http"
	"vivek-ray/modules/companies/helper"
	"vivek-ray/modules/companies/service"

	"github.com/gin-gonic/gin"
)

func GetCompaniesByFilter(c *gin.Context) {
	query, err := helper.BindAndValidateVQLQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}
	result, err := service.NewCompanyService().ListByFilters(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

func GetCompaniesCountByFilter(c *gin.Context) {
	query, err := helper.BindAndValidateVQLQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}
	count, err := service.NewCompanyService().CountByFilters(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count, "success": true})
}

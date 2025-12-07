package controller

import (
	"net/http"
	"vivek-ray/modules/companies/service"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
)

func GetCompaniesByFilter(c *gin.Context) {
	var query utilities.NQLQuery
	if err := c.ShouldBindJSON(&query); err != nil {
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
	var query utilities.NQLQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}
	count, err := service.NewCompanyService().CountByFilters(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": count, "success": true})
}

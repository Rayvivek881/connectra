package controller

import (
	"net/http"
	"vivek-ray/models"
	"vivek-ray/modules/contacts/service"

	"github.com/gin-gonic/gin"
)

func GetFilters(c *gin.Context) {
	result, err := service.NewFilterService().GetFilters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

func GetFilterData(c *gin.Context) {
	var query models.FiltersDataQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}
	result, err := service.NewFilterService().GetFilterData(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

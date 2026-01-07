package controller

import (
	"net/http"
	"vivek-ray/modules/common/helper"
	"vivek-ray/modules/common/service"

	"github.com/gin-gonic/gin"
)

func GetFilters(c *gin.Context) {
	serviceType := c.Param("service")

	result, err := service.NewFilterService().GetFilters(serviceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

func GetFilterData(c *gin.Context) {
	serviceType := c.Param("service")

	query, err := helper.BindAndValidateFiltersDataQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	result, err := service.NewFilterService().GetFilterData(serviceType, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}


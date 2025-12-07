package controller

import (
	"net/http"
	"vivek-ray/modules/contacts/service"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
)

func GetContactsByFilter(c *gin.Context) {
	var query utilities.NQLQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}
	result, err := service.NewContactService().ListByFilters(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "success": true})
}

func GetContactsCountByFilter(c *gin.Context) {
	var query utilities.NQLQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}
	count, err := service.NewContactService().CountByFilters(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": count, "success": true})
}

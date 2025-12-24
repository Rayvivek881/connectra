package controller

import (
	"net/http"
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
	result, err := service.NewContactService().ListByFilters(query)
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
	count, err := service.NewContactService().CountByFilters(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count, "success": true})
}

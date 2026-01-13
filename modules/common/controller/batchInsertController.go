package controller

import (
	"net/http"
	"vivek-ray/constants"
	"vivek-ray/modules/common/helper"
	"vivek-ray/modules/common/service"

	"github.com/gin-gonic/gin"
)

func BatchUpsert(c *gin.Context) {
	request, err := helper.BindAndValidateBatchInsert(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	batchService := service.NewBatchUpsertService()
	if batchService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.FailedToInitBatchServiceError.Error(), "success": false})
		return
	}
	companyUuids, contactUuids, err := batchService.ProcessBatchUpsert(request.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"company_uuids": companyUuids,
		"contact_uuids": contactUuids,
		"message":       "Batch upsert successful",
		"success":       true,
	})
}

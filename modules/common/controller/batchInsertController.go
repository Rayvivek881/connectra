package controller

import (
	"net/http"
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
	if err := batchService.ProcessBatchUpsert(request.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch upsert successful",
		"count":   len(request.Data),
		"success": true,
	})
}

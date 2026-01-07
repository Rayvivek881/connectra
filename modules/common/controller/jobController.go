package controller

import (
	"net/http"
	"vivek-ray/modules/common/helper"
	"vivek-ray/modules/common/service"

	"github.com/gin-gonic/gin"
)

func CreateJob(c *gin.Context) {
	request, err := helper.BindAndValidateCreateJob(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	err = service.NewJobService().CreateJob(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job created successfully",
		"success": true,
	})
}

func ListJobs(c *gin.Context) {
	request, err := helper.BindAndValidateListJobs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	jobs, err := service.NewJobService().ListJobs(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    jobs,
		"success": true,
	})
}

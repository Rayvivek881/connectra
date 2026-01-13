package controller

import (
	"net/http"
	"vivek-ray/constants"
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

func GetJobByUuid(c *gin.Context) {
	jobUuid := c.Param("job_uuid")
	if jobUuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.JobUuidRequiredError.Error(), "success": false})
		return
	}

	job, err := service.NewJobService().GetJobByUuid(jobUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    job,
		"success": true,
	})
}

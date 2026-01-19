package controller

import (
	"net/http"
	"vivek-ray/connections"
	"vivek-ray/models"
	"vivek-ray/modules/common/helper"
	"vivek-ray/modules/common/service"

	"github.com/gin-gonic/gin"
)

func BulkInsertCompleteGraph(c *gin.Context) {
	requested_nodes, err := helper.GetBulkInsertCompleteGraphRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = service.NewJobService(connections.PgDBConnection.Client).BulkInsertJobs(requested_nodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Jobs inserted successfully"})
}

func UpdateAndRetriggerJob(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	var req helper.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.NewJobService(connections.PgDBConnection.Client).UpdateAndRetriggerJob(uuid, req.Data, req.RetryCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job updated and retriggered successfully"})
}

func GetJobs(c *gin.Context) {
	var filters models.JobFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobs, err := service.NewJobService(connections.PgDBConnection.Client).GetJobs(&filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": jobs})
}

func GetJobByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	job, err := service.NewJobService(connections.PgDBConnection.Client).GetJobByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": job})
}

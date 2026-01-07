package helper

import (
	"encoding/json"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
)

type BatchInsertRequest struct {
	Data []map[string]string `json:"data" binding:"required"`
}

func BindAndValidateBatchInsert(c *gin.Context) (BatchInsertRequest, error) {
	var request BatchInsertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return request, err
	}

	if len(request.Data) == 0 {
		return request, constants.DataArrayEmptyError
	}

	if len(request.Data) > constants.MaxPageSize {
		return request, constants.BatchSizeExceededError
	}

	return request, nil
}

func BindAndValidateFiltersDataQuery(c *gin.Context) (models.FiltersDataQuery, error) {
	var query models.FiltersDataQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		return query, err
	}
	if err := utilities.ValidatePageSize(query.Limit); err != nil {
		return query, err
	}
	return query, nil
}

type CreateJobRequest struct {
	JobType    string          `json:"job_type" binding:"required"`
	JobData    json.RawMessage `json:"job_data" binding:"required"`
	RetryCount int             `json:"retry_count"`
}

func BindAndValidateCreateJob(c *gin.Context) (CreateJobRequest, error) {
	var request CreateJobRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return request, err
	}

	if request.JobType == "" {
		return request, constants.JobTypeRequiredError
	}
	if request.RetryCount < 0 {
		return request, constants.RetryCountNegativeError
	}
	if len(request.JobData) == 0 {
		return request, constants.JobDataRequiredError
	}

	return request, nil
}

type ListJobsRequest struct {
	JobType string   `json:"job_type"`
	Status  []string `json:"status"`
	Limit   int      `json:"limit"`
}

func BindAndValidateListJobs(c *gin.Context) (ListJobsRequest, error) {
	var request ListJobsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return request, err
	}

	if request.Limit < 0 {
		return request, constants.LimitNegativeError
	}
	if request.Limit > constants.MaxPageSize {
		return request, constants.LimitExceededError
	}

	return request, nil
}

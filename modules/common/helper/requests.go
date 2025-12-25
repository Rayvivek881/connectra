package helper

import (
	"encoding/json"
	"errors"
	"strconv"
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
		return request, errors.New("data array is empty")
	}

	if len(request.Data) > constants.MaxPageSize {
		return request, errors.New("batch size cannot exceed " + strconv.Itoa(constants.MaxPageSize))
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
		return request, errors.New("job_type is required")
	}
	if request.RetryCount < 0 {
		return request, errors.New("retry_count cannot be negative")
	}
	if len(request.JobData) == 0 {
		return request, errors.New("job_data is required")
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
		return request, errors.New("limit cannot be negative")
	}
	if request.Limit > 100 {
		return request, errors.New("limit cannot exceed 100")
	}

	return request, nil
}

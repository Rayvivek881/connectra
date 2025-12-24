package helper

import (
	"errors"
	"vivek-ray/models"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
)

const MaxBatchSize = 100

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

	if len(request.Data) > MaxBatchSize {
		return request, errors.New("batch size cannot exceed 100")
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


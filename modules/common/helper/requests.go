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

type UpdateJobRequest struct {
	Data       json.RawMessage `json:"data"`
	RetryCount *int            `json:"retry_count,omitempty"`
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

type BulkInsertGraphRequest struct {
	*models.ModelJobNodes
	Edges []string `json:"edges"`
}

func VerifyDagAndUpdateDegree(dag []*BulkInsertGraphRequest) bool {
	nodeMap := make(map[string]*BulkInsertGraphRequest)
	inDegree := make(map[string]int)

	for _, node := range dag {
		nodeMap[node.UUID] = node
		if _, exists := inDegree[node.UUID]; !exists {
			inDegree[node.UUID] = 0
		}
		for _, target := range node.Edges {
			inDegree[target]++
		}
	}
	for _, node := range dag {
		node.Degree = inDegree[node.UUID]
	}
	queue := make([]string, 0)
	workingDegree := make(map[string]int)
	for uuid, degree := range inDegree {
		workingDegree[uuid] = degree
		if degree == 0 {
			queue = append(queue, uuid)
		}
	}
	processed := 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		processed++

		node, exists := nodeMap[current]
		if !exists {
			continue
		}
		for _, target := range node.Edges {
			workingDegree[target]--
			if workingDegree[target] == 0 {
				queue = append(queue, target)
			}
		}
	}

	return processed == len(inDegree)
}

func GetBulkInsertCompleteGraphRequest(c *gin.Context) ([]*BulkInsertGraphRequest, error) {
	var request_nodes []*BulkInsertGraphRequest
	if err := c.ShouldBindJSON(&request_nodes); err != nil {
		return nil, err
	}
	if len(request_nodes) > constants.MaxNodesPerRequest {
		return nil, constants.ErrInvalidDAG
	}

	if !VerifyDagAndUpdateDegree(request_nodes) {
		return nil, constants.ErrInvalidDAG
	}
	return request_nodes, nil
}

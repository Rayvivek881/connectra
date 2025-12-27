package helper

import (
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
)

type FilterStatusUpdate struct {
	Active    bool   `json:"active"`
	FilterKey string `json:"filter_key"`
}

func BindAndValidateVQLQuery(c *gin.Context) (utilities.VQLQuery, error) {
	var query utilities.VQLQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		return query, err
	}
	if err := utilities.ValidateElasticPagination(query.Page, query.Limit); err != nil {
		return query, err
	}

	return query, nil
}

func BindAndValidateFiltersDataQuery(c *gin.Context) (models.FiltersDataQuery, error) {
	var query models.FiltersDataQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		return query, err
	}
	if err := utilities.ValidatePageSize(query.Limit); err != nil {
		return query, err
	}
	query.Service = constants.CompaniesService
	return query, nil
}

func BindFilterUpdateStatus(c *gin.Context) (FilterStatusUpdate, error) {
	var statusUpdate FilterStatusUpdate
	err := c.ShouldBindJSON(&statusUpdate)

	return statusUpdate, err
}

func BindBatchUpsertRequest(c *gin.Context) ([]*models.PgCompany, []*models.ElasticCompany, error) {
	pgCompanies := make([]*models.PgCompany, 0)
	esCompanies := make([]*models.ElasticCompany, 0)
	err := c.ShouldBindJSON(&pgCompanies)
	if err != nil {
		return nil, nil, err
	}
	if len(pgCompanies) > constants.MaxPageSize {
		return nil, nil, constants.PageSizeExceededError
	}
	for _, pgCompany := range pgCompanies {
		esCompanies = append(esCompanies, models.ElasticCompanyFromRawData(pgCompany))
	}
	return pgCompanies, esCompanies, nil
}

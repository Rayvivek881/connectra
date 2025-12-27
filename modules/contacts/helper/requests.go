package helper

import (
	"vivek-ray/constants"
	"vivek-ray/models"
	companyService "vivek-ray/modules/companies/service"
	"vivek-ray/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	query.Service = constants.ContactsService
	return query, nil
}

func BindFilterUpdateStatus(c *gin.Context) (FilterStatusUpdate, error) {
	var statusUpdate FilterStatusUpdate
	err := c.ShouldBindJSON(&statusUpdate)

	return statusUpdate, err
}

func BindBatchUpsertRequest(c *gin.Context) ([]*models.PgContact, []*models.ElasticContact, error) {
	pgContacts := make([]*models.PgContact, 0)
	esContacts, companyUuids := make([]*models.ElasticContact, 0), make([]string, 0)
	err := c.ShouldBindJSON(&pgContacts)
	if err != nil {
		return nil, nil, err
	}
	if len(pgContacts) > constants.MaxPageSize {
		return nil, nil, constants.PageSizeExceededError
	}
	for _, contact := range pgContacts {
		if _, err := uuid.Parse(contact.CompanyID); err == nil {
			companyUuids = append(companyUuids, contact.CompanyID)
		}
	}
	companies, err := companyService.NewCompanyService([]*models.ModelFilter{}).GetCompanyByUuids(companyUuids, []string{})
	if err != nil {
		return nil, nil, err
	}
	companyMap := make(map[string]*models.PgCompany)
	for _, company := range companies {
		companyMap[company.UUID] = company
	}
	for _, contact := range pgContacts {
		company := &models.PgCompany{}
		if companyData, ok := companyMap[contact.CompanyID]; ok {
			company = companyData
		}
		esContacts = append(esContacts, models.ElasticContactFromRawData(contact, company))
	}
	return pgContacts, esContacts, nil
}

package helper

import "vivek-ray/models"

type FilterDataResponse struct {
	Value        string `json:"value"`
	DisplayValue string `json:"display_value"`
}

func ToFilterDataResponses(data []*models.ModelFilterData) []FilterDataResponse {
	responses := make([]FilterDataResponse, 0)
	for _, item := range data {
		responses = append(responses, FilterDataResponse{
			Value:        item.Value,
			DisplayValue: item.DisplayValue,
		})
	}
	return responses
}

type CompanyResponse struct {
	*models.PgCompany
	Cursor []string `json:"cursor,omitempty"`
}

func ToCompanyResponses(companies []*models.PgCompany, orderedUuids []string, cursors map[string][]string) []CompanyResponse {
	responses := make([]CompanyResponse, 0)
	companiesMap := make(map[string]*models.PgCompany)
	for _, company := range companies {
		companiesMap[company.UUID] = company
	}
	for _, uuid := range orderedUuids {
		if company, ok := companiesMap[uuid]; ok {
			responses = append(responses, CompanyResponse{
				PgCompany: company,
				Cursor:    cursors[uuid],
			})
		}
	}
	return responses
}

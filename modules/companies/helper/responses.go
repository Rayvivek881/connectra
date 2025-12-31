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

func ToCompanyResponses(companies []*models.PgCompany, cursors map[string][]string) []CompanyResponse {
	responses := make([]CompanyResponse, 0)
	for _, company := range companies {
		responses = append(responses, CompanyResponse{
			PgCompany: company,
			Cursor:    cursors[company.UUID],
		})
	}
	return responses
}


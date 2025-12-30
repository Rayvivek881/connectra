package helper

import "vivek-ray/models"

type FilterDataResponse struct {
	Value        string `json:"value"`
	DisplayValue string `json:"display_value"`
}

func ToFilterDataResponses(data []*models.ModelFilterData) []FilterDataResponse {
	responses := make([]FilterDataResponse, 0)
	for _, curr_data := range data {
		responses = append(responses, FilterDataResponse{
			Value:        curr_data.Value,
			DisplayValue: curr_data.DisplayValue,
		})
	}
	return responses
}

type CompanyResponse struct {
	*models.PgCompany
	SearchAfter []string `json:"search_after,omitempty"`
}

func ToCompanyResponses(companies []*models.PgCompany, searchAfter []string) []CompanyResponse {
	responses := make([]CompanyResponse, 0)
	for _, company := range companies {
		responses = append(responses, CompanyResponse{
			PgCompany:   company,
			SearchAfter: searchAfter,
		})
	}
	return responses
}

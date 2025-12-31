package helper

import (
	"vivek-ray/models"
)

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

type ContactResponse struct {
	*models.PgContact
	Company *models.PgCompany `json:"company,omitempty"`
	Cursor  []string          `json:"cursor,omitempty"`
}

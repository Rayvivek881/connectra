package helper

import "vivek-ray/models"

type FilterDataResponse struct {
	Value        string `json:"value"`
	DisplayValue string `json:"display_value"`
}

func ToFilterDataResponses(data []*models.ModelFilterData) []FilterDataResponse {
	responses := make([]FilterDataResponse, 0, len(data))
	for _, d := range data {
		responses = append(responses, FilterDataResponse{
			Value:        d.Value,
			DisplayValue: d.DisplayValue,
		})
	}
	return responses
}


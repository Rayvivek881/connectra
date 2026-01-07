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

type BulkOperationError struct {
	Index int    `json:"index"`
	UUID  string `json:"uuid,omitempty"`
	Error string `json:"error"`
}

type BulkOperationResponse struct {
	TotalCount   int64                `json:"total_count"`
	SuccessCount int64                `json:"success_count"`
	ErrorCount   int64                `json:"error_count"`
	Errors       []BulkOperationError `json:"errors,omitempty"`
}

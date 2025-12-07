package service

import (
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/utilities"
)

type FilterService struct {
	filtersRepository     models.FiltersSvcRepo
	filtersDataRepository models.FiltersDataSvcRepo
	pgCompanyRepository   models.PgCompanySvcRepo
}

func NewFilterService() FilterSvcRepo {
	filtersRepository := models.FiltersRepository(connections.PgDBConnection.Client)
	filtersDataRepository := models.FiltersDataRepository(connections.PgDBConnection.Client)
	pgCompanyRepository := models.PgCompanyRepository(connections.PgDBConnection.Client)
	return &FilterService{
		filtersRepository:     filtersRepository,
		filtersDataRepository: filtersDataRepository,
		pgCompanyRepository:   pgCompanyRepository,
	}
}

type FilterSvcRepo interface {
	GetFilters() ([]*models.ModelFilter, error)
	GetFilterData(query models.FiltersDataQuery) ([]string, error)
}

func (s *FilterService) GetFilters() ([]*models.ModelFilter, error) {
	return s.filtersRepository.GetFiltersByService(constants.CompaniesService)
}

func (s *FilterService) GetFilterData(query models.FiltersDataQuery) ([]string, error) {
	filterData, err := s.filtersRepository.GetFilterByKeyAndService(query.Service, query.FilterKey)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	if !filterData.DirectDrived {
		curr_lst, err := s.filtersDataRepository.GetFiltersByQuery(query)
		if err == nil {
			for _, curr_data := range curr_lst {
				result = append(result, curr_data.Value)
			}
		}
		return result, err
	}
	curr_lst, err := s.pgCompanyRepository.GetFiltersByQuery(query)
	if err == nil {
		for _, curr_data := range curr_lst {
			fieldValue := utilities.GetFieldValue(curr_data, query.FilterKey)
			if fieldValue != "" {
				result = append(result, fieldValue)
			}
		}
	}
	return result, err
}

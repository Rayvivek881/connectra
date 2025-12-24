package service

import (
	"errors"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
	"vivek-ray/modules/common/helper"
	"vivek-ray/utilities"
)

type FilterSvc interface {
	GetFilters(serviceType string) ([]*models.ModelFilter, error)
	GetFilterData(serviceType string, query models.FiltersDataQuery) ([]helper.FilterDataResponse, error)
}

type filterService struct {
	filtersRepository     models.FiltersSvcRepo
	filtersDataRepository models.FiltersDataSvcRepo
	pgCompanyRepository   models.PgCompanySvcRepo
	pgContactRepository   models.PgContactSvcRepo
}

func NewFilterService() FilterSvc {
	return &filterService{
		filtersRepository:     models.FiltersRepository(connections.PgDBConnection.Client),
		filtersDataRepository: models.FiltersDataRepository(connections.PgDBConnection.Client),
		pgCompanyRepository:   models.PgCompanyRepository(connections.PgDBConnection.Client),
		pgContactRepository:   models.PgContactRepository(connections.PgDBConnection.Client),
	}
}

func (s *filterService) GetFilters(serviceType string) ([]*models.ModelFilter, error) {
	if !isValidService(serviceType) {
		return nil, errors.New("invalid service type")
	}
	return s.filtersRepository.GetFiltersByService(serviceType)
}

func (s *filterService) GetFilterData(serviceType string, query models.FiltersDataQuery) ([]helper.FilterDataResponse, error) {
	if !isValidService(serviceType) {
		return nil, errors.New("invalid service type")
	}

	query.Service = serviceType

	filterData, err := s.filtersRepository.GetFilterByKeyAndService(serviceType, query.FilterKey)
	if err != nil {
		return nil, err
	}

	if !filterData.DirectDrived {
		data, err := s.filtersDataRepository.GetFiltersByQuery(query)
		if err != nil {
			return nil, err
		}
		return helper.ToFilterDataResponses(data), nil
	}

	return s.getDirectDerivedFilterData(serviceType, query)
}

func (s *filterService) getDirectDerivedFilterData(serviceType string, query models.FiltersDataQuery) ([]helper.FilterDataResponse, error) {
	result := make([]helper.FilterDataResponse, 0)

	switch serviceType {
	case constants.CompaniesService:
		data, err := s.pgCompanyRepository.GetFiltersByQuery(query)
		if err != nil {
			return nil, err
		}
		for _, item := range data {
			if fieldValue := utilities.GetFieldValue(item, query.FilterKey); fieldValue != "" {
				result = append(result, helper.FilterDataResponse{
					Value:        fieldValue,
					DisplayValue: fieldValue,
				})
			}
		}

	case constants.ContactsService:
		data, err := s.pgContactRepository.GetFiltersByQuery(query)
		if err != nil {
			return nil, err
		}
		for _, item := range data {
			if fieldValue := utilities.GetFieldValue(item, query.FilterKey); fieldValue != "" {
				result = append(result, helper.FilterDataResponse{
					Value:        fieldValue,
					DisplayValue: fieldValue,
				})
			}
		}
	}

	return result, nil
}

func isValidService(serviceType string) bool {
	return serviceType == constants.CompaniesService || serviceType == constants.ContactsService
}


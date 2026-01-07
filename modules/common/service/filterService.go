package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
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
		return nil, constants.InvalidServiceTypeError
	}
	
	// Try to get from cache
	cacheKey := fmt.Sprintf("filters:%s", serviceType)
	if cached, found := utilities.FilterMetadataCache.Get(cacheKey); found {
		if filters, ok := cached.([]*models.ModelFilter); ok {
			return filters, nil
		}
	}
	
	// Fetch from database
	filters, err := s.filtersRepository.GetFiltersByService(serviceType)
	if err != nil {
		return nil, err
	}
	
	// Store in cache
	utilities.FilterMetadataCache.Set(cacheKey, filters)
	
	return filters, nil
}

func (s *filterService) GetFilterData(serviceType string, query models.FiltersDataQuery) ([]helper.FilterDataResponse, error) {
	if !isValidService(serviceType) {
		return nil, constants.InvalidServiceTypeError
	}

	query.Service = serviceType

	// Generate cache key from query parameters
	cacheKey := generateFilterDataCacheKey(serviceType, query)
	
	// Try to get from cache
	if cached, found := utilities.FilterDataCache.Get(cacheKey); found {
		if data, ok := cached.([]helper.FilterDataResponse); ok {
			return data, nil
		}
	}

	filterData, err := s.filtersRepository.GetFilterByKeyAndService(serviceType, query.FilterKey)
	if err != nil {
		return nil, err
	}

	var result []helper.FilterDataResponse
	
	if !filterData.DirectDerived {
		data, err := s.filtersDataRepository.GetFiltersByQuery(query)
		if err != nil {
			return nil, err
		}
		result = helper.ToFilterDataResponses(data)
	} else {
		result, err = s.getDirectDerivedFilterData(serviceType, query)
		if err != nil {
			return nil, err
		}
	}
	
	// Store in cache (shorter TTL for direct-derived as they change more frequently)
	ttl := 5 * time.Minute
	if filterData.DirectDerived {
		ttl = 2 * time.Minute // Shorter TTL for direct-derived filters
	}
	utilities.FilterDataCache.SetWithTTL(cacheKey, result, ttl)
	
	return result, nil
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
			if fieldValue, ok := utilities.GetFieldValue(item, query.FilterKey).(string); ok && fieldValue != "" {
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
			if fieldValue, ok := utilities.GetFieldValue(item, query.FilterKey).(string); ok && fieldValue != "" {
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

// generateFilterDataCacheKey creates a cache key from filter data query parameters
func generateFilterDataCacheKey(serviceType string, query models.FiltersDataQuery) string {
	// Create a stable key from query parameters
	keyData := map[string]interface{}{
		"service":    serviceType,
		"filter_key": query.FilterKey,
		"search_text": query.SearchText,
		"page":       query.Page,
		"limit":      query.Limit,
	}
	
	// Serialize to JSON for consistent key generation
	jsonData, err := json.Marshal(keyData)
	if err != nil {
		// Fallback to simple string concatenation
		return fmt.Sprintf("filter_data:%s:%s:%s:%d:%d", serviceType, query.FilterKey, query.SearchText, query.Page, query.Limit)
	}
	
	// Hash the JSON for shorter keys
	hash := md5.Sum(jsonData)
	return fmt.Sprintf("filter_data:%s:%x", serviceType, hash)
}

// InvalidateFilterCache invalidates cache entries for a specific service
func InvalidateFilterCache(serviceType string) {
	utilities.FilterMetadataCache.Delete(fmt.Sprintf("filters:%s", serviceType))
	utilities.FilterDataCache.InvalidateByPrefix(fmt.Sprintf("filter_data:%s:", serviceType))
}

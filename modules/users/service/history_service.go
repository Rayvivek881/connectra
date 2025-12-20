package service

import (
	"context"
	"database/sql"
	"vivek-ray/models"
	"vivek-ray/modules/users/helper"
	"vivek-ray/modules/users/repository"
)

type HistoryService struct {
	historyRepo *repository.UserHistoryRepository
	userRepo    *repository.UserRepository
}

func NewHistoryService() *HistoryService {
	return &HistoryService{
		historyRepo: repository.NewUserHistoryRepository(),
		userRepo:    repository.NewUserRepository(),
	}
}

// RecordRegistration records a user registration event
func (s *HistoryService) RecordRegistration(ctx context.Context, userID string, geolocation *helper.GeolocationData) error {
	history := &models.UserHistory{
		UserID:    userID,
		EventType: models.EventTypeRegistration,
	}

	if geolocation != nil {
		s.populateHistoryFromGeolocation(history, geolocation)
	}

	return s.historyRepo.CreateHistory(ctx, history)
}

// RecordLogin records a user login event
func (s *HistoryService) RecordLogin(ctx context.Context, userID string, geolocation *helper.GeolocationData) error {
	history := &models.UserHistory{
		UserID:    userID,
		EventType: models.EventTypeLogin,
	}

	if geolocation != nil {
		s.populateHistoryFromGeolocation(history, geolocation)
	}

	return s.historyRepo.CreateHistory(ctx, history)
}

// GetUserHistory retrieves user history with pagination and filtering
func (s *HistoryService) GetUserHistory(
	ctx context.Context,
	userID *string,
	eventType *string,
	limit, offset int,
) (*helper.UserHistoryListResponse, error) {
	var eventTypeEnum *models.UserHistoryEventType
	if eventType != nil {
		et := models.UserHistoryEventType(*eventType)
		eventTypeEnum = &et
	}

	historyRecords, total, err := s.historyRepo.ListHistory(ctx, userID, eventTypeEnum, limit, offset)
	if err != nil {
		return nil, err
	}

	items := make([]helper.UserHistoryItem, 0, len(historyRecords))
	for _, record := range historyRecords {
		item := s.toHistoryItem(ctx, record)
		items = append(items, *item)
	}

	return &helper.UserHistoryListResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// populateHistoryFromGeolocation populates history record from geolocation data
func (s *HistoryService) populateHistoryFromGeolocation(history *models.UserHistory, geolocation *helper.GeolocationData) {
	history.IP = geolocation.IP
	history.Continent = geolocation.Continent
	history.ContinentCode = geolocation.ContinentCode
	history.Country = geolocation.Country
	history.CountryCode = geolocation.CountryCode
	history.Region = geolocation.Region
	history.RegionName = geolocation.RegionName
	history.City = geolocation.City
	history.District = geolocation.District
	history.Zip = geolocation.Zip
	history.Timezone = geolocation.Timezone
	history.Currency = geolocation.Currency
	history.ISP = geolocation.ISP
	history.Org = geolocation.Org
	history.ASName = geolocation.ASName
	history.Reverse = geolocation.Reverse
	history.Device = geolocation.Device
	history.Proxy = geolocation.Proxy
	history.Hosting = geolocation.Hosting

	if geolocation.Lat != nil {
		history.Lat = &sql.NullFloat64{
			Float64: *geolocation.Lat,
			Valid:   true,
		}
	}
	if geolocation.Lon != nil {
		history.Lon = &sql.NullFloat64{
			Float64: *geolocation.Lon,
			Valid:   true,
		}
	}
}

// toHistoryItem converts UserHistory model to UserHistoryItem response
func (s *HistoryService) toHistoryItem(ctx context.Context, record models.UserHistory) *helper.UserHistoryItem {
	item := &helper.UserHistoryItem{
		ID:            int(record.ID),
		UserID:        record.UserID,
		EventType:     string(record.EventType),
		IP:            record.IP,
		Continent:     record.Continent,
		ContinentCode: record.ContinentCode,
		Country:       record.Country,
		CountryCode:   record.CountryCode,
		Region:        record.Region,
		RegionName:    record.RegionName,
		City:          record.City,
		District:      record.District,
		Zip:           record.Zip,
		Timezone:      record.Timezone,
		Currency:      record.Currency,
		ISP:           record.ISP,
		Org:           record.Org,
		ASName:        record.ASName,
		Reverse:       record.Reverse,
		Device:        record.Device,
		Proxy:         record.Proxy,
		Hosting:       record.Hosting,
		CreatedAt:     record.CreatedAt,
	}

	// Extract lat/lon
	if record.Lat != nil && record.Lat.Valid {
		item.Lat = &record.Lat.Float64
	}
	if record.Lon != nil && record.Lon.Valid {
		item.Lon = &record.Lon.Float64
	}

	// Get user info if available
	if record.User != nil {
		item.UserEmail = &record.User.Email
		item.UserName = record.User.Name
	} else {
		// Try to fetch user
		user, err := s.userRepo.GetByUUID(ctx, record.UserID)
		if err == nil && user != nil {
			item.UserEmail = &user.Email
			item.UserName = user.Name
		}
	}

	return item
}


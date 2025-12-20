package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vivek-ray/models"
	"vivek-ray/modules/usage/helper"
	userrepo "vivek-ray/modules/users/repository"

	"github.com/rs/zerolog/log"
)

// Feature limits configuration matching Python backend
var FEATURE_LIMITS = map[models.FeatureType]map[models.Role]int{
	models.FeatureTypeAIChat:           {models.RoleFreeUser: 0, models.RoleProUser: -1},  // -1 = unlimited
	models.FeatureTypeBulkExport:       {models.RoleFreeUser: 0, models.RoleProUser: -1},
	models.FeatureTypeAPIKeys:          {models.RoleFreeUser: 0, models.RoleProUser: -1},
	models.FeatureTypeTeamManagement:   {models.RoleFreeUser: 0, models.RoleProUser: -1},
	models.FeatureTypeEmailFinder:      {models.RoleFreeUser: 10, models.RoleProUser: -1},
	models.FeatureTypeVerifier:         {models.RoleFreeUser: 5, models.RoleProUser: -1},
	models.FeatureTypeLinkedIn:         {models.RoleFreeUser: 5, models.RoleProUser: -1},
	models.FeatureTypeDataSearch:       {models.RoleFreeUser: 20, models.RoleProUser: -1},
	models.FeatureTypeAdvancedFilters:  {models.RoleFreeUser: 0, models.RoleProUser: -1},
	models.FeatureTypeAISummaries:      {models.RoleFreeUser: 0, models.RoleProUser: -1},
	models.FeatureTypeSaveSearches:     {models.RoleFreeUser: 0, models.RoleProUser: -1},
	models.FeatureTypeBulkVerification: {models.RoleFreeUser: 0, models.RoleProUser: -1},
}

// UsageService handles feature usage tracking business logic
type UsageService struct {
	featureRepo *userrepo.FeatureUsageRepository
	profileRepo *userrepo.UserProfileRepository
}

// NewUsageService creates a new usage service
func NewUsageService() *UsageService {
	return &UsageService{
		featureRepo: userrepo.NewFeatureUsageRepository(),
		profileRepo: userrepo.NewUserProfileRepository(),
	}
}

// getUserRoleLevel normalizes user role to FreeUser or ProUser
func (s *UsageService) getUserRoleLevel(userRole string) models.Role {
	if userRole == string(models.RoleSuperAdmin) || userRole == string(models.RoleAdmin) {
		return models.RoleProUser // Admins get pro-level access
	}
	if userRole == string(models.RoleProUser) {
		return models.RoleProUser
	}
	return models.RoleFreeUser
}

// getFeatureLimit returns usage limit for a feature based on user role
func (s *UsageService) getFeatureLimit(feature models.FeatureType, userRole string) int {
	roleLevel := s.getUserRoleLevel(userRole)
	limits, exists := FEATURE_LIMITS[feature]
	if !exists {
		return 0
	}
	limit, exists := limits[roleLevel]
	if !exists {
		return 0
	}
	// -1 means unlimited for pro users
	if limit == -1 && roleLevel == models.RoleProUser {
		return -1 // Unlimited
	}
	return limit
}

// calculatePeriodEnd calculates the end of the current month
func (s *UsageService) calculatePeriodEnd(now time.Time) time.Time {
	// Get first day of next month
	year, month, _ := now.Date()
	if month == 12 {
		return time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
}

// getOrCreateUsage gets existing usage record or creates a new one
func (s *UsageService) getOrCreateUsage(ctx context.Context, userID string, feature models.FeatureType, userRole string) (*models.FeatureUsage, error) {
	// Check if usage record exists
	usage, err := s.featureRepo.GetUsage(ctx, userID, feature)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	limit := s.getFeatureLimit(feature, userRole)

	if usage != nil {
		// Check if we need to reset based on billing period
		if usage.PeriodEnd != nil && now.After(*usage.PeriodEnd) {
			// Reset usage for new period
			periodEnd := s.calculatePeriodEnd(now)
			err := s.featureRepo.ResetUsageWithPeriod(ctx, userID, feature, now, periodEnd, limit)
			if err != nil {
				return nil, err
			}
			// Reload usage after reset
			usage, err = s.featureRepo.GetUsage(ctx, userID, feature)
			if err != nil {
				return nil, err
			}
		} else if usage.PeriodEnd == nil {
			// Initialize period_end if not set
			periodEnd := s.calculatePeriodEnd(now)
			usage.PeriodEnd = &periodEnd
			usage.Limit = limit
			err := s.featureRepo.CreateOrUpdateUsage(ctx, usage)
			if err != nil {
				return nil, err
			}
		} else if usage.Limit != limit {
			// Update limit if it changed (e.g., user upgraded)
			usage.Limit = limit
			err := s.featureRepo.CreateOrUpdateUsage(ctx, usage)
			if err != nil {
				return nil, err
			}
		}
		return usage, nil
	}

	// Create new usage record
	periodEnd := s.calculatePeriodEnd(now)
	usage = &models.FeatureUsage{
		UserID:      userID,
		Feature:     feature,
		Used:        0,
		Limit:       limit,
		PeriodStart: now,
		PeriodEnd:   &periodEnd,
	}

	err = s.featureRepo.CreateOrUpdateUsage(ctx, usage)
	if err != nil {
		return nil, err
	}

	return usage, nil
}

// GetCurrentUsage returns current usage for all features for a user
func (s *UsageService) GetCurrentUsage(ctx context.Context, userID string) (helper.CurrentUsageResponse, error) {
	// Get user profile to determine role
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New(fmt.Sprintf("User profile not found for user_id: %s", userID))
	}

	userRole := "FreeUser"
	if profile.Role != nil {
		userRole = *profile.Role
	}

	// Get all existing usage records
	usages, err := s.featureRepo.GetAllUsage(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create a map of existing usages
	usageMap := make(map[models.FeatureType]*models.FeatureUsage)
	for _, usage := range usages {
		usageMap[usage.Feature] = usage
	}

	// Build response with all features
	response := make(helper.CurrentUsageResponse)
	now := time.Now()

	// Iterate through all feature types
	allFeatures := []models.FeatureType{
		models.FeatureTypeAIChat,
		models.FeatureTypeBulkExport,
		models.FeatureTypeAPIKeys,
		models.FeatureTypeTeamManagement,
		models.FeatureTypeEmailFinder,
		models.FeatureTypeVerifier,
		models.FeatureTypeLinkedIn,
		models.FeatureTypeDataSearch,
		models.FeatureTypeAdvancedFilters,
		models.FeatureTypeAISummaries,
		models.FeatureTypeSaveSearches,
		models.FeatureTypeBulkVerification,
	}

	for _, featureType := range allFeatures {
		usage := usageMap[featureType]

		if usage != nil {
			// Check if period needs reset
			if usage.PeriodEnd != nil && now.After(*usage.PeriodEnd) {
				// Reset usage for new period
				limit := s.getFeatureLimit(featureType, userRole)
				periodEnd := s.calculatePeriodEnd(now)
				err := s.featureRepo.ResetUsageWithPeriod(ctx, userID, featureType, now, periodEnd, limit)
				if err != nil {
					log.Error().Err(err).Msg("Error resetting usage period")
					continue
				}
				// Reload usage after reset
				usage, err = s.featureRepo.GetUsage(ctx, userID, featureType)
				if err != nil {
					log.Error().Err(err).Msg("Error reloading usage after reset")
					continue
				}
			}

			limit := usage.Limit
			if limit == -1 {
				limit = 999999 // Return 999999 for unlimited
			}
			response[string(featureType)] = helper.FeatureUsageItem{
				Used:  usage.Used,
				Limit: limit,
			}
		} else {
			// Create default entry for feature
			limit := s.getFeatureLimit(featureType, userRole)
			if limit == -1 {
				limit = 999999 // Return 999999 for unlimited
			}
			response[string(featureType)] = helper.FeatureUsageItem{
				Used:  0,
				Limit: limit,
			}
		}
	}

	return response, nil
}

// TrackUsage tracks feature usage for a user
func (s *UsageService) TrackUsage(ctx context.Context, userID string, featureStr string, amount int) (*helper.TrackUsageResponse, error) {
	// Validate feature name
	var featureType models.FeatureType
	switch featureStr {
	case string(models.FeatureTypeAIChat):
		featureType = models.FeatureTypeAIChat
	case string(models.FeatureTypeBulkExport):
		featureType = models.FeatureTypeBulkExport
	case string(models.FeatureTypeAPIKeys):
		featureType = models.FeatureTypeAPIKeys
	case string(models.FeatureTypeTeamManagement):
		featureType = models.FeatureTypeTeamManagement
	case string(models.FeatureTypeEmailFinder):
		featureType = models.FeatureTypeEmailFinder
	case string(models.FeatureTypeVerifier):
		featureType = models.FeatureTypeVerifier
	case string(models.FeatureTypeLinkedIn):
		featureType = models.FeatureTypeLinkedIn
	case string(models.FeatureTypeDataSearch):
		featureType = models.FeatureTypeDataSearch
	case string(models.FeatureTypeAdvancedFilters):
		featureType = models.FeatureTypeAdvancedFilters
	case string(models.FeatureTypeAISummaries):
		featureType = models.FeatureTypeAISummaries
	case string(models.FeatureTypeSaveSearches):
		featureType = models.FeatureTypeSaveSearches
	case string(models.FeatureTypeBulkVerification):
		featureType = models.FeatureTypeBulkVerification
	default:
		return nil, errors.New(fmt.Sprintf("Invalid feature: %s", featureStr))
	}

	// Get user profile to determine role
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New(fmt.Sprintf("User profile not found for user_id: %s", userID))
	}

	userRole := "FreeUser"
	if profile.Role != nil {
		userRole = *profile.Role
	}

	// Get or create usage record
	usage, err := s.getOrCreateUsage(ctx, userID, featureType, userRole)
	if err != nil {
		return nil, err
	}

	// Update limit if it changed (e.g., user upgraded)
	limit := s.getFeatureLimit(featureType, userRole)
	if usage.Limit != limit {
		usage.Limit = limit
	}

	// Increment usage (only if not unlimited)
	if usage.Limit == -1 || usage.Limit == 0 {
		// Unlimited or no access - keep used at 0
		usage.Used = 0
	} else {
		// Cap at limit
		newUsed := usage.Used + amount
		if newUsed > usage.Limit {
			usage.Used = usage.Limit
		} else {
			usage.Used = newUsed
		}
	}

	now := time.Now()
	usage.UpdatedAt = &now

	err = s.featureRepo.CreateOrUpdateUsage(ctx, usage)
	if err != nil {
		return nil, err
	}

	// Return response with limit as 999999 if unlimited
	responseLimit := usage.Limit
	if responseLimit == -1 {
		responseLimit = 999999
	}

	return &helper.TrackUsageResponse{
		Feature: featureStr,
		Used:    usage.Used,
		Limit:   responseLimit,
		Success: true,
	}, nil
}


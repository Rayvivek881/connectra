package service

import (
	"context"
	"errors"
	"time"
	"vivek-ray/models"
	"vivek-ray/modules/users/helper"
	"vivek-ray/modules/users/repository"
	"vivek-ray/utilities"
)

type UserService struct {
	userRepo    *repository.UserRepository
	profileRepo *repository.UserProfileRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:    repository.NewUserRepository(),
		profileRepo: repository.NewUserProfileRepository(),
	}
}

// GetUserProfile retrieves user profile, creating one if it doesn't exist
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*helper.ProfileResponse, error) {
	user, err := s.userRepo.GetByUUID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Auto-create profile if it doesn't exist
	if profile == nil {
		defaultProfile := &models.UserProfile{
			UserID:             userID,
			Role:               stringPtr("FreeUser"),
			Credits:            50,
			SubscriptionPlan:   stringPtr("free"),
			SubscriptionStatus: stringPtr("active"),
			Notifications: &models.NotificationPreferences{
				WeeklyReports: boolPtr(true),
				NewLeadAlerts: boolPtr(true),
			},
		}
		if err := s.profileRepo.CreateProfile(ctx, defaultProfile); err != nil {
			return nil, err
		}
		profile = defaultProfile
	}

	return s.toProfileResponse(user, profile), nil
}

// UpdateUserProfile updates user profile with partial update
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, update *helper.ProfileUpdateRequest) (*helper.ProfileResponse, error) {
	user, err := s.userRepo.GetByUUID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update user name if provided
	if update.Name != nil {
		user.Name = update.Name
		if err := s.userRepo.UpdateUser(ctx, user); err != nil {
			return nil, err
		}
	}

	// Get or create profile
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		defaultProfile := &models.UserProfile{
			UserID:             userID,
			Role:               stringPtr("FreeUser"),
			Credits:            50,
			SubscriptionPlan:   stringPtr("free"),
			SubscriptionStatus: stringPtr("active"),
			Notifications: &models.NotificationPreferences{
				WeeklyReports: boolPtr(true),
				NewLeadAlerts: boolPtr(true),
			},
		}
		if err := s.profileRepo.CreateProfile(ctx, defaultProfile); err != nil {
			return nil, err
		}
		profile = defaultProfile
	}

	// Update profile fields
	now := time.Now()
	profile.UpdatedAt = &now

	if update.JobTitle != nil {
		profile.JobTitle = update.JobTitle
	}
	if update.Bio != nil {
		profile.Bio = update.Bio
	}
	if update.Timezone != nil {
		profile.Timezone = update.Timezone
	}
	if update.AvatarURL != nil {
		profile.AvatarURL = update.AvatarURL
	}
	if update.Role != nil {
		profile.Role = update.Role
	}

	// Merge notifications
	if update.Notifications != nil {
		if profile.Notifications == nil {
			profile.Notifications = &models.NotificationPreferences{}
		}
		if update.Notifications.WeeklyReports != nil {
			profile.Notifications.WeeklyReports = update.Notifications.WeeklyReports
		}
		if update.Notifications.NewLeadAlerts != nil {
			profile.Notifications.NewLeadAlerts = update.Notifications.NewLeadAlerts
		}
	}

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.toProfileResponse(user, profile), nil
}

// PromoteUserToAdmin promotes a user to admin role
func (s *UserService) PromoteUserToAdmin(ctx context.Context, userID string) (*helper.ProfileResponse, error) {
	user, err := s.userRepo.GetByUUID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	profile, err := s.profileRepo.GetOrCreate(ctx, userID, &models.UserProfile{
		Role:               stringPtr("FreeUser"),
		Credits:            50,
		SubscriptionPlan:   stringPtr("free"),
		SubscriptionStatus: stringPtr("active"),
		Notifications: &models.NotificationPreferences{
			WeeklyReports: boolPtr(true),
			NewLeadAlerts: boolPtr(true),
		},
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	profile.Role = stringPtr("Admin")
	profile.UpdatedAt = &now

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.toProfileResponse(user, profile), nil
}

// PromoteUserToSuperAdmin promotes a user to super admin role
func (s *UserService) PromoteUserToSuperAdmin(ctx context.Context, userID string) (*helper.ProfileResponse, error) {
	user, err := s.userRepo.GetByUUID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	profile, err := s.profileRepo.GetOrCreate(ctx, userID, &models.UserProfile{
		Role:               stringPtr("FreeUser"),
		Credits:            50,
		SubscriptionPlan:   stringPtr("free"),
		SubscriptionStatus: stringPtr("active"),
		Notifications: &models.NotificationPreferences{
			WeeklyReports: boolPtr(true),
			NewLeadAlerts: boolPtr(true),
		},
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	profile.Role = stringPtr("SuperAdmin")
	profile.UpdatedAt = &now

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.toProfileResponse(user, profile), nil
}

// GetUserStats returns aggregated user statistics
func (s *UserService) GetUserStats(ctx context.Context) (*helper.UserStatsResponse, error) {
	// This is a simplified version - implement actual aggregation queries
	users, total, err := s.userRepo.ListAllUsers(ctx, 10000, 0)
	if err != nil {
		return nil, err
	}

	activeUsers := 0
	usersByRole := make(map[string]int)
	usersByPlan := make(map[string]int)

	for _, user := range users {
		if user.IsActive {
			activeUsers++
		}

		profile, _ := s.profileRepo.GetByUserID(ctx, user.UUID)
		if profile != nil {
			if profile.Role != nil {
				usersByRole[*profile.Role]++
			}
			if profile.SubscriptionPlan != nil {
				usersByPlan[*profile.SubscriptionPlan]++
			}
		}
	}

	return &helper.UserStatsResponse{
		TotalUsers:  total,
		ActiveUsers: activeUsers,
		UsersByRole: usersByRole,
		UsersByPlan: usersByPlan,
	}, nil
}

// toProfileResponse converts User and UserProfile to ProfileResponse
func (s *UserService) toProfileResponse(user *models.User, profile *models.UserProfile) *helper.ProfileResponse {
	avatarURL := profile.AvatarURL
	if avatarURL != nil {
		fullURL := utilities.GetFullAvatarURL(avatarURL)
		avatarURL = fullURL
	}

	var notifications *helper.NotificationPreferences
	if profile.Notifications != nil {
		notifications = &helper.NotificationPreferences{
			WeeklyReports: profile.Notifications.WeeklyReports,
			NewLeadAlerts: profile.Notifications.NewLeadAlerts,
		}
	}

	return &helper.ProfileResponse{
		UUID:          user.UUID,
		Name:          user.Name,
		Email:         user.Email,
		Role:          profile.Role,
		AvatarURL:     avatarURL,
		IsActive:      user.IsActive,
		JobTitle:      profile.JobTitle,
		Bio:           profile.Bio,
		Timezone:      profile.Timezone,
		Notifications: notifications,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     profile.UpdatedAt,
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}


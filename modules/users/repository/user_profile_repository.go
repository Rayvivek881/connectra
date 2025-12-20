package repository

import (
	"context"
	"vivek-ray/connections"
	"vivek-ray/models"

	"github.com/rs/zerolog/log"
)

type UserProfileRepository struct{}

func NewUserProfileRepository() *UserProfileRepository {
	return &UserProfileRepository{}
}

// CreateProfile creates a new user profile
func (r *UserProfileRepository) CreateProfile(ctx context.Context, profile *models.UserProfile) error {
	_, err := connections.PgDBConnection.Client.NewInsert().Model(profile).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user profile")
		return err
	}
	return nil
}

// GetByUserID retrieves a profile by user ID
func (r *UserProfileRepository) GetByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
	profile := new(models.UserProfile)
	err := connections.PgDBConnection.Client.NewSelect().
		Model(profile).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error finding profile by user ID")
		return nil, err
	}
	return profile, nil
}

// GetOrCreate retrieves a profile or creates one if it doesn't exist
func (r *UserProfileRepository) GetOrCreate(ctx context.Context, userID string, defaults *models.UserProfile) (*models.UserProfile, error) {
	profile, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	if profile != nil {
		return profile, nil
	}
	
	// Create with defaults
	if defaults == nil {
		defaults = &models.UserProfile{
			UserID: userID,
		}
	} else {
		defaults.UserID = userID
	}
	
	if err := r.CreateProfile(ctx, defaults); err != nil {
		return nil, err
	}
	
	return defaults, nil
}

// UpdateProfile updates profile fields
func (r *UserProfileRepository) UpdateProfile(ctx context.Context, profile *models.UserProfile) error {
	_, err := connections.PgDBConnection.Client.NewUpdate().
		Model(profile).
		Where("user_id = ?", profile.UserID).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error updating user profile")
		return err
	}
	return nil
}


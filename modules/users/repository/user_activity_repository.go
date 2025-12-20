package repository

import (
	"context"
	"vivek-ray/connections"
	"vivek-ray/models"

	"github.com/rs/zerolog/log"
)

type UserActivityRepository struct{}

func NewUserActivityRepository() *UserActivityRepository {
	return &UserActivityRepository{}
}

// CreateActivity creates a new activity record
func (r *UserActivityRepository) CreateActivity(ctx context.Context, activity *models.UserActivity) error {
	_, err := connections.PgDBConnection.Client.NewInsert().Model(activity).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user activity")
		return err
	}
	return nil
}

// ListByUser returns paginated activities for a user
func (r *UserActivityRepository) ListByUser(
	ctx context.Context,
	userID string,
	limit, offset int,
) ([]models.UserActivity, int, error) {
	var activities []models.UserActivity
	
	// Get total count
	count, err := connections.PgDBConnection.Client.NewSelect().
		Model((*models.UserActivity)(nil)).
		Where("user_id = ?", userID).
		Count(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error counting user activities")
		return nil, 0, err
	}
	
	// Get paginated results
	err = connections.PgDBConnection.Client.NewSelect().
		Model(&activities).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing user activities")
		return nil, 0, err
	}
	
	return activities, count, nil
}

// GetStats returns aggregated statistics for user activities
func (r *UserActivityRepository) GetStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	// This is a placeholder - implement actual aggregation queries as needed
	stats := make(map[string]interface{})
	
	// Example: Get total count
	total, err := connections.PgDBConnection.Client.NewSelect().
		Model((*models.UserActivity)(nil)).
		Where("user_id = ?", userID).
		Count(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error getting activity stats")
		return nil, err
	}
	
	stats["total_activities"] = total
	return stats, nil
}


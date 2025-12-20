package repository

import (
	"context"
	"time"
	"vivek-ray/connections"
	"vivek-ray/models"

	"github.com/rs/zerolog/log"
)

type FeatureUsageRepository struct{}

func NewFeatureUsageRepository() *FeatureUsageRepository {
	return &FeatureUsageRepository{}
}

// GetUsage retrieves feature usage for a user
func (r *FeatureUsageRepository) GetUsage(ctx context.Context, userID string, feature models.FeatureType) (*models.FeatureUsage, error) {
	usage := new(models.FeatureUsage)
	err := connections.PgDBConnection.Client.NewSelect().
		Model(usage).
		Where("user_id = ? AND feature = ?", userID, feature).
		Scan(ctx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error finding feature usage")
		return nil, err
	}
	return usage, nil
}

// GetAllUsage retrieves all feature usage records for a user
func (r *FeatureUsageRepository) GetAllUsage(ctx context.Context, userID string) ([]*models.FeatureUsage, error) {
	var usages []*models.FeatureUsage
	err := connections.PgDBConnection.Client.NewSelect().
		Model(&usages).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error finding all feature usage")
		return nil, err
	}
	return usages, nil
}

// CreateOrUpdateUsage creates or updates feature usage
func (r *FeatureUsageRepository) CreateOrUpdateUsage(ctx context.Context, usage *models.FeatureUsage) error {
	now := time.Now()
	usage.UpdatedAt = &now
	
	_, err := connections.PgDBConnection.Client.NewInsert().
		Model(usage).
		On("CONFLICT (user_id, feature) DO UPDATE").
		Set("used = EXCLUDED.used").
		Set("\"limit\" = EXCLUDED.limit").
		Set("period_start = EXCLUDED.period_start").
		Set("period_end = EXCLUDED.period_end").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating/updating feature usage")
		return err
	}
	return nil
}

// IncrementUsage increments the used count for a feature
func (r *FeatureUsageRepository) IncrementUsage(ctx context.Context, userID string, feature models.FeatureType) error {
	_, err := connections.PgDBConnection.Client.NewUpdate().
		Model((*models.FeatureUsage)(nil)).
		Set("used = used + 1").
		Set("updated_at = NOW()").
		Where("user_id = ? AND feature = ?", userID, feature).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error incrementing feature usage")
		return err
	}
	return nil
}

// IncrementUsageByAmount increments the used count by a custom amount
func (r *FeatureUsageRepository) IncrementUsageByAmount(ctx context.Context, userID string, feature models.FeatureType, amount int) error {
	_, err := connections.PgDBConnection.Client.NewUpdate().
		Model((*models.FeatureUsage)(nil)).
		Set("used = used + ?", amount).
		Set("updated_at = NOW()").
		Where("user_id = ? AND feature = ?", userID, feature).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error incrementing feature usage by amount")
		return err
	}
	return nil
}

// ResetUsageWithPeriod resets usage with new period dates and limit
func (r *FeatureUsageRepository) ResetUsageWithPeriod(ctx context.Context, userID string, feature models.FeatureType, newPeriodStart, newPeriodEnd time.Time, newLimit int) error {
	_, err := connections.PgDBConnection.Client.NewUpdate().
		Model((*models.FeatureUsage)(nil)).
		Set("used = 0").
		Set("\"limit\" = ?", newLimit).
		Set("period_start = ?", newPeriodStart).
		Set("period_end = ?", newPeriodEnd).
		Set("updated_at = NOW()").
		Where("user_id = ? AND feature = ?", userID, feature).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error resetting feature usage with period")
		return err
	}
	return nil
}

// CheckLimit checks if usage is within limit
func (r *FeatureUsageRepository) CheckLimit(ctx context.Context, userID string, feature models.FeatureType) (bool, error) {
	usage, err := r.GetUsage(ctx, userID, feature)
	if err != nil {
		return false, err
	}
	if usage == nil {
		return true, nil // No limit set
	}
	// If limit is -1 (unlimited) or 0 (no access), always return true for unlimited
	if usage.Limit == -1 || usage.Limit == 0 {
		return true, nil
	}
	return usage.Used < usage.Limit, nil
}

// ResetUsage resets usage for a new period
func (r *FeatureUsageRepository) ResetUsage(ctx context.Context, userID string, feature models.FeatureType) error {
	_, err := connections.PgDBConnection.Client.NewUpdate().
		Model((*models.FeatureUsage)(nil)).
		Set("used = 0").
		Set("period_start = NOW()").
		Set("updated_at = NOW()").
		Where("user_id = ? AND feature = ?", userID, feature).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error resetting feature usage")
		return err
	}
	return nil
}


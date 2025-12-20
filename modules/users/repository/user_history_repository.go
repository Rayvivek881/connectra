package repository

import (
	"context"
	"vivek-ray/connections"
	"vivek-ray/models"

	"github.com/rs/zerolog/log"
)

type UserHistoryRepository struct{}

func NewUserHistoryRepository() *UserHistoryRepository {
	return &UserHistoryRepository{}
}

// CreateHistory creates a new history record
func (r *UserHistoryRepository) CreateHistory(ctx context.Context, history *models.UserHistory) error {
	_, err := connections.PgDBConnection.Client.NewInsert().Model(history).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user history")
		return err
	}
	return nil
}

// ListHistory returns paginated history records with optional filtering
func (r *UserHistoryRepository) ListHistory(
	ctx context.Context,
	userID *string,
	eventType *models.UserHistoryEventType,
	limit, offset int,
) ([]models.UserHistory, int, error) {
	var history []models.UserHistory
	
	query := connections.PgDBConnection.Client.NewSelect().Model(&history)
	
	// Apply filters
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if eventType != nil {
		query = query.Where("event_type = ?", *eventType)
	}
	
	// Get total count
	count, err := query.Count(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error counting user history")
		return nil, 0, err
	}
	
	// Get paginated results
	err = query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing user history")
		return nil, 0, err
	}
	
	return history, count, nil
}


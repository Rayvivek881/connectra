package repository

import (
	"context"
	"vivek-ray/connections"
	"vivek-ray/models"

	"github.com/rs/zerolog/log"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := connections.PgDBConnection.Client.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user")
		return err
	}
	return nil
}

// GetByUUID retrieves a user by UUID
func (r *UserRepository) GetByUUID(ctx context.Context, uuid string) (*models.User, error) {
	user := new(models.User)
	err := connections.PgDBConnection.Client.NewSelect().
		Model(user).
		Where("uuid = ?", uuid).
		Scan(ctx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error finding user by UUID")
		return nil, err
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	err := connections.PgDBConnection.Client.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error finding user by email")
		return nil, err
	}
	return user, nil
}

// UpdateUser updates user fields
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := connections.PgDBConnection.Client.NewUpdate().
		Model(user).
		Where("uuid = ?", user.UUID).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error updating user")
		return err
	}
	return nil
}

// DeleteUser deletes a user by UUID
func (r *UserRepository) DeleteUser(ctx context.Context, uuid string) error {
	_, err := connections.PgDBConnection.Client.NewDelete().
		Model((*models.User)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting user")
		return err
	}
	return nil
}

// ListAllUsers returns a paginated list of all users
func (r *UserRepository) ListAllUsers(ctx context.Context, limit, offset int) ([]models.User, int, error) {
	var users []models.User
	
	// Get total count
	count, err := connections.PgDBConnection.Client.NewSelect().
		Model((*models.User)(nil)).
		Count(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error counting users")
		return nil, 0, err
	}
	
	// Get paginated users
	err = connections.PgDBConnection.Client.NewSelect().
		Model(&users).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing users")
		return nil, 0, err
	}
	
	return users, count, nil
}


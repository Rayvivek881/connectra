package repository

import (
	"context"
	"time"
	"vivek-ray/connections"
	"vivek-ray/models"

	"github.com/rs/zerolog/log"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := connections.PgDBConnection.Client.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user")
		return err
	}
	return nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	err := connections.PgDBConnection.Client.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error finding user by email")
		return nil, err
	}
	return user, nil
}

func (r *AuthRepository) AddToBlacklist(ctx context.Context, token string, expiresAt time.Time) error {
	blacklist := &models.TokenBlacklist{
		Token:     token,
		ExpiresAt: expiresAt,
	}
	_, err := connections.PgDBConnection.Client.NewInsert().Model(blacklist).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error adding token to blacklist")
		return err
	}
	return nil
}

func (r *AuthRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := connections.PgDBConnection.Client.NewSelect().
		Model((*models.TokenBlacklist)(nil)).
		Where("token = ?", token).
		Exists(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error checking if token is blacklisted")
		return false, err
	}
	return exists, nil
}

// EnsureTables creates the necessary tables if they don't exist.
// This is a helper for setup; in production, use migrations.
func (r *AuthRepository) EnsureTables(ctx context.Context) error {
	db := connections.PgDBConnection.Client

	log.Info().Msg("Ensuring auth tables exist...")

	// Drop tables for development/testing to ensure schema is correct
	// In production, use proper migrations
	if _, err := db.ExecContext(ctx, "DROP TABLE IF EXISTS users CASCADE"); err != nil {
		log.Error().Err(err).Msg("Error dropping users table")
	}
	if _, err := db.ExecContext(ctx, "DROP TABLE IF EXISTS token_blacklist CASCADE"); err != nil {
		log.Error().Err(err).Msg("Error dropping token_blacklist table")
	}

	log.Info().Msg("Creating users table...")
	_, err := db.NewCreateTable().Model((*models.User)(nil)).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating users table")
		return err
	}

	log.Info().Msg("Creating token_blacklist table...")
	_, err = db.NewCreateTable().Model((*models.TokenBlacklist)(nil)).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error creating token_blacklist table")
		return err
	}

	log.Info().Msg("Auth tables created successfully")
	return nil
}

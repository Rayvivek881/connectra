package utilities

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

// WithTransaction executes a function within a database transaction.
// It automatically handles rollback on error and commit on success.
//
// Example:
//   result, err := WithTransaction(ctx, db, func(tx bun.Tx) (*MyModel, error) {
//       // Perform operations with tx
//       return &MyModel{}, nil
//   })
func WithTransaction[T any](
	ctx context.Context,
	db *bun.DB,
	fn func(tx bun.Tx) (T, error),
) (T, error) {
	var zero T

	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return zero, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure rollback on error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Log rollback error but return original error
				err = fmt.Errorf("%w (rollback error: %v)", err, rbErr)
			}
		}
	}()

	// Execute function
	result, err := fn(tx)
	if err != nil {
		return zero, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return zero, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// WithTransactionVoid executes a function within a database transaction that returns no value.
// It automatically handles rollback on error and commit on success.
//
// Example:
//   err := WithTransactionVoid(ctx, db, func(tx bun.Tx) error {
//       // Perform operations with tx
//       return nil
//   })
func WithTransactionVoid(
	ctx context.Context,
	db *bun.DB,
	fn func(tx bun.Tx) error,
) error {
	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure rollback on error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Log rollback error but return original error
				err = fmt.Errorf("%w (rollback error: %v)", err, rbErr)
			}
		}
	}()

	// Execute function
	if err = fn(tx); err != nil {
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}


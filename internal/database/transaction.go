package database

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Transaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
func Transaction(ctx context.Context, db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// channel to monitor context cancellation
	done := make(chan struct{})
	defer close(done)

	// goroutine to monitor context cancellation and rollback transaction if needed
	go func() {
		select {
		case <-ctx.Done():
			if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
				log.Ctx(ctx).
					Error().
					Err(rollbackErr).
					Msg("failed to rollback transaction on context cancellation")
			}
		case <-done:
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Ctx(ctx).Error().Interface("panic", r).Msg("transaction panicked, rolling back")
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			log.Ctx(ctx).Error().
				Err(rollbackErr).
				Msg("failed to rollback transaction")
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.Ctx(ctx).Error().
			Err(err).
			Msg("failed to commit transaction")
		return err
	}

	return nil
}

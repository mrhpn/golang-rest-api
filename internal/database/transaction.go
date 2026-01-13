package database

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// Transaction executes the given function within a database transaction.
//
// Behavior:
//   - Begins a transaction
//   - Commits if fn returns nil
//   - Rolls back if fn returns an error
//   - Rolls back if fn panics
//   - Propagates context cancellation to all queries
//
// This is a thin wrapper around gorm.DB.Transaction.
func Transaction(
	ctx context.Context,
	db *gorm.DB,
	fn func(context.Context) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// inject the transaction in the context
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

// GetTx retrieves a transaction from the context if it exists
func GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// Usage example in service layer
// func (s *userService) UpdateProfile(ctx context.Context, id string) error {
// 	return database.Transaction(ctx, s.db, func(txCtx context.Context) error {
// 		// Both of these will automatically use the SAME transaction
// 		// because r.DB(txCtx) detects the transaction in the context!
// 		user, err := s.repo.FindByID(txCtx, id)
// 		if err != nil {
// 			return err
// 		}

// 		user.Status = "Active"
// 		return s.repo.Update(txCtx, user)
// 	})
// }

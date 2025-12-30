package users

// This file demonstrates how to use database transactions in service layer
// This is an EXAMPLE - not actual production code

import (
	"context"

	"github.com/mrhpn/go-rest-api/internal/database"
	"gorm.io/gorm"
)

// Example: TransferUserRole demonstrates using transactions
// when you need to perform multiple database operations atomically
//
// To use this pattern, you need to pass the *gorm.DB from AppContext to your service:
//
//	type service struct {
//		repo Repository
//		db   *gorm.DB  // Add this field
//	}
//
//	func NewService(repo Repository, db *gorm.DB) Service {
//		return &service{repo: repo, db: db}
//	}
//
// Then in your service method:
func ExampleTransferUserRole(ctx context.Context, db *gorm.DB, fromUserID, toUserID string, newRole string) error {
	// Example: Update multiple users in a transaction
	// If any operation fails, all changes are rolled back
	return database.Transaction(ctx, db, func(tx *gorm.DB) error {
		// Operation 1: Update source user
		if err := tx.WithContext(ctx).Model(&User{}).
			Where("id = ?", fromUserID).
			Update("role", "employee").Error; err != nil {
			return err // This will trigger rollback
		}

		// Operation 2: Update target user
		if err := tx.WithContext(ctx).Model(&User{}).
			Where("id = ?", toUserID).
			Update("role", newRole).Error; err != nil {
			return err // This will trigger rollback
		}

		// If we reach here, both operations succeeded
		// Transaction will be committed automatically
		return nil
	})
}

// Example: CreateUserWithProfile demonstrates using transactions
// when creating related records across multiple tables
func ExampleCreateUserWithProfile(ctx context.Context, db *gorm.DB, user *User, profileData map[string]interface{}) error {
	// This is just an example - you would need a Profile model and repository

	return database.Transaction(ctx, db, func(tx *gorm.DB) error {
		// Step 1: Create user
		if err := tx.WithContext(ctx).Create(user).Error; err != nil {
			return err // Rollback if user creation fails
		}

		// Step 2: Create user profile (example - you'd have a Profile model)
		// profile := &Profile{
		// 	UserID: user.ID,
		// 	Data:   profileData,
		// }
		// if err := tx.WithContext(ctx).Create(profile).Error; err != nil {
		// 	return err // Rollback if profile creation fails
		// }

		// All operations succeeded - transaction will commit
		return nil
	})
}

// Example: Using transactions in a service method
// This shows how you would integrate it into your existing service
//
// In your routes.go, pass db to service:
//   userS := users.NewService(userR, ctx.DB)
//
// In your service:
//   type service struct {
//       repo Repository
//       db   *gorm.DB
//   }
//
//   func (s *service) ComplexOperation(ctx context.Context, ...) error {
//       return database.Transaction(ctx, s.db, func(tx *gorm.DB) error {
//           // All your database operations using tx instead of s.repo
//           // If any operation fails, everything rolls back
//           return nil
//       })
//   }

// When to use transactions:
// 1. Multiple related database operations that must succeed or fail together
// 2. Data consistency requirements (e.g., transfer operations)
// 3. Creating related records across multiple tables
// 4. Complex business logic requiring atomicity
//
// When NOT to use transactions:
// 1. Single database operation (already atomic)
// 2. Read-only operations
// 3. Operations that can be safely retried individually

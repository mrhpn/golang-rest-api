# Answers to Your Questions

## Question 1: How is Query Timeout Applied to Repository Queries?

### Answer

The query timeout is **automatically applied to all database operations** using
GORM callbacks. This means you don't need to manually wrap each query - it's
handled automatically at the database connection level.

### Implementation

1. **Automatic via GORM Callbacks**: When the database connection is established
   in `database.Connect()`, GORM callbacks are registered that automatically
   apply query timeout to all operations:

   ```go
   // In database/postgres.go
   registerQueryTimeoutCallback(db, dbCfg)
   ```

2. **Clean Repository Code**: Repository methods are simple and clean - no
   timeout wrappers needed:

   ```go
   func (r *repository) FindById(ctx context.Context, id string) (*User, error) {
       var user User
       err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
       // Query timeout is automatically applied by GORM callback
       // ... error handling
   }
   ```

3. **How It Works**:
   - GORM callbacks intercept all database operations (Query, Create, Update,
     Delete)
   - Before each operation executes, the callback checks if the context has a
     timeout
   - If no timeout exists, it automatically applies the configured
     `DB_QUERY_TIMEOUT_SECOND`
   - If the query takes longer than the timeout, the context is cancelled
   - GORM respects context cancellation and returns an error

### Configuration

Set the query timeout via environment variable:

```bash
DB_QUERY_TIMEOUT_SECOND=30  # Default: 30 seconds
```

### Files Modified

- `internal/database/postgres.go` - Query timeout callbacks registered
  automatically
- `internal/modules/users/repository.go` - Clean code, timeout applied
  automatically

---

## Question 2: How to Use the Transaction Helper?

### Answer

The `database.Transaction()` helper is used when you need to perform multiple
database operations atomically (all succeed or all fail).

### When to Use Transactions

✅ **Use transactions when:**

- Multiple related database operations must succeed or fail together
- Data consistency is critical (e.g., transfer operations)
- Creating related records across multiple tables
- Complex business logic requiring atomicity

❌ **Don't use transactions when:**

- Single database operation (already atomic)
- Read-only operations
- Operations that can be safely retried individually

### How to Use

#### Step 1: Pass DB to Your Service

In `internal/routes/routes.go`, update service creation:

```go
// Current (without transactions):
userS := users.NewService(userR)

// With transactions support:
userS := users.NewService(userR, ctx.DB)
```

#### Step 2: Update Service to Accept DB

In your service struct:

```go
type service struct {
    repo Repository
    db   *gorm.DB  // Add this field
}

func NewService(repo Repository, db *gorm.DB) Service {
    return &service{repo: repo, db: db}
}
```

#### Step 3: Use Transaction in Service Method

```go
func (s *service) TransferUserRole(ctx context.Context, fromUserID, toUserID string, newRole string) error {
    return database.Transaction(ctx, s.db, func(tx *gorm.DB) error {
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
```

### Example Use Cases

1. **Transfer Operations**: Transfer data between users/accounts
2. **Multi-Table Inserts**: Create user + profile + settings in one transaction
3. **Complex Updates**: Update multiple related records atomically
4. **Financial Operations**: Debit/credit operations that must be atomic

### Important Notes

- **Use `tx` instead of `s.db`**: Inside the transaction function, use the
  `tx *gorm.DB` parameter, not the original `s.db`
- **Automatic Rollback**: If the function returns an error, the transaction is
  automatically rolled back
- **Automatic Commit**: If the function returns `nil`, the transaction is
  automatically committed
- **Panic Recovery**: The transaction helper includes panic recovery with
  rollback

### Files

- `internal/database/transaction.go` - Transaction helper implementation
- `internal/modules/users/transaction_example.go` - Detailed examples and
  patterns

---

## Question 3: Should We Check MinIO Health?

### Answer

**Yes!** MinIO health check has been added to the readiness probe since it's a
critical dependency for media uploads.

### Implementation

1. **Added HealthCheck Method to Media Service Interface**:

   ```go
   type Service interface {
       Upload(file *multipart.FileHeader, subDir FileCategory) (string, error)
       HealthCheck(ctx context.Context) error // New method
   }
   ```

2. **Implemented in MinIO Service**:

   ```go
   func (s *minioService) HealthCheck(ctx context.Context) error {
       ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
       defer cancel()

       // Check if bucket exists and is accessible
       exists, err := s.client.BucketExists(ctx, s.bucketName)
       if err != nil {
           return fmt.Errorf("failed to check bucket existence: %w", err)
       }

       if !exists {
           return fmt.Errorf("bucket '%s' does not exist", s.bucketName)
       }

       return nil
   }
   ```

3. **Integrated into Readiness Probe**:
   ```go
   // Check MinIO storage connectivity
   if h.appCtx.MediaService != nil {
       if err := h.appCtx.MediaService.HealthCheck(ctx); err != nil {
           checks["storage"] = "unhealthy: " + err.Error()
           allHealthy = false
       } else {
           checks["storage"] = "healthy"
       }
   }
   ```

### Health Check Response

The `/health/ready` endpoint now includes storage status:

```json
{
  "status": "ready",
  "timestamp": "2024-01-01T12:00:00Z",
  "checks": {
    "database": "healthy",
    "db_open_conns": "5",
    "db_idle_conns": "3",
    "storage": "healthy"
  }
}
```

### Why This Matters

- **Kubernetes/Docker**: Readiness probe will fail if MinIO is down, preventing
  traffic to unhealthy pods
- **Monitoring**: You can alert on storage health failures
- **Debugging**: Easy to see if storage issues are causing problems
- **Production**: Critical for ensuring the service is truly ready to accept
  requests

### Files Modified

- `internal/modules/media/service.go` - Added `HealthCheck()` to interface
- `internal/modules/media/service_minio.go` - Implemented MinIO health check
- `internal/modules/health/handler.go` - Added storage check to readiness probe

---

## Summary

✅ **Query Timeout**: All repository methods now use GORM callbacks to
automatically enforce query timeouts on all database operations  
✅ **Transactions**: Use `database.Transaction()` for atomic multi-operation
workflows  
✅ **MinIO Health**: Storage health check added to readiness probe for
production reliability

All three improvements are now production-ready and properly integrated!

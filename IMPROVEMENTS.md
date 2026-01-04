# Production-Grade Code Improvements

This document outlines critical improvements needed to make this codebase
production-ready at a senior level.

## 🔴 Critical Issues

### 1. **Context Cancellation in Database Transaction**

**File:** `internal/database/transaction.go` **Issue:** The transaction function
doesn't properly handle context cancellation. If the context is cancelled, the
transaction should be rolled back. **Fix:**

```go
func Transaction(ctx context.Context, db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create a channel to monitor context cancellation
	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
				log.Ctx(ctx).Error().Err(rollbackErr).Msg("failed to rollback transaction on context cancellation")
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
			log.Ctx(ctx).Error().Err(rollbackErr).Msg("failed to rollback transaction")
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to commit transaction")
		return err
	}

	return nil
}
```

### 2. **Context Leak in Query Timeout Callback**

**File:** `internal/database/postgres.go:185` **Issue:** The `cancel()` function
is intentionally ignored, which can cause context leaks. **Fix:** Store cancel
functions and call them when possible, or use a different approach:

```go
applyTimeout := func(db *gorm.DB) {
	if db.Statement != nil && db.Statement.Context != nil {
		if _, hasDeadline := db.Statement.Context.Deadline(); !hasDeadline {
			timeoutCtx, cancel := context.WithTimeout(db.Statement.Context, queryTimeout)
			// Store cancel in statement's context value or use a cleanup mechanism
			db.Statement.Context = timeoutCtx
			// Note: We can't call cancel here, but the timeout will be respected
			// The parent context cancellation will propagate
			_ = cancel // Suppress linter warning, but document the limitation
		}
	}
}
```

### 3. **Missing Index on Role Column**

**File:** `migrations/20251224093615_create_users_table.sql` **Issue:** The
`role` column is frequently used in queries but lacks an index. **Fix:** Add
index:

```sql
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
```

### 4. **SQL Injection Risk in Pagination**

**File:** `internal/pagination/pagination.go:100` **Issue:** Direct string
concatenation for ORDER BY is safe only because we validate against allowed
columns, but it's still fragile. **Fix:** Use parameterized queries or whitelist
validation (already done, but add comment):

```go
// SortBy is validated against SortableColumns, so this is safe
// However, consider using a map lookup for better performance
sortField := "created_at" // default
for _, allowed := range opts.SortableColumns {
	if stringx.ToSnakeCase(opts.SortBy) == stringx.ToSnakeCase(allowed) {
		sortField = stringx.ToSnakeCase(allowed)
		break
	}
}
// Use parameterized approach if possible, or ensure SortableColumns is always validated
db = db.Order(sortField + " " + strings.ToUpper(opts.Order))
```

### 5. **Race Condition in Auth Service RefreshToken**

**File:** `internal/modules/auth/service.go:65` **Issue:** TODO comment
indicates missing user status check. This is a security issue. **Fix:**
Implement user status validation:

```go
// 2. hit the db: ensure the user still exists and isn't blocked
user, err := s.userService.GetByID(ctx, claims.UserID)
if err != nil {
	return "", security.ErrInvalidToken
}

// Check if user is active/not blocked
if user.Status != "active" { // Assuming you add a Status field
	return "", security.ErrBlockedUser
}
```

## ⚠️ Performance Issues

### 6. **N+1 Query Potential in List Operations**

**File:** `internal/modules/users/repository.go:86` **Issue:** If you add
relationships later, this could cause N+1 queries. **Fix:** Use `Preload` for
relationships:

```go
func (r *repository) List(ctx context.Context, opts *pagination.QueryOptions) ([]*User, int64, error) {
	var users []*User
	var total int64

	// 1. Get total count using the SearchScope
	err := r.db.WithContext(ctx).Model(&User{}).
		Scopes(pagination.SearchScope(opts)).
		Count(&total).Error
	if err != nil {
		return nil, 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to count users",
			err,
		)
	}

	// 2. Fetch data using the Paginate Scope
	// If you add relationships, use Preload here
	err = r.db.WithContext(ctx).
		Scopes(pagination.Paginate(opts)).
		Find(&users).Error
	if err != nil {
		return nil, 0, apperror.Wrap(
			apperror.Internal,
			apperror.ErrDatabaseError.Code,
			"failed to build query",
			err,
		)
	}

	return users, total, nil
}
```

### 7. **Inefficient Search Query**

**File:** `internal/pagination/pagination.go:140` **Issue:** Using OR conditions
can be slow on large datasets. Consider full-text search for better performance.
**Fix:** For large datasets, consider PostgreSQL full-text search:

```go
// For large datasets, consider:
// db = db.Where("to_tsvector('english', " + columnName + ") @@ plainto_tsquery('english', ?)", opts.Search)
// Or use a search index
```

### 8. **Missing Composite Index for Common Queries**

**File:** `migrations/20251224093615_create_users_table.sql` **Issue:** Queries
filtering by role and deleted_at together need a composite index. **Fix:**

```sql
CREATE INDEX idx_users_role_deleted_at ON users(role, deleted_at) WHERE deleted_at IS NULL;
```

### 9. **Context Background in Media Upload**

**File:** `internal/modules/media/service_minio.go:78` **Issue:** Using
`context.Background()` loses request context, breaking tracing and cancellation.
**Fix:** Accept context as parameter:

```go
func (s *minioService) Upload(ctx context.Context, file *multipart.FileHeader, subDir fileCategory) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", apperror.Wrap(
			apperror.Internal,
			errFileOpen.Code,
			errFileOpen.Message,
			err,
		)
	}
	defer func() { _ = src.Close() }()

	// ... rest of the code
}
```

### 10. **No Connection Pooling Metrics**

**File:** `internal/database/postgres.go` **Issue:** No way to monitor
connection pool health in production. **Fix:** Add metrics endpoint or logging:

```go
// Add periodic health check
go func() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		stats := sqlDB.Stats()
		log.Info().
			Int("open_connections", stats.OpenConnections).
			Int("in_use", stats.InUse).
			Int("idle", stats.Idle).
			Int("wait_count", stats.WaitCount).
			Dur("wait_duration", stats.WaitDuration).
			Msg("database connection pool stats")
	}
}()
```

## 🔒 Security Issues

### 11. **Password Hash Cost Not Configurable**

**File:** `internal/modules/users/service.go:46` **Issue:** Using
`bcrypt.DefaultCost` may be too low for production. **Fix:** Make it
configurable:

```go
const (
	bcryptCost = 12 // Production should use 12-14
)

hash, err := bcrypt.GenerateFromPassword(
	[]byte(req.Password),
	bcryptCost,
)
```

### 12. **Missing Rate Limiting on Auth Endpoints**

**File:** `internal/routes/routes.go` (check if auth routes have rate limiting)
**Issue:** Auth endpoints should have stricter rate limiting to prevent brute
force. **Fix:** Ensure auth routes use `RateLimitRedisWithConfig` with stricter
limits.

### 13. **JWT Secret Should Be Validated**

**File:** `internal/security/jwt.go:31` **Issue:** No validation that secret is
strong enough. **Fix:** Add validation:

```go
func NewJWTHandler(secret string, accessTokenExpirySecond, refreshTokenExpirySecond int) (*JWTHandler, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT secret must be at least 32 characters")
	}
	// ... rest of code
}
```

### 14. **Missing Input Sanitization**

**File:** `internal/pagination/pagination.go:136` **Issue:** Search input should
be sanitized to prevent injection. **Fix:** Add input sanitization:

```go
// Sanitize search input
search := strings.TrimSpace(opts.Search)
// Remove potentially dangerous characters
search = strings.ReplaceAll(search, "%", "\\%")
search = strings.ReplaceAll(search, "_", "\\_")
```

### 15. **File Upload Size Check Before Processing**

**File:** `internal/modules/media/handler.go:97` **Issue:** ContentLength can be
spoofed. Always verify actual file size. **Fix:** The current implementation
checks `file.Size` which is correct, but add comment:

```go
// Note: ContentLength can be spoofed, so we also check file.Size after parsing
```

## 🏗️ Code Quality Issues

### 16. **Error Wrapping Inconsistency**

**File:** `internal/modules/users/service.go:51` **Issue:** Using
`apperror.ErrInternal` directly instead of wrapping with context. **Fix:**

```go
if err != nil {
	return nil, apperror.Wrap(
		apperror.Internal,
		apperror.ErrInternal.Code,
		"failed to hash password",
		err,
	)
}
```

### 17. **Missing Validation for ULID Format**

**File:** `internal/modules/users/repository.go:48` **Issue:** Should validate
ULID format before querying. **Fix:** Add validation in service layer or use a
custom validator.

### 18. **Hardcoded Timeout Values**

**File:** `internal/modules/media/service_minio.go:19` **Issue:** Timeout values
should be configurable. **Fix:** Move to config or constants with documentation.

### 19. **Missing Graceful Shutdown for Background Goroutines**

**File:** Check for any background goroutines **Issue:** If you add background
workers, ensure they respect shutdown signals. **Fix:** Use context cancellation
in all background operations.

### 20. **Inconsistent Error Messages**

**File:** Throughout codebase **Issue:** Some errors use "failed to X" while
others don't. **Fix:** Standardize error message format.

## 📊 Monitoring & Observability

### 21. **Missing Request Duration Metrics**

**File:** `internal/middlewares/request.go` **Issue:** No metrics for request
duration, which is critical for performance monitoring. **Fix:** Add Prometheus
metrics or structured logging with duration:

```go
// Add request duration to logs
duration := time.Since(start)
log.Ctx(c.Request.Context()).Info().
	Dur("duration_ms", duration).
	Int("status", c.Writer.Status()).
	Msg("request completed")
```

### 22. **Missing Database Query Metrics**

**File:** `internal/database/postgres.go:46` **Issue:** While slow queries are
logged, there's no aggregation of query performance. **Fix:** Consider adding
metrics collection for:

- Query count by type
- Average query duration
- Error rate

### 23. **Health Check Doesn't Verify All Dependencies**

**File:** `internal/modules/health/handler.go:48` **Issue:** Health checks
should verify all critical dependencies. **Fix:** Already good, but ensure all
are checked.

## 🚀 Best Practices

### 24. **Use Database Transactions for Multi-Step Operations**

**File:** `internal/modules/users/service.go:35` **Issue:** User creation
involves multiple steps but doesn't use transactions. **Fix:** Wrap in
transaction if you add more steps:

```go
// If you add more operations, use:
return database.Transaction(ctx, s.repo.db, func(tx *gorm.DB) error {
	// ... operations
})
```

### 25. **Add Request ID to All Logs**

**File:** Throughout codebase **Issue:** Ensure all logs include request ID for
tracing. **Fix:** Already implemented via middleware, but verify all log calls
use `log.Ctx(ctx)`.

### 26. **Implement Circuit Breaker for External Services**

**File:** `internal/modules/media/service_minio.go` **Issue:** No circuit
breaker for MinIO calls. **Fix:** Consider adding circuit breaker pattern for
external service calls.

### 27. **Add Database Query Result Caching**

**File:** `internal/modules/users/repository.go` **Issue:** Frequently accessed
data (like user by ID) could be cached. **Fix:** Add Redis caching layer for
read-heavy operations:

```go
// Check cache first
cached, err := cache.Get(ctx, "user:"+id)
if err == nil {
	return cached, nil
}
// ... fetch from DB and cache
```

### 28. **Implement Retry Logic for Transient Failures**

**File:** External service calls **Issue:** No retry logic for transient
failures. **Fix:** Add retry with exponential backoff for:

- Database connection failures (already done)
- MinIO operations
- Redis operations

### 29. **Add Request Validation Middleware**

**File:** `internal/httpx/bind.go` **Issue:** Validation happens but could be
more centralized. **Fix:** Already good, but ensure all endpoints validate
input.

### 30. **Implement Proper API Versioning**

**File:** `internal/routes/routes.go` **Issue:** No API versioning strategy
visible. **Fix:** Add version prefix: `/api/v1/...`

## 📝 Documentation

### 31. **Add API Documentation for Error Responses**

**File:** Swagger annotations **Issue:** Some endpoints don't document all
possible error responses. **Fix:** Ensure all endpoints document:

- 400 Bad Request
- 401 Unauthorized
- 403 Forbidden
- 404 Not Found
- 500 Internal Server Error

### 32. **Add Code Comments for Complex Logic**

**File:** `internal/pagination/pagination.go:115` **Issue:** Complex search
logic needs better documentation. **Fix:** Add detailed comments explaining the
search algorithm.

## 🔧 Configuration

### 33. **Make All Magic Numbers Configurable**

**File:** Throughout codebase **Issue:** Hardcoded values like timeouts, limits
should be in config. **Fix:** Move to config file with sensible defaults.

### 34. **Add Environment-Specific Configurations**

**File:** `internal/config/config.go` **Issue:** Ensure different configs for
dev/staging/prod. **Fix:** Use environment variables with validation.

## 🧪 Testing

### 35. **Add Integration Tests**

**Issue:** No visible integration tests. **Fix:** Add tests for:

- Database operations
- Authentication flow
- File upload
- Pagination

### 36. **Add Load Testing**

**Issue:** No performance benchmarks. **Fix:** Add load tests for critical
endpoints.

## Summary Priority

**Immediate (Before Production):**

1. Context cancellation in transactions (#1)
2. User status check in refresh token (#5)
3. Context parameter in media upload (#9)
4. Password hash cost configuration (#11)
5. JWT secret validation (#13)

**High Priority (First Sprint):** 6. Database indexes (#3, #8) 7. Request
duration metrics (#21) 8. Error wrapping consistency (#16) 9. Input sanitization
(#14) 10. API versioning (#30)

**Medium Priority:**

- Performance optimizations (#6, #7)
- Monitoring improvements (#22)
- Caching layer (#27)

**Low Priority (Technical Debt):**

- Documentation improvements
- Code style consistency
- Additional test coverage

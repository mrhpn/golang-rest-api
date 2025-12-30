# Production Deployment Guide

This document outlines the production-grade features and best practices
implemented in this Go REST API.

## 🚀 Production-Grade Features

### 1. **Rate Limiting**

- In-memory rate limiter using token bucket algorithm
- Configurable rate and time window
- Per-IP address tracking
- Automatic cleanup of old entries
- Default: 100 requests per minute per IP

**Configuration:**

```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RATE=100
RATE_LIMIT_WINDOW_SECOND=60
```

### 2. **Security Headers**

- X-Frame-Options: DENY (prevents clickjacking)
- X-Content-Type-Options: nosniff (prevents MIME sniffing)
- X-XSS-Protection: 1; mode=block
- Referrer-Policy: strict-origin-when-cross-origin
- Content-Security-Policy: configured for secure content loading
- Permissions-Policy: restricts browser features

### 3. **Request Timeout**

- Prevents long-running requests from consuming resources
- Configurable per-request timeout
- Default: 30 seconds

**Configuration:**

```bash
REQUEST_TIMEOUT_SECOND=30
```

### 4. **Enhanced Health Checks**

- **Liveness Probe** (`/health/live`): Checks if service is alive
- **Readiness Probe** (`/health/ready`): Checks if service is ready (includes DB
  health)
- **Health Check** (`/health`): Basic health status
- Database connectivity checks
- Connection pool statistics

### 5. **Database Connection Management**

#### Connection Pool Configuration

```bash
DB_MAX_OPEN_CONNS=25        # Maximum open connections
DB_MAX_IDLE_CONNS=10        # Maximum idle connections
DB_CONN_MAX_LIFETIME_MINUTE=60  # Connection max lifetime
DB_CONN_MAX_IDLE_TIME_MINUTE=30 # Connection max idle time
DB_QUERY_TIMEOUT_SECOND=30      # Query timeout
```

#### Retry Logic

- Exponential backoff retry for database connections
- Configurable retry attempts and delay
- Default: 3 attempts with 2-second initial delay

**Configuration:**

```bash
DB_RETRY_ATTEMPTS=3
DB_RETRY_DELAY_SECOND=2
```

#### Query Timeout

- All database queries respect configured timeout
- Prevents runaway queries
- Context-aware timeout handling

### 6. **Transaction Support**

- Helper function for database transactions
- Automatic rollback on error
- Panic recovery with rollback
- Located in `internal/database/transaction.go`

**Usage:**

```go
err := database.Transaction(ctx, db, func(tx *gorm.DB) error {
    // Your transaction logic here
    return nil
})
```

### 7. **Graceful Shutdown**

- Handles SIGTERM and SIGINT signals
- 15-second timeout for graceful shutdown
- Proper database connection cleanup
- Prevents connection leaks

### 8. **API Versioning**

- Current version: `/api/v1`
- Legacy routes: `/api` (for backward compatibility)
- Easy to add new versions in the future

### 9. **Enhanced Logging**

- Structured logging with Zerolog
- Request ID tracking for distributed tracing
- Context-aware logging
- Log rotation with Lumberjack
- Separate log levels for development/production

**Configuration:**

```bash
LOG_LEVEL=INFO
LOG_PATH=./logs
LOG_MAX_SIZE_MB=100
LOG_MAX_BACKUP_COUNT=8
LOG_MAX_DAY=30
LOG_COMPRESS=true
```

### 10. **Error Handling**

- Standardized error responses
- Custom error types with HTTP status mapping
- Field-level validation errors
- Proper error propagation

### 11. **CORS Configuration**

- Configurable allowed origins
- Proper credentials handling
- Preflight request support

**Configuration:**

```bash
ALLOWED_ORIGINS=https://example.com,https://app.example.com
# Or use "*" for development (not recommended for production)
```

### 12. **Request Body Size Limits**

- Configurable maximum request body size
- Default: 50MB
- Prevents memory exhaustion

**Configuration:**

```bash
MAX_REQUEST_BODY_SIZE_MB=50
```

## 📋 Environment Variables

### Required Variables

```bash
DATABASE_URL=postgres://user:password@host:port/dbname?sslmode=require
JWT_SECRET=your-secret-key-min-32-chars
STORAGE_HOST=minio.example.com
```

### Optional Variables (with defaults)

```bash
APP_ENV=production
APP_PORT=8080
REQUEST_TIMEOUT_SECOND=30

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RATE=100
RATE_LIMIT_WINDOW_SECOND=60

# Database
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME_MINUTE=60
DB_CONN_MAX_IDLE_TIME_MINUTE=30
DB_QUERY_TIMEOUT_SECOND=30
DB_RETRY_ATTEMPTS=3
DB_RETRY_DELAY_SECOND=2

# Logging
LOG_LEVEL=INFO
LOG_PATH=./logs
LOG_MAX_SIZE_MB=100
LOG_MAX_BACKUP_COUNT=8
LOG_MAX_DAY=30
LOG_COMPRESS=true

# CORS
ALLOWED_ORIGINS=*

# Storage
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
STORAGE_BUCKET_NAME=app_assets
STORAGE_USE_SSL=false

# JWT
ACCESS_TOKEN_EXPIRATION_TIME_SECOND=3600
REFRESH_TOKEN_EXPIRATION_TIME_SECOND=86400
```

## 🐳 Docker Deployment

### Build

```bash
docker build -t go-rest-api .
```

### Run

```bash
docker run -d \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://... \
  -e JWT_SECRET=... \
  -e STORAGE_HOST=... \
  --name go-rest-api \
  go-rest-api
```

## ☸️ Kubernetes Deployment

### Health Checks

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### Resource Limits

```yaml
resources:
  requests:
    memory: '256Mi'
    cpu: '250m'
  limits:
    memory: '512Mi'
    cpu: '500m'
```

## 🔒 Security Best Practices

1. **Never commit secrets** - Use environment variables or secret management
2. **Use HTTPS in production** - Enable HSTS header
3. **Rotate JWT secrets regularly**
4. **Set appropriate CORS origins** - Never use "\*" in production
5. **Enable rate limiting** - Protect against DDoS
6. **Monitor logs** - Set up log aggregation
7. **Database credentials** - Use strong passwords and connection pooling
8. **Regular updates** - Keep dependencies updated

## 📊 Monitoring Recommendations

1. **Application Metrics**

   - Request rate
   - Error rate
   - Response time (p50, p95, p99)
   - Database connection pool usage

2. **Infrastructure Metrics**

   - CPU usage
   - Memory usage
   - Disk I/O
   - Network I/O

3. **Health Checks**
   - Monitor `/health/ready` endpoint
   - Alert on 503 responses
   - Track database connectivity

## 🚨 Troubleshooting

### High Database Connection Usage

- Check `DB_MAX_OPEN_CONNS` setting
- Monitor connection pool stats via `/health/ready`
- Review slow query logs

### Rate Limit Issues

- Adjust `RATE_LIMIT_RATE` and `RATE_LIMIT_WINDOW_SECOND`
- Consider per-endpoint rate limits for specific routes

### Timeout Issues

- Increase `REQUEST_TIMEOUT_SECOND` if needed
- Check `DB_QUERY_TIMEOUT_SECOND` for slow queries
- Review slow query logs

### Memory Issues

- Reduce `MAX_REQUEST_BODY_SIZE_MB`
- Check for memory leaks in handlers
- Monitor goroutine count

## 📝 Additional Notes

- All database operations use context for cancellation
- Prepared statements are cached for performance
- Global updates/deletes are prevented (safety feature)
- Request IDs are included in all logs for tracing
- Structured logging enables easy log parsing

package constants

import "time"

// Time constants
const (
	APIDateTimeLayout = time.RFC3339
)

// Common constants
const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

// Pagination constants
const (
	PaginationDefaultLimit = 10 // Default items per page
	PaginationMaxLimit     = 50 // Maximum items per page
	PaginationDefaultPage  = 1  // Default page number
)

// API constants
const (
	APIPrefix         = "api"
	CurrentAPIVersion = "v1"
	APIAuthPrefix     = "auth"

	APIVersionPrefix = "/" + APIPrefix + "/" + CurrentAPIVersion
	APIAuthPath      = APIVersionPrefix + "/" + APIAuthPrefix
)

// Security constants
const (
	DefaultBcryptCost = 14

	RateLimit     = "100-M" // 100 requests per minute
	RateLimitAuth = "7-M"   // 7 requests per minute for auth routes

	AccessTokenExpirationSecond  = 3600
	RefreshTokenExpirationSecond = 86400

	JWTSecretMinLength = 32

	RateLimitKeyPrefix = "ratelimit:"
)

// Media constants
const (
	MaxProfileImageWidth   = 400
	MaxProfileImageHeight  = 400
	MaxProfileImageQuality = 75

	MaxThumbnailImageWidth   = 800
	MaxThumbnailImageHeight  = 600
	MaxThumbnailImageQuality = 80

	MaxImageSize    = 5 * MB
	MaxVideoSize    = 50 * MB
	MaxDocumentSize = 10 * MB
)

// Server constants
const (
	EnvDev  = "development"
	EnvProd = "production"
	EnvTest = "testing"

	RequestTimeoutSecond     = 30
	RequestMaxBodySizeMB     = 50 * MB
	ServerReadTimeoutSecond  = 20
	ServerWriteTimeoutSecond = 30
	ServerIdleTimeoutSecond  = 120

	MaxMultipartMemoryMB = 8
	MaxMultipartMemory   = MaxMultipartMemoryMB * MB
)

// DB constants
const (
	DBMaxQueryTimeoutSecond        = 30
	DBMaxRetryAttempts             = 3
	DBRetryDelaySecond             = 2
	DBMaxOpenConns                 = 25
	DBMaxIdleConns                 = 10
	DBMaxLifetimeMinute            = 60
	DBConnMaxIdleTimeMinute        = 30
	DBPoolMetricsEnabled           = true
	DBPoolMetricsLogIntervalSecond = 30
)

// Redis constants
const (
	RedisPoolSize                 = 10
	RedisMinIdleConns             = 5
	RedisDialTimeoutSecond        = 5 * time.Second
	RedisReadTimeoutSecond        = 3 * time.Second
	RedisWriteTimeoutSecond       = 3 * time.Second
	RedisConnMaxIdleTimeMinute    = 5 * time.Minute
	RedisConnMaxLifetimeMinute    = 30 * time.Minute
	RedisHealthCheckTimeoutSecond = 5 * time.Second
)

// Logger constants
const (
	LogMaxDay    = 30
	LogMaxBackup = 8
	LogMaxSizeMB = 100
)

package constants

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

const (
	APIPrefix = "api"
	CurrentAPIVersion = "v1"
	APIAuthPrefix = "auth"

	APIVersionPrefix = "/" + APIPrefix + "/" + CurrentAPIVersion
	APIAuthPath      = APIVersionPrefix + "/" + APIAuthPrefix
)

const (
	DefaultBcryptCost = 14
)

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

const (
	RequestTimeoutSecond     = 30
	RequestMaxBodySizeMB     = 50 * MB
	ServerReadTimeoutSecond  = 20
	ServerWriteTimeoutSecond = 30
	ServerIdleTimeoutSecond  = 120

	MaxMultipartMemoryMB = 8
	MaxMultipartMemory   = MaxMultipartMemoryMB * MB

	DBMaxQueryTimeoutSecond = 30
	DBMaxRetryAttempts      = 3
	DBRetryDelaySecond      = 2
	DBMaxOpenConns          = 25
	DBMaxIdleConns          = 10
	DBMaxLifetimeMinute     = 60
	DBConnMaxIdleTimeMinute = 30

	LogMaxDay    = 30
	LogMaxBackup = 8
	LogMaxSizeMB = 100

	RateLimit       = 100
	RateLimitAuth   = 7
	RateLimitWindow = 60

	AccessTokenExpirationSecond  = 3600
	RefreshTokenExpirationSecond = 86400
)

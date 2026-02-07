// Package config provides application configuration loading and access.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mrhpn/go-rest-api/internal/constants"
)

// Config represents the full application configuration loaded at startup.
type Config struct {
	AppEnv    string
	Port      string
	DBURL     string
	HTTP      HTTPConfig
	RateLimit RateLimitConfig
	DB        DBConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Log       LogConfig
	Storage   StorageConfig
}

// HTTPConfig represents the http-related config
type HTTPConfig struct {
	AllowedOrigins       []string
	MaxRequestBodySize   int64
	RequestTimeoutSecond int
}

// RateLimitConfig represents rate limit related config
// Rate and AuthRate use ulule/limiter format: "100-M" (100 per minute), "50-H" (50 per hour), "10-S" (10 per second)
type RateLimitConfig struct {
	Enabled  bool
	Rate     string // rate limit in ulule/limiter format (e.g., "100-M" for 100 per minute)
	AuthRate string // auth route rate limit in ulule/limiter format (e.g., "7-M" for 7 per minute)
}

// DBConfig represents database related config
type DBConfig struct {
	MaxOpenConns                   int  // maximum open connections
	MaxIdleConns                   int  // maximum idle connections
	ConnMaxLifetimeMinute          int  // connection max lifetime in minutes
	ConnMaxIdleTimeMinute          int  // connection max idle time in minutes
	QueryTimeoutSecond             int  // query timeout in seconds
	RetryAttempts                  int  // number of retry attempts
	RetryDelaySecond               int  // retry delay in seconds
	DBPoolMetricsEnabled           bool // db pool metrics log enabled or not
	DBPoolMetricsLogIntervalSecond int  // db pool metrics log interval second
}

// RedisConfig represents redis related config
type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	DB       int // Redis database number (0-15)
}

// JWTConfig represents app's auth (jwt) related config
type JWTConfig struct {
	Secret                       string
	AccessTokenExpirationSecond  int // in seconds
	RefreshTokenExpirationSecond int // in seconds
}

// LogConfig represents app's logger related config
type LogConfig struct {
	Path           string
	Level          string
	MaxSizeMB      int
	MaxBackupCount int
	MaxAgeDay      int
	Compress       bool
}

// StorageConfig represents app's storage (minio/local) related config
type StorageConfig struct {
	Provider   string
	Host       string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
	LocalPath  string
}

// Load loads the application configuration from environment variables.
// It returns an error if any required configuration is missing.
func Load() (*Config, error) {
	originsRaw := getEnv("ALLOWED_ORIGINS", "*")
	var allowedOrigins []string
	if originsRaw == "*" {
		allowedOrigins = []string{"*"}
	} else {
		allowedOrigins = strings.Split(originsRaw, ",")
	}

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", constants.EnvDev),
		Port:   getEnv("APP_PORT", "8080"),
		DBURL:  getEnv("DATABASE_URL", ""),

		HTTP: HTTPConfig{
			AllowedOrigins:       allowedOrigins,
			MaxRequestBodySize:   constants.RequestMaxBodySizeMB,
			RequestTimeoutSecond: constants.RequestTimeoutSecond,
		},

		RateLimit: RateLimitConfig{
			Enabled:  getEnvAsBool("RATE_LIMIT_ENABLED", true),
			Rate:     getEnv("RATE_LIMIT_RATE", constants.RateLimit),
			AuthRate: getEnv("RATE_LIMIT_AUTH_RATE", constants.RateLimitAuth),
		},

		DB: DBConfig{
			MaxOpenConns:                   getEnvAsInt("DB_MAX_OPEN_CONNS", constants.DBMaxOpenConns),
			MaxIdleConns:                   getEnvAsInt("DB_MAX_IDLE_CONNS", constants.DBMaxIdleConns),
			ConnMaxLifetimeMinute:          getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTE", constants.DBMaxLifetimeMinute),
			ConnMaxIdleTimeMinute:          getEnvAsInt("DB_CONN_MAX_IDLE_TIME_MINUTE", constants.DBConnMaxIdleTimeMinute),
			QueryTimeoutSecond:             getEnvAsInt("DB_QUERY_TIMEOUT_SECOND", constants.DBMaxQueryTimeoutSecond),
			RetryAttempts:                  getEnvAsInt("DB_RETRY_ATTEMPTS", constants.DBMaxRetryAttempts),
			RetryDelaySecond:               getEnvAsInt("DB_RETRY_DELAY_SECOND", constants.DBRetryDelaySecond),
			DBPoolMetricsEnabled:           getEnvAsBool("DB_POOL_METRICS_ENABLED", constants.DBPoolMetricsEnabled),
			DBPoolMetricsLogIntervalSecond: getEnvAsInt("DB_POOL_METRICS_INTERVAL_SECOND", constants.DBPoolMetricsLogIntervalSecond),
		},

		Redis: RedisConfig{
			// Default to enabled in production (can be overridden via env)
			// In Docker, this will be enabled via environment variable
			Enabled:  getEnvAsBool("REDIS_ENABLED", false),
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},

		JWT: JWTConfig{
			Secret:                       getEnv("JWT_SECRET", ""),
			AccessTokenExpirationSecond:  getEnvAsInt("ACCESS_TOKEN_EXPIRATION_TIME_SECOND", constants.AccessTokenExpirationSecond),
			RefreshTokenExpirationSecond: getEnvAsInt("REFRESH_TOKEN_EXPIRATION_TIME_SECOND", constants.RefreshTokenExpirationSecond),
		},

		Log: LogConfig{
			Path:           getEnv("LOG_PATH", "./logs"),
			Level:          getEnv("LOG_LEVEL", "INFO"),
			MaxSizeMB:      getEnvAsInt("LOG_MAX_SIZE_MB", constants.LogMaxSizeMB),
			MaxBackupCount: getEnvAsInt("LOG_MAX_BACKUP_COUNT", constants.LogMaxBackup),
			MaxAgeDay:      getEnvAsInt("LOG_MAX_DAY", constants.LogMaxDay),
			Compress:       getEnvAsBool("LOG_COMPRESS", true),
		},

		Storage: StorageConfig{
			Provider:   getEnv("STORAGE_PROVIDER", "minio"),
			Host:       getEnv("STORAGE_HOST", ""),
			AccessKey:  getEnv("STORAGE_ACCESS_KEY", "minioadmin"),
			SecretKey:  getEnv("STORAGE_SECRET_KEY", "minioadmin"),
			BucketName: getEnv("STORAGE_BUCKET_NAME", "app_assets"),
			UseSSL:     getEnvAsBool("STORAGE_USE_SSL", false),
			LocalPath:  getEnv("STORAGE_LOCAL_PATH", "./uploads"),
		},
	}

	if cfg.DBURL == "" {
		return nil, errors.New("env: DATABASE_URL is missing")
	}
	if cfg.JWT.Secret == "" || len(cfg.JWT.Secret) < constants.JWTSecretMinLength {
		return nil, fmt.Errorf("env: JWT_SECRET is missing or less than %d characters", constants.JWTSecretMinLength)
	}
	switch strings.ToLower(cfg.Storage.Provider) {
	case "minio":
		if cfg.Storage.Host == "" {
			return nil, errors.New("env: STORAGE_HOST is missing")
		}
	case "local":
		if cfg.Storage.LocalPath == "" {
			return nil, errors.New("env: STORAGE_LOCAL_PATH is missing")
		}
	default:
		return nil, errors.New("env: STORAGE_PROVIDER is invalid (should be minio | local)")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return fallback
}

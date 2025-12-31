// Package config provides application configuration loading and access.
package config

import (
	"os"
	"strconv"
	"strings"
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
	Storage   MinioConfig
}

// HTTPConfig represents the http-related config
type HTTPConfig struct {
	AllowedOrigins       []string
	MaxRequestBodySize   int64
	RequestTimeoutSecond int
}

// RateLimitConfig represents rate limit related config
type RateLimitConfig struct {
	Enabled  bool
	Rate     int // requests per window
	AuthRate int
	Window   int // window in seconds
}

// DBConfig represents database related config
type DBConfig struct {
	MaxOpenConns          int // maximum open connections
	MaxIdleConns          int // maximum idle connections
	ConnMaxLifetimeMinute int // connection max lifetime in minutes
	ConnMaxIdleTimeMinute int // connection max idle time in minutes
	QueryTimeoutSecond    int // query timeout in seconds
	RetryAttempts         int // number of retry attempts
	RetryDelaySecond      int // retry delay in seconds
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

// MinioConfig represents app's storage (minio) related config
type MinioConfig struct {
	Host       string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

// MustLoad loads the application configuration from environment variables.
// It panics if any required configuration is missing.
func MustLoad() *Config {
	originsRaw := getEnv("ALLOWED_ORIGINS", "*")
	var allowedOrigins []string
	if originsRaw == "*" {
		allowedOrigins = []string{"*"}
	} else {
		allowedOrigins = strings.Split(originsRaw, ",")
	}

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("APP_PORT", "8080"),
		DBURL:  getEnv("DATABASE_URL", ""),

		HTTP: HTTPConfig{
			AllowedOrigins:       allowedOrigins,
			MaxRequestBodySize:   int64(getEnvAsInt("MAX_REQUEST_BODY_SIZE_MB", 50)) * 1024 * 1024,
			RequestTimeoutSecond: getEnvAsInt("REQUEST_TIMEOUT_SECOND", 30),
		},

		RateLimit: RateLimitConfig{
			Enabled:  getEnvAsBool("RATE_LIMIT_ENABLED", true),
			Rate:     getEnvAsInt("RATE_LIMIT_RATE", 100),
			AuthRate: getEnvAsInt("RATE_LIMIT_AUTH_RATE", 5),
			Window:   getEnvAsInt("RATE_LIMIT_WINDOW_SECOND", 60),
		},

		DB: DBConfig{
			MaxOpenConns:          getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:          getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetimeMinute: getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTE", 60),
			ConnMaxIdleTimeMinute: getEnvAsInt("DB_CONN_MAX_IDLE_TIME_MINUTE", 30),
			QueryTimeoutSecond:    getEnvAsInt("DB_QUERY_TIMEOUT_SECOND", 30),
			RetryAttempts:         getEnvAsInt("DB_RETRY_ATTEMPTS", 3),
			RetryDelaySecond:      getEnvAsInt("DB_RETRY_DELAY_SECOND", 2),
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
			AccessTokenExpirationSecond:  getEnvAsInt("ACCESS_TOKEN_EXPIRATION_TIME_SECOND", 3600),
			RefreshTokenExpirationSecond: getEnvAsInt("REFRESH_TOKEN_EXPIRATION_TIME_SECOND", 86400),
		},

		Log: LogConfig{
			Path:           getEnv("LOG_PATH", "./logs"),
			Level:          getEnv("LOG_LEVEL", "INFO"),
			MaxSizeMB:      getEnvAsInt("LOG_MAX_SIZE_MB", 100),
			MaxBackupCount: getEnvAsInt("LOG_MAX_BACKUP_COUNT", 8),
			MaxAgeDay:      getEnvAsInt("LOG_MAX_DAY", 30),
			Compress:       getEnvAsBool("LOG_COMPRESS", true),
		},

		Storage: MinioConfig{
			Host:       getEnv("STORAGE_HOST", ""),
			AccessKey:  getEnv("STORAGE_ACCESS_KEY", "minioadmin"),
			SecretKey:  getEnv("STORAGE_SECRET_KEY", "minioadmin"),
			BucketName: getEnv("STORAGE_BUCKET_NAME", "app_assets"),
			UseSSL:     getEnvAsBool("STORAGE_USE_SSL", false),
		},
	}

	if cfg.DBURL == "" {
		panic("env: DATABASE_URL is missing")
	}
	if cfg.JWT.Secret == "" {
		panic("env: JWT_SECRET is missing")
	}
	if cfg.Storage.Host == "" {
		panic("env: STORAGE_HOST is missing")
	}

	return cfg
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

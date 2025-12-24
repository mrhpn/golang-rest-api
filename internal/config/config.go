package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv         string
	Port           string
	DBUrl          string
	AllowedOrigins []string
	JWT            JWTConfig
	Log            LogConfig
}

type JWTConfig struct {
	Secret                       string
	AccessTokenExpirationSecond  int // in seconds
	RefreshTokenExpirationSecond int // in seconds
}

type LogConfig struct {
	Path           string
	Level          string
	MaxSizeMB      int
	MaxBackupCount int
	MaxAgeDay      int
	Compress       bool
}

func MustLoad() *Config {
	originsRaw := getEnv("ALLOWED_ORIGINS", "*")
	var allowedOrigins []string
	if originsRaw == "*" {
		allowedOrigins = []string{"*"}
	} else {
		allowedOrigins = strings.Split(originsRaw, ",")
	}

	cfg := &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		Port:           getEnv("APP_PORT", "8080"),
		DBUrl:          getEnv("DATABASE_URL", ""),
		AllowedOrigins: allowedOrigins,

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
	}

	if cfg.DBUrl == "" {
		panic("env: DATABASE_URL is missing")
	}
	if cfg.JWT.Secret == "" {
		panic("env: JWT_SECRET is missing")
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

package config

import "os"

type Config struct {
	AppEnv    string
	Port      string
	DBUrl     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		AppEnv:    getEnv("APP_ENV", "development"),
		Port:      getEnv("APP_PORT", "8080"),
		DBUrl:     getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
		JWTSecret: getEnv("JW_SECRET", "2026J@NU@RY"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

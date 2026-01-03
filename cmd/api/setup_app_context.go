package main

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/modules/media"
	"github.com/mrhpn/go-rest-api/internal/security"
)

func setupAppContext(cfg *config.Config, db *gorm.DB, redis *redis.Client, logger zerolog.Logger, media media.Service) *app.Context {
	securityHandler := security.NewJWTHandler(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpirationSecond,
		cfg.JWT.RefreshTokenExpirationSecond,
	)

	return &app.Context{
		DB:              db,
		Redis:           redis,
		Cfg:             cfg,
		Logger:          logger,
		SecurityHandler: securityHandler,
		MediaService:    media,
	}
}

package main

import (
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/security"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func setupAppContext(cfg *config.Config, db *gorm.DB, logger zerolog.Logger) *app.AppContext {
	securityHandler := security.NewJWTHandler(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpirationSecond,
		cfg.JWT.RefreshTokenExpirationSecond,
	)

	return &app.AppContext{
		DB:              db,
		Cfg:             cfg,
		Logger:          logger,
		SecurityHandler: securityHandler,
	}
}

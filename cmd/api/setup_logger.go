package main

import (
	"github.com/rs/zerolog"

	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/config"
)

func setupLogger(cfg *config.Config) zerolog.Logger {
	return app.SetupLogger(&cfg.Log, cfg.AppEnv)
}

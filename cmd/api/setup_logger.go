package main

import (
	"github.com/mrhpn/go-rest-api/internal/app"
	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogger(cfg *config.Config) zerolog.Logger {
	logger := app.SetupLogger(&cfg.Log, cfg.AppEnv)
	log.Logger = logger
	return logger
}

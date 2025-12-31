package main

import (
	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func setupDatabase(cfg *config.Config) *gorm.DB {
	db, err := database.Connect(cfg.DBURL, &cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Database connection failed")
	}
	return db
}

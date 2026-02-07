package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/database"
)

const (
	timeoutSecond = 5 * time.Second
)

func setupDatabase(ctx context.Context, cfg *config.Config) (*gorm.DB, func(), error) {
	db, err := database.Connect(ctx, cfg.DBURL, &cfg.DB)
	if err != nil {
		return nil, nil, fmt.Errorf("database connection failed: %w", err)
	}

	cleanup := func() {
		sqlDB, dbErr := db.DB()
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("❌ Failed to get database connection for cleanup")
		} else if sqlDB != nil {
			closeCtx, cancel := context.WithTimeout(context.Background(), timeoutSecond)
			defer cancel()

			closed := make(chan error, 1)
			go func() {
				closed <- sqlDB.Close()
			}()

			select {
			case dbErr = <-closed:
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("❌ Failed to close database connection")
				} else {
					log.Info().Msg("✓ Database connection closed")
				}
			case <-closeCtx.Done():
				log.Warn().Msg("⚠ Database connection close timeout")
			}
		}
	}

	return db, cleanup, nil
}

package main

import (
	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/modules/media"
	"github.com/rs/zerolog/log"
)

func setupMedia(cfg *config.Config) media.Service {
	svc, err := media.NewMinioService(
		cfg.Storage.Host,
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.BucketName,
		cfg.Storage.UseSSL,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("❌ failed to initialize MinIO storage service")
	}

	return svc
}

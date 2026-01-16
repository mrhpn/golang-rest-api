package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/modules/media"
)

func setupMedia(cfg *config.Config) (media.Service, func(), error) {
	svc, err := media.NewMinioService(
		cfg.Storage.Host,
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.BucketName,
		cfg.Storage.UseSSL,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup media service: %w", err)
	}

	cleanup := func() {
		log.Info().
			Str("service", "minio").
			Str("host", cfg.Storage.Host).
			Str("bucket", cfg.Storage.BucketName).
			Bool("ssl", cfg.Storage.UseSSL).
			Msg("✓ MinIO client cleanup completed")
	}

	return svc, cleanup, nil
}

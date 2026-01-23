package media

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

const bucketCreateTimeout = 10 * time.Second

// MinioConfig holds the raw connection details
type MinioConfig struct {
	Host      string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func SetupMinio(cfg MinioConfig) (*minio.Client, error) {
	// 1. Initialize client
	client, err := minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	log.Info().Msg("✅ MinIO — connected successfully")

	// 2. Prepare bucket
	ctx, cancel := context.WithTimeout(context.Background(), bucketCreateTimeout)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		if err = makePublicBucket(ctx, client, cfg.Bucket); err != nil {
			return nil, err
		}
	}

	log.Info().Msg("✅ MinIO — Bucket checked successfully")

	return client, nil
}

// makePublicBucket is a private helper for MinIO-specific policy configurations
func makePublicBucket(ctx context.Context, client *minio.Client, bucketName string) error {
	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
				"Action": ["s3:GetObject"],
				"Effect":"Allow",
				"Principal":"*",
				"Resource":["arn:aws:s3:::%s/*"]
			}]
		}`, bucketName)

	if err = client.SetBucketPolicy(ctx, bucketName, policy); err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}
	return nil
}

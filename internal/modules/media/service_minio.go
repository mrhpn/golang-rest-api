package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioService struct {
	client     *minio.Client
	bucketName string
}

func NewMinioService(host, accessKey, secretKey, bucketName string, useSSL bool) (Service, error) {
	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	// 1. auto-create bucket if it doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if minio bucket exists: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		// 2. set public policy so images can be viewed in browser via direct lin[k
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
					"Action": ["s3:GetObject"],
					"Effect":"Allow",
					"Principal":"*",
					"Resource":["arn:aws:s3:::%s/*"]
				}]
			}`, bucketName)

		err = client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	return &minioService{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// Upload streams the file to MinIO and returns the path
func (s *minioService) Upload(file *multipart.FileHeader, subDir string) (string, error) {
	ctx := context.Background()

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("%s/%s%s", subDir, uuid.New().String(), ext)

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to minio: %w", err)
	}

	return fmt.Sprintf("/%s", objectName), nil
}

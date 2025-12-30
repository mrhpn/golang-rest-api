package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	appErr "github.com/mrhpn/go-rest-api/internal/errors"
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
func (s *minioService) Upload(file *multipart.FileHeader, subDir FileCategory) (string, error) {
	ctx := context.Background()

	src, err := file.Open()
	if err != nil {
		return "", appErr.Wrap(
			appErr.Internal,
			ErrFileOpen.Code,
			ErrFileOpen.Message,
			err,
		)
	}
	defer src.Close()

	// 1. initialize vars with raw upload data
	var reader io.Reader = src
	size := file.Size
	contentType := file.Header.Get("Content-Type")
	ext := filepath.Ext(file.Filename)

	// 2. define processing rules based on directory
	var opts *ImageOptions
	switch subDir {
	case FileCategoryProfiles:
		opts = &ImageOptions{MaxWidth: 400, MaxHeight: 400, Quality: 75}
	case FileCategoryThumbnails:
		opts = &ImageOptions{MaxWidth: 800, MaxHeight: 600, Quality: 80}
	}

	// 3. apply processing if options were found
	if opts != nil {
		processed, newSize, err := ProcessImage(src, *opts)
		if err == nil {
			reader = processed
			size = newSize
			contentType = "image/jpeg"
			ext = ".jpg"
		}
		// Note: If image processing fails, we continue with the original file
	}

	// 4. construct the obj name
	objectName := fmt.Sprintf("%s/%s%s", subDir, uuid.New().String(), ext)

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", appErr.Wrap(
			appErr.Internal,
			ErrUploadToStorage.Code,
			ErrUploadToStorage.Message,
			err,
		)
	}

	return fmt.Sprintf("/%s", objectName), nil
}

// HealthCheck verifies that MinIO is accessible and the bucket exists
func (s *minioService) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if bucket exists and is accessible
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return appErr.Wrap(
			appErr.Internal,
			ErrStorageHealthCheck.Code,
			ErrStorageHealthCheck.Message,
			err,
		)
	}

	if !exists {
		return appErr.New(
			appErr.Internal,
			ErrStorageBucketMissing.Code,
			ErrStorageBucketMissing.Message,
		)
	}

	return nil
}
